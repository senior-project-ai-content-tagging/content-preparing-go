package publisher

import (
	"encoding/json"

	"cloud.google.com/go/pubsub"
	"github.com/senior-project-ai-content-tagging/content-preparing/entity"
)

type PreparedPublisher interface {
	Publish(message entity.TicketPubSub) error
}

type preparedPublisher struct {
	BasePublisher
}

func (p *preparedPublisher) Publish(message entity.TicketPubSub) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	p.topic.Publish(p.ctx, &pubsub.Message{
		Data: data,
	})

	return nil
}
