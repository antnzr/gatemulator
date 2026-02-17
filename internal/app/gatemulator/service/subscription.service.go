package service

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/antnzr/gatemulator/internal/app/gatemulator/apperror"
	"github.com/antnzr/gatemulator/internal/app/gatemulator/constant"
	"github.com/antnzr/gatemulator/internal/app/gatemulator/dto"
	"github.com/antnzr/gatemulator/internal/app/gatemulator/repository"
	"github.com/antnzr/gatemulator/internal/pkg/gatemulator/job"
	"github.com/antnzr/gatemulator/internal/pkg/gatemulator/qrcode"
	"github.com/google/uuid"
)

type subscriptionService struct {
	dao        repository.DAO
	jobManager *job.JobManager
}

func NewSubscriptionService(dao repository.DAO, jobManager *job.JobManager) SubscriptionService {
	return &subscriptionService{dao, jobManager}
}

func (s *subscriptionService) Enable(subscriber *dto.SubscriberResponse, subscriptionId string) error {
	subscription, err := s.dao.Subscription().GetBySubscriptionID(subscriptionId)
	if err != nil {
		return apperror.NewHTTPError(http.StatusNotFound, err.Error())
	}

	if subscription.State != nil && *subscription.State == constant.StateActive {
		return apperror.NewHTTPError(http.StatusBadRequest, "already active")
	}

	state := constant.StateQR
	updateData := dto.UpdateSubscriptionRequest{
		SubscriptionID: uuid.MustParse(subscriptionId),
		State:          &state,
	}
	s.dao.Subscription().Update(updateData)

	jobId := "qr:" + subscriptionId
	jobDelay := 4 * time.Second
	jobHandler := s.sendQr(subscription.SubscriptionID, subscriber)
	s.jobManager.ScheduleRepeatable(jobId, jobDelay, jobHandler)

	return nil
}

func (s *subscriptionService) Delete(subscriptionId string) error {
	state := constant.StateDeleted
	updateData := dto.UpdateSubscriptionRequest{
		SubscriptionID: uuid.MustParse(subscriptionId),
		State:          &state,
	}
	s.dao.Subscription().Update(updateData)
	return nil
}

func (s *subscriptionService) ScanQr(subscriber *dto.SubscriberResponse, payload *dto.ScanQrRequest) error {
	slog.Info("SubscriptionService", slog.String("tag", "ScanQr"), slog.Any("payload", payload))

	state := constant.StateActive
	updateData := dto.UpdateSubscriptionRequest{
		SubscriptionID: payload.SubscriptionID,
		Phone:          &payload.Phone,
		State:          &state,
	}
	entity, err := s.dao.Subscription().Update(updateData)
	if err != nil {
		return apperror.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	s.jobManager.Delete("qr:" + payload.SubscriptionID.String())

	if entity == nil {
		return apperror.NewHTTPError(http.StatusBadRequest, "subscription not found")
	}

	time.AfterFunc(4*time.Second, func() {
		var webhookBody dto.UpdateChannelsWehhook
		needReadRecentChats := true

		webhookUpdateData := dto.UpdateChannelWehhook{
			Type:           "subscriptionUpdate",
			SubscriptionID: entity.SubscriptionID,
			Timestamp:      time.Now().Unix(),
			Whatsapp: &dto.TransportData{
				State: entity.State,
				Phone: entity.Phone,
				Details: &dto.TransportDetails{
					NeedReadRecentChats: &needReadRecentChats,
				},
			},
		}

		webhookBody = append(webhookBody, webhookUpdateData)
		s.send(webhookBody, subscriber.WebhookUrl)
	})

	return nil
}

func (s *subscriptionService) CreateSubscription(subscriber *dto.SubscriberResponse, data *dto.CreateSubscriptionRequest) (*dto.CreateSubscriptionResponse, error) {
	subscriptionId := uuid.New()
	entity, err := s.dao.Subscription().Create(dto.CreateSubscriptionEntity{
		SubscriberID:   subscriber.ID,
		SubscriptionID: subscriptionId,
		Transport:      data.Transport,
	})
	if err != nil {
		return nil, apperror.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	subscription := dto.CreateSubscriptionResponse{
		SubscriptionID: subscriptionId,
	}

	if entity.Transport == constant.TransportWAPI {
		subscription.Details = &dto.SubscriptionDetails{
			OnboardingLink: constant.OnboardingLink,
		}
	}

	jobId := "qr:" + subscriptionId.String()
	jobDelay := 2 * time.Second
	jobHandler := s.sendQr(subscriptionId, subscriber)
	s.jobManager.ScheduleRepeatable(jobId, jobDelay, jobHandler)

	return &subscription, nil
}

func (s *subscriptionService) sendQr(subscriptionId uuid.UUID, subscriber *dto.SubscriberResponse) func() {
	return func() {
		jobId := "qr:" + subscriptionId.String()

		// no need to send qr if it's in the 'active' state
		channelSubscription, err := s.dao.Subscription().GetBySubscriptionID(subscriptionId.String())
		if err != nil || (channelSubscription.State != nil && *channelSubscription.State == constant.StateActive) {
			s.jobManager.Delete(jobId)
			return
		}

		count := s.jobManager.GetJobState(jobId)
		if count > 5 {
			state := constant.StateQRIdle
			payload := dto.UpdateSubscriptionRequest{
				SubscriptionID: subscriptionId,
				State:          &state,
			}
			s.dao.Subscription().Update(payload)
			s.jobManager.Delete(jobId)
			return
		}

		qr, err := qrcode.Generate(subscriptionId.String())
		if err != nil {
			slog.Info("SubscriptionService", slog.String("tag", "sendQr"), slog.Any("err", err))
			return
		}

		state := constant.StateQR
		var webhookBody dto.UpdateChannelsWehhook
		updateData := dto.UpdateChannelWehhook{
			Type:           "subscriptionUpdate",
			SubscriptionID: subscriptionId,
			Timestamp:      time.Now().Unix(),
			Whatsapp: &dto.TransportData{
				State:  &state,
				QrCode: &qr,
			},
		}
		webhookBody = append(webhookBody, updateData)
		s.send(webhookBody, subscriber.WebhookUrl)

		s.jobManager.ScheduleRepeatable("qr:"+subscriptionId.String(), 7*time.Second, s.sendQr(subscriptionId, subscriber))
	}
}

func (s *subscriptionService) send(webhookBody interface{}, url string) {
	postBody, _ := json.Marshal(webhookBody)
	body := bytes.NewBuffer(postBody)

	res, err := http.Post(url, "application/json", body)
	if err != nil {
		slog.Warn("SubscriptionService", slog.String("tag", "send"), slog.Any("err", err))
	} else {
		slog.Info("SubscriptionService", slog.String("tag", "send"), slog.String("url", url), slog.String("responseStatus", res.Status))
	}
}
