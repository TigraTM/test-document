package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func New() *viper.Viper {
	cfg := viper.New()

	// Postgres part of configibis.basis-plus.ru
	cfg.SetDefault("DB_HOST", "postgres")
	cfg.SetDefault("DB_PORT", "5432")
	cfg.SetDefault("DB_USER", "postgresql")
	cfg.SetDefault("DB_PASSWORD", "postgresql")
	cfg.SetDefault("DB_NAME", "postgresql")

	// Servers part of config
	cfg.SetDefault("LISTEN", ":8000")

	cfg.AutomaticEnv()

	cfg.Set("DB", dbConnectionString(cfg.GetString("DB_HOST"),
		cfg.GetString("DB_PORT"),
		cfg.GetString("DB_USER"),
		cfg.GetString("DB_PASSWORD"),
		cfg.GetString("DB_NAME")))

	return cfg
}

func dbConnectionString(host string, port string, user string, password string, dbname string) string {
	return fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}
