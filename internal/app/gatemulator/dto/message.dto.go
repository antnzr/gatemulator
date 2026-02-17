package dto

import (
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"
)

type MessengerSendsMessageRequest struct {
	SubscriptionID uuid.UUID `json:"subscriptionId"`
	ReceiverPhone  string    `json:"receiverPhone"`
	MessageText    *string   `json:"messageText"`
}

type MessageUpdates struct {
	GUID        string `json:"guid"`
	MessengerID string `json:"messengerId"`
	Timestamp   int64  `json:"timestamp"`
	Status      string `json:"status"`
}

type UpdateMessagesBeforeTs struct {
	ChannelId string `json:"channelId"`
	OldStatus string `json:"oldStatus"`
	NewStatus string `json:"newStatus"`
	Timestamp int64  `json:"timestamp"`
}

type ChatUpdateWebhook struct {
	Type                   string                    `json:"type"`
	Timestamp              int64                     `json:"timestamp"`
	SubscriberID           string                    `json:"subscriberId"`
	SubscriptionID         string                    `json:"subscriptionId"`
	Chat                   Chat                      `json:"chat"`
	Sender                 *MessageSender            `json:"sender,omitempty"`
	Messages               *[]Message                `json:"messages,omitempty"`
	MessageUpdates         *[]MessageUpdates         `json:"messageUpdates,omitempty"`
	UpdateMessagesBeforeTs *[]UpdateMessagesBeforeTs `json:"updateMessagesBeforeTs,omitempty"`
}

type Chat struct {
	ChatType string   `json:"chatType"`
	ChatID   string   `json:"chatId"`
	Contact  *Contact `json:"contact,omitempty"`
}

type Contact struct {
	ID       string         `json:"id"`
	ChatType string         `json:"chatType"`
	Name     string         `json:"name,omitempty"`
	Avatar   ContactAvatar  `json:"avatar,omitempty"`
	Details  ContactDetails `json:"details,omitempty"`
	Phone    string         `json:"phone,omitempty"`
}

type ContactAvatar struct {
	SHA1 string `json:"sha1,omitempty"`
}

type ContactDetails struct {
	Phone string `json:"phone,omitempty"`
}

type MessageSender struct {
	Title       string `json:"title,omitempty"`
	IsMessenger bool   `json:"isMessenger,omitempty"`
}

type Message struct {
	GUID        string        `json:"guid"`
	MessengerID string        `json:"messengerId"`
	Type        string        `json:"type"`
	Text        string        `json:"text"`
	Attachments []string      `json:"attachments,omitempty"`
	Timestamp   int64         `json:"timestamp"`
	Status      string        `json:"status,omitempty"`
	Author      MessageAuthor `json:"author,omitempty"`
	Details     interface{}   `json:"details,omitempty"`
}

type MessageAuthor struct {
	Name string `json:"name,omitempty"`
}

type ReadRecentChatMessagesRequest struct {
	ChatsCount *int `json:"chatsCount,omitempty"`
}

type MessageEntryRequest struct {
	Type           string       `json:"type"`
	GUID           *string      `json:"guid,omitempty"`
	RefMessengerId *string      `json:"refMessengerId,omitempty"`
	Text           *string      `json:"text,omitempty"`
	Timestamp      *int64       `json:"timestamp,omitempty"`
	Status         *string      `json:"status,omitempty"`
	Attachments    []Attachment `json:"attachments,omitempty"`
}

type ChatEntryRequest struct {
	ChatType string `json:"chatType"`
	ChatId   string `json:"chatId"`
}

type PostMessageRequest struct {
	Message MessageEntryRequest `json:"message"`
	Chat    ChatEntryRequest    `json:"chat"`
}

type PostMessageResponse struct {
	MessageId string `json:"messageId"`
}

type ReadMessagesRequest struct {
	SubscriptionID string `json:"subscriptionId"`
	Phone          string `json:"phone"`
}

type MediaAttachmentType string

const (
	MediaAttachmentTypeImage    MediaAttachmentType = "image"
	MediaAttachmentTypeAudio    MediaAttachmentType = "audio"
	MediaAttachmentTypeVideo    MediaAttachmentType = "video"
	MediaAttachmentTypeDocument MediaAttachmentType = "document"
)

type LocationAttachment struct {
	Type      string  `json:"type"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type MediaAttachment struct {
	Type     MediaAttachmentType `json:"type"`
	Name     string              `json:"name"`
	Size     int                 `json:"size"`
	Sha1     string              `json:"sha1"`
	MimeType string              `json:"mimetype"`
	Link     string              `json:"link"`
}

type CreateMessageFileDto struct {
	MessageId string
	SHA1      string
	MimeType  string
}

type MessageFileResponse struct {
	ID        uuid.UUID `json:"id"`
	MessageId string    `json:"messageId"`
	SHA1      string    `json:"sha1"`
	MimeType  string    `json:"mimetype"`
}

type MessageAttachment interface{}

type Attachment struct {
	MessageAttachment
}

func (a *Attachment) UnmarshalJSON(data []byte) error {
	var typeCheck struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &typeCheck); err != nil {
		return err
	}

	switch typeCheck.Type {
	case "image", "audio", "video", "document":
		var media MediaAttachment
		if err := json.Unmarshal(data, &media); err != nil {
			return err
		}
		a.MessageAttachment = media
	case "location":
		var location LocationAttachment
		if err := json.Unmarshal(data, &location); err != nil {
			return err
		}
		a.MessageAttachment = location
	default:
		slog.Info("UnmarshalJSON", "unknown attachment type", slog.String("typeCheck.Type", typeCheck.Type))
		return nil
	}
	return nil
}
