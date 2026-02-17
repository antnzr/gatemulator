package dto

import "github.com/google/uuid"

type CreateSubscriberRequest struct {
	Title      string `json:"title"`
	Token      string `json:"token"`
	WebhookUrl string `json:"webhookUrl"`
}

type CreateSubscriberEntity struct {
	Title      string `json:"title"`
	Token      string `json:"token"`
	WebhookUrl string `json:"webhookUrl"`
}

type UpdateSubscriberRequest struct {
	WebhookUrl *string `json:"webhookUrl"`
}

type UpdateSubscriberEntity struct {
	WebhookUrl *string `json:"webhookUrl"`
}

type SubscriberResponse struct {
	ID         uuid.UUID `json:"id"`
	Title      string    `json:"title"`
	Token      string    `json:"token"`
	WebhookUrl string    `json:"webhookUrl"`
}
