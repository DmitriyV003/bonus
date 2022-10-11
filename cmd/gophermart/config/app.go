package config

import (
	"context"
	"github.com/DmitriyV003/bonus/cmd/gophermart/container"
	"github.com/DmitriyV003/bonus/cmd/gophermart/routes/api"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
)

type App struct {
	pool      *pgxpool.Pool
	Conf      Config
	Container *container.Container
}

func (app *App) Start() http.Handler {
	router := chi.NewRouter()
	app.pool = app.connectToDB()

	if app.Conf.DatabaseDsn != "" && app.pool != nil {
		app.migrate()
	}

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Compress(5))
	router.Use(middleware.Heartbeat("/heartbeat"))

	privateApiRoutes := routes.Private{
		Container: app.Container,
	}

	router.Route("/api", func(r chi.Router) {
		r.Mount("/", privateApiRoutes.Routes())
	})

	return router
}

func (app *App) Config() {
	app.Conf.ParseEnv()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func (app *App) Close() {
	if app.pool != nil {
		app.pool.Close()
	}
}

func (app *App) connectToDB() (pool *pgxpool.Pool) {
	if app.Conf.DatabaseDsn == "" {
		log.Warn().Msg("Database URl not provided")
		return nil
	}

	var err error
	conf, err := pgxpool.ParseConfig(app.Conf.DatabaseDsn)
	if err != nil {
		log.Error().Err(err).Msg("Unable to parse Database config")
		return
	}
	pool, err = pgxpool.NewWithConfig(context.Background(), conf)

	if err != nil {
		log.Error().Err(err).Msg("Unable to connect to database")
		return
	}

	return pool
}

func (app *App) migrate() {

	//parsedDbUrl, _ := url.Parse(container.conf.DatabaseDsn)
	//cmd := exec.Command("tern", "migrate", "--migrations", "./migrations")
	//cmd.Env = append(cmd.Env, fmt.Sprintf("DATABASE=%s", strings.Trim(parsedDbUrl.Path, "/")))
	//cmd.Env = append(cmd.Env, fmt.Sprintf("DATABASE_DSN=%s", container.conf.DatabaseDsn))
	//out, err := cmd.CombinedOutput()
	//if err != nil {
	//	log.Error("Error during migrations: ", err)
	//	return
	//}
	//
	//log.Info("Migrating: ", string(out))

	sql := `CREATE TABLE IF NOT EXISTS metrics(
    	id serial PRIMARY KEY,
    	name VARCHAR (255) NOT NULL,
    	type VARCHAR (255) NOT NULL,
    	int_value BIGINT,
    	float_value DOUBLE PRECISION
	)`
	_, err := app.pool.Query(context.Background(), sql)
	if err != nil {
		log.Error().Err(err).Msg("Error during migration")
		return
	}

	log.Info().Msgf("Migrating: %s", sql)
}