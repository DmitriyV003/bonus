package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Address        string `env:"ADDRESS"`
	DatabaseURI    string `env:"DATABASE_URI"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"http://localhost:8080"`
	JwtSecret      string `env:"JWT_SECRET" envDefault:"jvf48g57h348f493fol-9m[=mp634b3p[n-89--fh23498gh4fgj3f4i[g4["`
}

const defaultAddress = "localhost:8080"
const defaultAccrualSystemAddress = "http://localhost:8080"
const defaultDatabaseDsn = ""

func (conf *Config) ParseEnv() {
	err := env.Parse(conf)
	if err != nil {
		log.Error().Err(err).Msg("Unable to parse ENV")
	}

	address := flag.String("a", defaultAddress, "Server address")
	databaseURI := flag.String("d", defaultDatabaseDsn, "connection string to database")
	accrualAddress := flag.String("r", defaultAccrualSystemAddress, "accrual address")
	flag.PrintDefaults()
	flag.Parse()

	log.Info().Fields(map[string]interface{}{
		"defaultAccrualSystemAddress": *accrualAddress,
	}).Msg("Env")

	if conf.Address == "" {
		conf.Address = *address
	}
	if conf.DatabaseURI == "" {
		conf.DatabaseURI = *databaseURI
	}
	if conf.AccrualAddress == "" {
		conf.AccrualAddress = *accrualAddress
	}
}
