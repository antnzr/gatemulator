package service

import (
	"net/http"

	"github.com/antnzr/gatemulator/internal/app/gatemulator/apperror"
	"github.com/antnzr/gatemulator/internal/app/gatemulator/dto"
	"github.com/antnzr/gatemulator/internal/app/gatemulator/repository"
)

type subscriberService struct {
	dao repository.DAO
}

func NewSubscriberService(dao repository.DAO) SubscriberService {
	return &subscriberService{dao}
}

func (s *subscriberService) GetOneByToken(token string) (*dto.SubscriberResponse, error) {
	entity, err := s.dao.Subscriber().FindOneByToken(token)
	if err != nil {
		return nil, err
	}

	if entity == nil {
		return nil, apperror.NewHTTPError(http.StatusNotFound, "subscriber not found")
	}

	return &dto.SubscriberResponse{
		ID:         entity.ID,
		Title:      entity.Title,
		Token:      entity.Token,
		WebhookUrl: entity.WebhookUrl,
	}, nil
}

func (s *subscriberService) CreateSubscriber(payload *dto.CreateSubscriberRequest) (*dto.SubscriberResponse, error) {
	entity, err := s.dao.Subscriber().Create(dto.CreateSubscriberEntity{
		Title:      payload.Title,
		Token:      payload.Token,
		WebhookUrl: payload.WebhookUrl,
	})
	if err != nil {
		return nil, apperror.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return &dto.SubscriberResponse{
		ID:         entity.ID,
		Title:      entity.Title,
		Token:      entity.Token,
		WebhookUrl: entity.WebhookUrl,
	}, nil
}

func (s *subscriberService) UpdateSubscriber(subscriberId string, payload *dto.UpdateSubscriberRequest) (*dto.SubscriberResponse, error) {
	entity, err := s.dao.Subscriber().Update(subscriberId, dto.UpdateSubscriberEntity{
		WebhookUrl: payload.WebhookUrl,
	})
	if err != nil {
		return nil, apperror.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return &dto.SubscriberResponse{
		ID:         entity.ID,
		Title:      entity.Title,
		Token:      entity.Token,
		WebhookUrl: entity.WebhookUrl,
	}, nil
}

func (s *subscriberService) DeleteSubscriber(subscriberId string) error {
	err := s.dao.Subscriber().Delete(subscriberId)
	if err != nil {
		return apperror.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}
