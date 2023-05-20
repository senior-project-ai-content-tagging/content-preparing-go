package config

import "github.com/caarlos0/env/v8"

type AzureConfig struct {
	ApiKey string `env:"AZURE_API_KEY"`
}

func GetAzureConfig() AzureConfig {
	config := AzureConfig{}
	err := env.Parse(&config)
	if err != nil {
		panic(err)
	}

	return config
}
