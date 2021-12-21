package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port     string
	Postgres Postgres
	Jwt      JWT
}

type Postgres struct {
	Dsn string
}

type JWT struct {
	SigningKey         string
	AccessTimeExpired  int
	RefreshTimeExpires int
}

func Parse() *Config {
	envPort, exists := os.LookupEnv("PORT")
	if !exists {
		envPort = "8005"
	}

	pgDsn, exists := os.LookupEnv("POSTGRES_DSN")
	if !exists {
		pgDsn = "host=127.0.0.1 port=5434 user=postgres password=postgres dbname=library sslmode=disable"
	}

	signingKey, exists := os.LookupEnv("JWT_SIGNING_KEY")
	if !exists {
		signingKey = "asdlfn"
	}

	accExpToken, err := strconv.Atoi(os.Getenv("JWT_ACC_TOKEN_EXPIRED"))
	if err != nil {
		accExpToken = 60
	}
	refreshExpToken, err := strconv.Atoi(os.Getenv("JWT_REFRESH_TOKEN_EXPIRED"))
	if err != nil {
		refreshExpToken = 360
	}

	return &Config{
		Port:     envPort,
		Postgres: Postgres{Dsn: pgDsn},
		Jwt: JWT{
			SigningKey:         signingKey,
			AccessTimeExpired:  accExpToken,
			RefreshTimeExpires: refreshExpToken,
		},
	}
}
