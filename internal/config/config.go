package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Address     string `env:"ADDRESS"`
	DatabaseDsn string `env:"DATABASE_DSN"`
	JwtSecret   string `env:"JWT_SECRET" envDefault:"jvf48g57h348f493fol-9m[=mp634b3p[n-89--fh23498gh4fgj3f4i[g4["`
}

const defaultAddress = "localhost:8081"
const defaultDatabaseDsn = ""

func (conf *Config) ParseEnv() {
	err := env.Parse(conf)
	if err != nil {
		log.Error().Err(err).Msg("Unable to parse ENV")
	}

	address := flag.String("a", defaultAddress, "Server address")
	databaseDsn := flag.String("d", defaultDatabaseDsn, "connection string to database")
	flag.PrintDefaults()
	flag.Parse()

	if conf.Address == "" {
		conf.Address = *address
	}
	if conf.DatabaseDsn == "" {
		conf.DatabaseDsn = *databaseDsn
	}
}
