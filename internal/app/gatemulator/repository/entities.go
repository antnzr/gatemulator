package repository

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID             uuid.UUID `json:"id"`
	SubscriptionID uuid.UUID `json:"subscription_id"`
	SubscriberID   uuid.UUID `json:"subscriber_id"`
	Transport      string    `json:"transport"`
	Phone          *string   `json:"phone"`
	State          *string   `json:"state"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Subscriber struct {
	ID         uuid.UUID `json:"id"`
	Title      string    `json:"title"`
	Token      string    `json:"token"`
	WebhookUrl string    `json:"webhook_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type MessageFile struct {
	ID        uuid.UUID `json:"id"`
	MessageId uuid.UUID `json:"message_id"`
	SHA1      string    `json:"sha1"`
	MimeType  string    `json:"mime_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
