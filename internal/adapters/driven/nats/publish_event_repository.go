package nats

import (
	broker "github.com/Puena/auction-house-message-broker"
)

type publishEventRepository struct {
	b broker.Publisher
}

func NewPublishEventRepository(broker broker.Publisher) *publishEventRepository {
	return &publishEventRepository{b: broker}
}
