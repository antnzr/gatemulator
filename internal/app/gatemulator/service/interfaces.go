package service

import "github.com/antnzr/gatemulator/internal/app/gatemulator/dto"

type SubscriptionService interface {
	CreateSubscription(subscriber *dto.SubscriberResponse, payload *dto.CreateSubscriptionRequest) (*dto.CreateSubscriptionResponse, error)
	ScanQr(subscriber *dto.SubscriberResponse, payload *dto.ScanQrRequest) error
	Delete(subscriptionId string) error
	Enable(subscriber *dto.SubscriberResponse, subscriptionId string) error
}

type SubscriberService interface {
	CreateSubscriber(payload *dto.CreateSubscriberRequest) (*dto.SubscriberResponse, error)
	UpdateSubscriber(subscriberId string, payload *dto.UpdateSubscriberRequest) (*dto.SubscriberResponse, error)
	DeleteSubscriber(subscriberId string) error
	GetOneByToken(token string) (*dto.SubscriberResponse, error)
}

type MessageService interface {
	MessengerSendsMessage(subscriber *dto.SubscriberResponse, paylaod dto.MessengerSendsMessageRequest) dto.ChatUpdateWebhook
	NeedReadRecentChats(subscriber *dto.SubscriberResponse, subscriptionId string, payload *dto.ReadRecentChatMessagesRequest)
	PostMessage(subscriber *dto.SubscriberResponse, subscriptionId string, payload *dto.PostMessageRequest) (*dto.PostMessageResponse, error)
	ReadMessages(subscriber *dto.SubscriberResponse, payload *dto.ReadMessagesRequest)
}

type StoreService interface {
	GetFile(sha1 string) (*dto.MessageFileResponse, interface{})
}
