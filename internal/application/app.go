package application

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/DmitriyV003/bonus/internal/config"
	"github.com/DmitriyV003/bonus/internal/container"
	"github.com/DmitriyV003/bonus/internal/repository"
	"github.com/DmitriyV003/bonus/internal/routes/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
)

type App struct {
	pool         *pgxpool.Pool
	Conf         config.Config
	Repositories *container.Repositories
	Services     *container.Services
}

func (app *App) CreateHandler() http.Handler {
	router := chi.NewRouter()
	app.pool = app.connectToDB()

	if app.Conf.DatabaseUri != "" && app.pool != nil {
		app.migrate()
	}

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Compress(5))
	router.Use(middleware.Heartbeat("/heartbeat"))

	app.Repositories.Users = repository.NewUserRepository(app.pool)
	app.Repositories.Orders = repository.NewOrderRepository(app.pool)
	app.Repositories.Payments = repository.NewPaymentRepository(app.pool)

	privateApiRoutes := routes.Private{
		Repositories: app.Repositories,
		Services:     app.Services,
		Conf:         &app.Conf,
	}

	router.Route("/api", func(r chi.Router) {
		r.Mount("/", privateApiRoutes.Routes())
	})

	return router
}

func (app *App) ApplyConfig() {
	app.Conf.ParseEnv()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func (app *App) Close() {
	if app.pool != nil {
		app.pool.Close()
	}
}

func (app *App) connectToDB() (pool *pgxpool.Pool) {
	if app.Conf.DatabaseUri == "" {
		log.Warn().Msg("Database URl not provided")
		return nil
	}

	var err error
	conf, err := pgxpool.ParseConfig(app.Conf.DatabaseUri)
	if err != nil {
		log.Error().Err(err).Msg("Unable to parse Database config")
		return
	}
	pool, err = pgxpool.ConnectConfig(context.Background(), conf)

	if err != nil {
		log.Error().Err(err).Msg("Unable to connect to database")
		return
	}

	return pool
}

// TODO: rewrite method normally
func (app *App) migrate() {
	sql := `CREATE TABLE IF NOT EXISTS migrations(
    	id serial PRIMARY KEY,
    	name VARCHAR (255) NOT NULL UNIQUE)`
	_, err := app.pool.Exec(context.Background(), sql)
	if err != nil {
		log.Error().Err(err).Msg("Error during migrationFile")
		return
	}

	log.Info().Msgf("Creating migrations table")

	migrations, err := os.ReadDir("migrations")
	if err != nil {
		log.Error().Err(err).Msg("unable to read migrations directory")
		return
	}

	var file *os.File
	defer file.Close()

	for _, migrationFile := range migrations {
		file, err = os.Open(fmt.Sprintf("migrations/%s", migrationFile.Name()))
		if err != nil {
			log.Error().Err(err).Msgf("unable to open migrationFile: %s", migrationFile.Name())
			return
		}

		sql = `SELECT id, name FROM migrations WHERE name = $1`
		var dbMigration migration

		err := app.pool.QueryRow(context.Background(), sql, migrationFile.Name()).Scan(&dbMigration.Id, &dbMigration.Name)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msg("Error to query migration")
			return
		}

		if err == nil {
			log.Info().Msgf("Migration passed: %s", migrationFile.Name())
			continue
		}

		wr := bytes.Buffer{}
		sc := bufio.NewScanner(file)
		for sc.Scan() {
			text := sc.Text()
			if text == "---- create above / drop below ----" {
				break
			}
			wr.WriteString(sc.Text())
		}

		sql = `INSERT INTO migrations (name) VALUES ($1)`
		_, err = app.pool.Exec(context.Background(), sql, migrationFile.Name())
		if err != nil {
			log.Error().Err(err).Msgf("unable to write migration to database: %s", migrationFile.Name())
			return
		}

		_, err = app.pool.Exec(context.Background(), wr.String())
		if err != nil {
			log.Error().Err(err).Msg("Error during migrationFile")
			return
		}

		log.Info().Msgf("Migrating: %s", wr.String())
	}
}

type migration struct {
	Id   int64
	Name string
}
