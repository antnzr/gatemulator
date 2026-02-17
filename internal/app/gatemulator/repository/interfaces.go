package repository

import (
	"github.com/antnzr/gatemulator/internal/app/gatemulator/dto"
)

type SubscriptionRepository interface {
	Create(payload dto.CreateSubscriptionEntity) (*Subscription, error)
	GetBySubscriptionID(subscriptionId string) (*Subscription, error)
	Update(payload dto.UpdateSubscriptionRequest) (*Subscription, error)
	Delete(subscriptionId string) error
}

type SubscriberRepository interface {
	Create(payload dto.CreateSubscriberEntity) (*Subscriber, error)
	Update(subscriberId string, payload dto.UpdateSubscriberEntity) (*Subscriber, error)
	Delete(subscriberId string) error
	FindOneByToken(token string) (*Subscriber, error)
}

type MessageFileRepository interface {
	Create(payload dto.CreateMessageFileDto) (*MessageFile, error)
	GetBySha1(sha1 string) *MessageFile
}
