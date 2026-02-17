package repository

import (
	"database/sql"
	"strings"
	"time"

	"github.com/antnzr/gatemulator/internal/app/gatemulator/dto"
	"github.com/google/uuid"
)

const returnSubscriber = "RETURNING id, title, token, webhook_url, created_at, updated_at;"

type subscriberRepository struct {
	db *sql.DB
}

func NewSubscriberRepository(db *sql.DB) SubscriberRepository {
	return &subscriberRepository{db}
}

func (s *subscriberRepository) FindOneByToken(token string) (*Subscriber, error) {
	row := s.db.QueryRow(`SELECT id, title, token, webhook_url, created_at, updated_at FROM subscribers WHERE token = ?`, token)

	subscriber, err := scanRowIntoSubscriber(row)
	if err != nil {
		return nil, err
	}

	return subscriber, nil
}

func (s *subscriberRepository) Create(payload dto.CreateSubscriberEntity) (*Subscriber, error) {
	query := `INSERT INTO subscribers (id, title, token, webhook_url, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?) ` + returnSubscriber
	row := s.db.QueryRow(query, uuid.New(), payload.Title, payload.Token, payload.WebhookUrl, time.Now(), time.Now())

	subscriber, err := scanRowIntoSubscriber(row)
	if err != nil {
		return nil, err
	}

	return subscriber, nil
}

func (s *subscriberRepository) Delete(subscriberId string) error {
	_, err := s.db.Exec(`DELETE FROM subscribers WHERE id = ?;`, subscriberId)
	if err != nil {
		return err
	}
	return nil
}

func (s *subscriberRepository) Update(subscriberId string, payload dto.UpdateSubscriberEntity) (*Subscriber, error) {
	query := `UPDATE subscribers SET `
	qParts := make([]string, 0)
	args := make([]interface{}, 0)

	qParts = append(qParts, `updated_at = ?`)
	args = append(args, time.Now())

	if payload.WebhookUrl != nil {
		qParts = append(qParts, `webhook_url = ?`)
		args = append(args, payload.WebhookUrl)
	}

	query += strings.Join(qParts, ",") + ` WHERE id = ? ` + returnSubscriber
	args = append(args, subscriberId)
	row := s.db.QueryRow(query, args...)

	subscriber, err := scanRowIntoSubscriber(row)
	if err != nil {
		return nil, err
	}

	return subscriber, nil
}

func scanRowIntoSubscriber(row *sql.Row) (*Subscriber, error) {
	var subscriber Subscriber
	err := row.Scan(
		&subscriber.ID,
		&subscriber.Title,
		&subscriber.Token,
		&subscriber.WebhookUrl,
		&subscriber.CreatedAt,
		&subscriber.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &subscriber, nil
}
