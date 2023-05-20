package config

import "github.com/caarlos0/env/v8"

type DatabaseConfig struct {
	Host string `env:"DB_HOST"`
	Port string `env:"DB_PORT"`
	User string `env:"DB_USER"`
	Pass string `env:"DB_PASS"`
	Name string `env:"DB_NAME"`
	SSL  string `env:"DB_SSL"`
}

func GetDatabaseConfig() DatabaseConfig {
	config := DatabaseConfig{}
	err := env.Parse(&config)
	if err != nil {
		panic(err)
	}

	return config
}
