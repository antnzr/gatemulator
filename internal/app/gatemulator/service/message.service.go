package service

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	"github.com/antnzr/gatemulator/internal/app/gatemulator/dto"
	"github.com/antnzr/gatemulator/internal/app/gatemulator/repository"
	"github.com/antnzr/gatemulator/internal/pkg/gatemulator/utils"
	"github.com/google/uuid"
)

type messageService struct {
	dao repository.DAO
}

func NewMessageService(dao repository.DAO) MessageService {
	return &messageService{dao}
}

func (m *messageService) ReadMessages(subscriber *dto.SubscriberResponse, payload *dto.ReadMessagesRequest) {
	var body []dto.ChatUpdateWebhook

	// sent -> read
	chatUpdate := dto.ChatUpdateWebhook{
		Type:           "chatUpdate",
		Timestamp:      time.Now().Unix(),
		SubscriberID:   subscriber.ID.String(),
		SubscriptionID: payload.SubscriptionID,
		Chat: dto.Chat{
			ChatType: "whatsapp",
			ChatID:   payload.Phone,
		},
		UpdateMessagesBeforeTs: &[]dto.UpdateMessagesBeforeTs{
			{
				ChannelId: payload.SubscriptionID,
				OldStatus: "sent",
				NewStatus: "read",
				Timestamp: time.Now().UnixMilli(),
			},
		},
	}
	body = append(body, chatUpdate)
	m.send(body, subscriber.WebhookUrl)

	// delivered -> read
	time.AfterFunc(1*time.Second, func() {
		var body []dto.ChatUpdateWebhook

		chatUpdate := dto.ChatUpdateWebhook{
			Type:           "chatUpdate",
			Timestamp:      time.Now().Unix(),
			SubscriberID:   subscriber.ID.String(),
			SubscriptionID: payload.SubscriptionID,
			Chat: dto.Chat{
				ChatType: "whatsapp",
				ChatID:   payload.Phone,
			},
			UpdateMessagesBeforeTs: &[]dto.UpdateMessagesBeforeTs{
				{
					ChannelId: payload.SubscriptionID,
					OldStatus: "delivered",
					NewStatus: "read",
					Timestamp: time.Now().UnixMilli(),
				},
			},
		}
		body = append(body, chatUpdate)
		m.send(body, subscriber.WebhookUrl)
	})
}

func (m *messageService) PostMessage(subscriber *dto.SubscriberResponse, subscriptionId string, payload *dto.PostMessageRequest) (*dto.PostMessageResponse, error) {
	var response dto.PostMessageResponse
	if payload.Message.GUID != nil {
		response.MessageId = *payload.Message.GUID
	} else {
		response.MessageId = uuid.NewString()
	}

	for _, attachment := range payload.Message.Attachments {
		if media, ok := attachment.MessageAttachment.(dto.MediaAttachment); ok {
			m.dao.MessageFile().Create(dto.CreateMessageFileDto{
				MessageId: response.MessageId,
				SHA1:      media.Sha1,
				MimeType:  media.MimeType,
			})
		}
	}

	time.AfterFunc(3*time.Second, func() {
		var body []dto.ChatUpdateWebhook

		chatUpdate := dto.ChatUpdateWebhook{
			Type:           "chatUpdate",
			Timestamp:      time.Now().Unix(),
			SubscriberID:   subscriber.ID.String(),
			SubscriptionID: subscriptionId,
			Chat: dto.Chat{
				ChatType: payload.Chat.ChatType,
				ChatID:   payload.Chat.ChatId,
			},
			MessageUpdates: &[]dto.MessageUpdates{
				{
					GUID:        response.MessageId,
					MessengerID: utils.RandomWhatsappMessengerId(),
					Timestamp:   time.Now().UnixMilli(),
					Status:      "sent",
				},
			},
		}

		body = append(body, chatUpdate)
		m.send(body, subscriber.WebhookUrl)
	})

	time.AfterFunc(2*time.Second, func() {
		var body []dto.ChatUpdateWebhook

		chatUpdate := dto.ChatUpdateWebhook{
			Type:           "chatUpdate",
			Timestamp:      time.Now().Unix(),
			SubscriberID:   subscriber.ID.String(),
			SubscriptionID: subscriptionId,
			Chat: dto.Chat{
				ChatType: payload.Chat.ChatType,
				ChatID:   payload.Chat.ChatId,
			},
			UpdateMessagesBeforeTs: &[]dto.UpdateMessagesBeforeTs{
				{
					ChannelId: subscriptionId,
					OldStatus: "sent",
					NewStatus: "delivered",
					Timestamp: time.Now().UnixMilli(),
				},
			},
		}
		body = append(body, chatUpdate)
		m.send(body, subscriber.WebhookUrl)
	})

	return &response, nil
}

