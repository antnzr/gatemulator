package repository

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/antnzr/gatemulator/internal/app/gatemulator/dto"
	"github.com/google/uuid"
)

const returnSubscription = "RETURNING id, subscription_id, subscriber_id, phone, state, transport, created_at, updated_at;"

type subscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (s *subscriptionRepository) Create(data dto.CreateSubscriptionEntity) (*Subscription, error) {
	query := `INSERT INTO subscriptions (id, subscription_id, subscriber_id, transport, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?, ?) ` + returnSubscription
	row := s.db.QueryRow(query, uuid.New(), data.SubscriptionID, data.SubscriberID, data.Transport, time.Now(), time.Now())

	subscription, err := scanRowIntoSubscription(row)
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

func (s *subscriptionRepository) GetBySubscriptionID(subscriptionId string) (*Subscription, error) {
	row := s.db.QueryRow(`
		SELECT id, subscription_id, subscriber_id, phone, state, transport, created_at, updated_at
		FROM subscriptions
		WHERE subscription_id = ?;
	`, subscriptionId)

	subscription, err := scanRowIntoSubscription(row)
	if err != nil {
		return nil, err
	}

	if subscription == nil {
		return nil, errors.New("subscription not found")
	}

	return subscription, nil
}

func (s *subscriptionRepository) Update(payload dto.UpdateSubscriptionRequest) (*Subscription, error) {
	query := `UPDATE subscriptions SET `
	qParts := make([]string, 0, 3)
	args := make([]interface{}, 0, 4)

	qParts = append(qParts, `updated_at = ?`)
	args = append(args, time.Now())

	if payload.Phone != nil {
		qParts = append(qParts, `phone = ?`)
		args = append(args, payload.Phone)
	}

	if payload.State != nil {
		qParts = append(qParts, `state = ?`)
		args = append(args, payload.State)
	}

	query += strings.Join(qParts, ",") + ` WHERE subscription_id = ? ` + returnSubscription
	args = append(args, payload.SubscriptionID)
	row := s.db.QueryRow(query, args...)

	subscription, err := scanRowIntoSubscription(row)
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

func (s *subscriptionRepository) Delete(subscriptionId string) error {
	panic("unimplemented")
}

func scanRowIntoSubscription(row *sql.Row) (*Subscription, error) {
	var subscription Subscription
	err := row.Scan(
		&subscription.ID,
		&subscription.SubscriptionID,
		&subscription.SubscriberID,
		&subscription.Phone,
		&subscription.State,
		&subscription.Transport,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}
