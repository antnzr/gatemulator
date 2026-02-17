package repository

import (
	"database/sql"
	"time"

	"github.com/antnzr/gatemulator/internal/app/gatemulator/dto"
	"github.com/google/uuid"
)

const returnMessageFile = "RETURNING id, message_id, sha1, mime_type, created_at, updated_at;"

type messageFileRepository struct {
	db *sql.DB
}

func NewMessageFileRepository(db *sql.DB) MessageFileRepository {
	return &messageFileRepository{db}
}

func (m *messageFileRepository) GetBySha1(sha1 string) *MessageFile {
	row := m.db.QueryRow(`SELECT id, message_id, sha1, mime_type, created_at, updated_at FROM message_files WHERE sha1 = ?;`, sha1)

	entity, err := scanRowIntoMessageFile(row)
	if err != nil {
		return nil
	}

	return entity
}

func (m *messageFileRepository) Create(payload dto.CreateMessageFileDto) (*MessageFile, error) {
	query := `INSERT INTO message_files (id, message_id, sha1, mime_type, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?) ` + returnMessageFile
	row := m.db.QueryRow(query, uuid.New(), payload.MessageId, payload.SHA1, payload.MimeType, time.Now(), time.Now())

	entity, err := scanRowIntoMessageFile(row)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func scanRowIntoMessageFile(row *sql.Row) (*MessageFile, error) {
	var result MessageFile
	err := row.Scan(
		&result.ID,
		&result.MessageId,
		&result.SHA1,
		&result.MimeType,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
