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

type SubscriptionController struct {
	service service.SubscriptionService
}

func NewSubscriptionController(service service.SubscriptionService) *SubscriptionController {
	return &SubscriptionController{service: service}
}

func (c *SubscriptionController) Create(w http.ResponseWriter, r *http.Request) error {
	subscriber := r.Context().Value(middleware.SubscriberKey("subscriber")).(*dto.SubscriberResponse)
	slog.Info("SubscriptionController", slog.String("tag", "Create"), slog.String("subscriberID", subscriber.ID.String()))

	var body dto.CreateSubscriptionRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return apperror.NewHTTPError(http.StatusBadRequest, "invalid request format "+err.Error())
	}

	subscription, err := c.service.CreateSubscription(subscriber, &body)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(subscription)
	return nil
}

func (c *SubscriptionController) ScanQr(w http.ResponseWriter, r *http.Request) error {
	subscriber := r.Context().Value(middleware.SubscriberKey("subscriber")).(*dto.SubscriberResponse)
	slog.Info("SubscriptionController", slog.String("tag", "ScanQr"), slog.String("subscriberID", subscriber.ID.String()))

	var body dto.ScanQrRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return apperror.NewHTTPError(http.StatusBadRequest, "invalid request format "+err.Error())
	}

	err = c.service.ScanQr(subscriber, &body)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (c *SubscriptionController) Delete(w http.ResponseWriter, r *http.Request) error {
	subscriber := r.Context().Value(middleware.SubscriberKey("subscriber")).(*dto.SubscriberResponse)
	subscriptionId := r.PathValue("subscriptionId")
	slog.Info("SubscriptionController", slog.String("tag", "Delete"), slog.String("subscriptionId", subscriptionId), slog.String("subscriberID", subscriber.ID.String()))

	c.service.Delete(subscriptionId)

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (c *SubscriptionController) Enable(w http.ResponseWriter, r *http.Request) error {
	subscriber := r.Context().Value(middleware.SubscriberKey("subscriber")).(*dto.SubscriberResponse)
	subscriptionId := r.PathValue("subscriptionId")
	slog.Info("SubscriptionController", slog.String("tag", "Enable"), slog.String("subscriptionId", subscriptionId), slog.String("subscriberID", subscriber.ID.String()))

	var body dto.EnableSubscriptionRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return apperror.NewHTTPError(http.StatusBadRequest, "invalid request format "+err.Error())
	}

	err = c.service.Enable(subscriber, subscriptionId)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
