package config

import "github.com/caarlos0/env/v8"

type PublisherConfig struct {
	ProjectID       string `env:"PUBSUB_PROJECT_ID"`
	PreparedTopicID string `env:"PUBSUB_PREPARED_TOPIC_ID"`
}

func GetPublisherConfig() PublisherConfig {
	config := PublisherConfig{}
	err := env.Parse(&config)
	if err != nil {
		panic(err)
	}

	return config
}
