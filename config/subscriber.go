package config

import "github.com/caarlos0/env/v8"

type SubscriberConfig struct {
	ProjectId        string `env:"PUBSUB_PROJECT_ID"`
	SubSubmitWeblink string `env:"PUBSUB_SUBSCRIPTION_WEBLINK_NAME"`
}

func GetSubConfig() SubscriberConfig {
	config := SubscriberConfig{}
	err := env.Parse(&config)
	if err != nil {
		panic(err)
	}

	return config
}
