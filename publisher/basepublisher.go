package publisher

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
)

type BasePublisher struct {
	topic *pubsub.Topic
	ctx   context.Context
}

func NewPublisher(ctx context.Context, projectID string, topicID string) (*BasePublisher, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	topic := client.Topic(topicID)

	return &BasePublisher{topic: topic, ctx: ctx}, nil
}

func (p *BasePublisher) Publish(message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	p.topic.Publish(p.ctx, &pubsub.Message{
		Data: data,
	})

	return nil
}
