package controller

import (
	"encoding/json"
	"net/http"

	"github.com/antnzr/gatemulator/internal/app/gatemulator/apperror"
	"github.com/antnzr/gatemulator/internal/app/gatemulator/dto"
	"github.com/antnzr/gatemulator/internal/app/gatemulator/service"
)

type SubscriberController struct {
	subscriberService service.SubscriberService
}

func NewSubscriberController(subscriberService service.SubscriberService) *SubscriberController {
	return &SubscriberController{subscriberService}
}

func (c *SubscriberController) CreateSubscriber(w http.ResponseWriter, r *http.Request) error {
	var body dto.CreateSubscriberRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return apperror.NewHTTPError(http.StatusBadRequest, "invalid request format "+err.Error())
	}

	result, err := c.subscriberService.CreateSubscriber(&body)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
	return nil
}

func (c *SubscriberController) UpdateSubscriber(w http.ResponseWriter, r *http.Request) error {
	var body dto.UpdateSubscriberRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return apperror.NewHTTPError(http.StatusBadRequest, "invalid request format "+err.Error())
	}

	subscriberId := r.PathValue("subscriberId")
	result, err := c.subscriberService.UpdateSubscriber(subscriberId, &body)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
	return nil
}

func (c *SubscriberController) DeleteSubscriber(w http.ResponseWriter, r *http.Request) error {
	subscriberId := r.PathValue("subscriberId")

	err := c.subscriberService.DeleteSubscriber(subscriberId)
	if err != nil {
		return apperror.NewHTTPError(http.StatusBadRequest, "invalid request format "+err.Error())
	}

	return nil
}
