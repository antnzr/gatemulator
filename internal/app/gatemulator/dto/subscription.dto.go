package dto

import "github.com/google/uuid"

type CreateSubscriptionRequest struct {
	Transport string `json:"transport"`
}

type SubscriptionDetails struct {
	OnboardingLink string `json:"onboardingLink,omitempty"`
}

type CreateSubscriptionResponse struct {
	SubscriptionID uuid.UUID            `json:"subscriptionId"`
	Details        *SubscriptionDetails `json:"details,omitempty"`
}

type UpdateSubscriptionRequest struct {
	SubscriptionID uuid.UUID            `json:"subscriptionId"`
	Details        *SubscriptionDetails `json:"details"`
	Phone          *string              `json:"phone"`
	State          *string              `json:"state"`
}

type UpdateChannelsWehhook []UpdateChannelWehhook

type UpdateChannelWehhook struct {
	Type           string         `json:"type"`
	SubscriptionID uuid.UUID      `json:"subscriptionId"`
	Timestamp      int64          `json:"timestamp,omitempty"`
	Whatsapp       *TransportData `json:"whatsapp,omitempty"`
	Tgapi          *TransportData `json:"tgapi,omitempty"`
	Wapi           *TransportData `json:"wapi,omitempty"`
}

type TelegramChannelUser struct {
	ID        *string
	Phone     *string
	Username  *string
	FirstName *string
	LastName  *string
}

type TransportData struct {
	QrCode  *string           `json:"qrCode,omitempty"`
	Phone   *string           `json:"phone,omitempty"`
	State   *string           `json:"state,omitempty"`
	Details *TransportDetails `json:"details,omitempty"`
}

type TransportDetails struct {
	Hint                *string              `json:"hint,omitempty"`
	User                *TelegramChannelUser `json:"user,omitempty"`
	NeedReadRecentChats *bool                `json:"needReadRecentChats,omitempty"`
	WabaId              *string              `json:"wabaId,omitempty"`
	Tier                *string              `json:"tier,omitempty"`
	WabaName            *string              `json:"wabaName,omitempty"`
}

type CreateSubscriptionEntity struct {
	SubscriptionID uuid.UUID `json:"subscriptionId"`
	SubscriberID   uuid.UUID `json:"subscriberId"`
	Transport      string    `json:"transport"`
}

type SubscriptionResponse struct {
	ID             uuid.UUID `json:"id"`
	SubscriptionID uuid.UUID `json:"subscriptionId"`
	SubscriberID   uuid.UUID `json:"subscriberId"`
	Transport      string    `json:"transport,omitempty"`
	Phone          *string   `json:"phone,omitempty"`
	State          *string   `json:"state,omitempty"`
	CreatedAt      uuid.Time `json:"createdAt,omitempty"`
}

type ScanQrRequest struct {
	SubscriptionID uuid.UUID `json:"subscriptionId"`
	Phone          string    `json:"phone"`
}

type EnableSubscriptionRequest struct {
	SubscriptionID uuid.UUID `json:"subscriptionId"`
	Phone          string    `json:"phone"`
	Transport      string    `json:"transport"`
}
