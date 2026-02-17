package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/antnzr/gatemulator/internal/app/gatemulator/apperror"
	"github.com/antnzr/gatemulator/internal/app/gatemulator/dto"
	"github.com/antnzr/gatemulator/internal/app/gatemulator/middleware"
	"github.com/antnzr/gatemulator/internal/app/gatemulator/service"
)

type MessageController struct {
	messageService service.MessageService
}

func NewMessageController(messageService service.MessageService) *MessageController {
	return &MessageController{messageService}
}

func (c *MessageController) ReadMessages(w http.ResponseWriter, r *http.Request) error {
	subscriber := r.Context().Value(middleware.SubscriberKey("subscriber")).(*dto.SubscriberResponse)
	slog.Info("ReadMessages", slog.String("read messages for subscriber", subscriber.ID.String()))

	var body dto.ReadMessagesRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return apperror.NewHTTPError(http.StatusBadRequest, "invalid request format "+err.Error())
	}

	c.messageService.ReadMessages(subscriber, &body)

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (c *MessageController) MessengerSendsMessage(w http.ResponseWriter, r *http.Request) error {
	subscriber := r.Context().Value(middleware.SubscriberKey("subscriber")).(*dto.SubscriberResponse)
	slog.Info("MessengerSendsMessage", slog.String("messenger sends message for subscriber", subscriber.ID.String()))

	var body dto.MessengerSendsMessageRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return apperror.NewHTTPError(http.StatusBadRequest, "invalid request format "+err.Error())
	}

	response := c.messageService.MessengerSendsMessage(subscriber, body)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	return nil
}

func (c *MessageController) NeedReadRecentChats(w http.ResponseWriter, r *http.Request) error {
	subscriber := r.Context().Value(middleware.SubscriberKey("subscriber")).(*dto.SubscriberResponse)
	subscriptionId := r.PathValue("subscriptionId")
	slog.Info("NeedReadRecentChats", slog.String("subscriberId", subscriber.ID.String()), slog.String("subscriptionId", subscriptionId))

	var body dto.ReadRecentChatMessagesRequest
	json.NewDecoder(r.Body).Decode(&body)

	c.messageService.NeedReadRecentChats(subscriber, subscriptionId, &body)

	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (c *MessageController) PostMessage(w http.ResponseWriter, r *http.Request) error {
	subscriber := r.Context().Value(middleware.SubscriberKey("subscriber")).(*dto.SubscriberResponse)
	subscriptionId := r.PathValue("subscriptionId")
	slog.Info("PostMessage", slog.String("subscriberId", subscriber.ID.String()), slog.String("subscriptionId", subscriptionId))

	var body dto.PostMessageRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return apperror.NewHTTPError(http.StatusBadRequest, "invalid request format "+err.Error())
	}

	response, err := c.messageService.PostMessage(subscriber, subscriptionId, &body)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	return nil
}