func (m *messageService) NeedReadRecentChats(subscriber *dto.SubscriberResponse, subscriptionId string, payload *dto.ReadRecentChatMessagesRequest) {
	baseDelay := 2 * time.Second

	chatsCount := 3
	if payload.ChatsCount != nil && *payload.ChatsCount >= 3 && *payload.ChatsCount <= 5 {
		chatsCount = *payload.ChatsCount
	}

	for i := 0; i < chatsCount; i++ {
		delay := time.Duration(i+1) * baseDelay
		fn := m.readRecentChatFn(subscriber, subscriptionId)
		time.AfterFunc(delay, fn)
	}
}

func (m *messageService) MessengerSendsMessage(
	subscriber *dto.SubscriberResponse,
	payload dto.MessengerSendsMessageRequest,
) dto.ChatUpdateWebhook {
	var messageText string

	if payload.MessageText == nil {
		messageText = utils.RandomString(12)
	} else {
		messageText = *payload.MessageText
	}

	chatType := "whatsapp"
	webhookType := "chatUpdate"
	contactName := utils.RandomName()
	webhookBody := dto.ChatUpdateWebhook{
		Type:           webhookType,
		Timestamp:      time.Now().Unix(),
		SubscriberID:   subscriber.ID.String(),
		SubscriptionID: payload.SubscriptionID.String(),
		Chat: dto.Chat{
			ChatType: chatType,
			ChatID:   payload.ReceiverPhone,
			Contact: &dto.Contact{
				ID:       payload.ReceiverPhone,
				ChatType: chatType,
				Name:     contactName,
				Phone:    payload.ReceiverPhone,
				Details: dto.ContactDetails{
					Phone: payload.ReceiverPhone,
				},
			},
		},
		Sender: &dto.MessageSender{
			Title:       "messenger",
			IsMessenger: true,
		},
		Messages: &[]dto.Message{
			{
				GUID:        uuid.NewString(),
				Text:        messageText,
				MessengerID: utils.RandomWhatsappMessengerId(),
				Type:        "text",
				Timestamp:   time.Now().UnixMilli(),
				Status:      "incoming",
				Author: dto.MessageAuthor{
					Name: contactName,
				},
			},
		},
	}

	var body []dto.ChatUpdateWebhook
	body = append(body, webhookBody)
	m.send(body, subscriber.WebhookUrl)

	return webhookBody
}

func (m *messageService) readRecentChatFn(subscriber *dto.SubscriberResponse, subscriptionId string) func() {
	return func() {
		chatUpdate := m.generateRandomChat(subscriber.ID.String(), subscriptionId)
		var body []dto.ChatUpdateWebhook
		body = append(body, chatUpdate)
		m.send(body, subscriber.WebhookUrl)
	}
}

func (m *messageService) generateRandomChat(subscriberId string, subscriptionId string) dto.ChatUpdateWebhook {
	numMessages := 100
	currentTime := time.Now()
	contactPhone := utils.RandomPhoneNumber()
	contactName := utils.RandomName()
	chatType := "whatsapp"
	webhookType := "chatUpdate"

	chatUpdate := dto.ChatUpdateWebhook{
		Type:           webhookType,
		Timestamp:      currentTime.Unix(),
		SubscriberID:   subscriberId,
		SubscriptionID: subscriptionId,
		Chat: dto.Chat{
			ChatType: chatType,
			ChatID:   contactPhone,
			Contact: &dto.Contact{
				ID:       contactPhone,
				ChatType: chatType,
				Name:     contactName,
				Details: dto.ContactDetails{
					Phone: contactPhone,
				},
				Phone: contactPhone,
			},
		},
		Sender: &dto.MessageSender{
			Title:       "messenger",
			IsMessenger: true,
		},
	}
	messages := make([]dto.Message, 0)

	for i := 0; i < numMessages; i++ {
		messageTime := currentTime.Add(-time.Duration(i)*time.Hour*14).UnixNano() / int64(time.Millisecond)

		message := dto.Message{
			GUID:        uuid.New().String(),
			MessengerID: utils.RandomWhatsappMessengerId(),
			Type:        "text",
			Text:        utils.RandomString(12),
			Attachments: []string{},
			Timestamp:   messageTime,
			Status:      m.randomStatus(),
			Author: dto.MessageAuthor{
				Name: contactName,
			},
			Details: dto.ContactDetails{},
		}

		messages = append(messages, message)
	}
	chatUpdate.Messages = &messages

	return chatUpdate
}

func (m *messageService) randomStatus() string {
	statuses := []string{"read", "incoming"}
	return statuses[rand.Intn(len(statuses))]
}

func (m *messageService) send(webhookBody interface{}, url string) {
	postBody, _ := json.Marshal(webhookBody)
	body := bytes.NewBuffer(postBody)

	res, err := http.Post(url, "application/json", body)
	if err != nil {
		slog.Warn("MessageService", slog.String("tag", "send"), slog.Any("err", err))
	} else {
		slog.Info("MessageService", slog.String("tag", "send"), slog.String("url", url), slog.String("responseStatus", res.Status))
	}
}
