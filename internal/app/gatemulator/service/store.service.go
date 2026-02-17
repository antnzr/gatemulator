package service

import (
	"strings"

	"github.com/antnzr/gatemulator/internal/app/gatemulator/dto"
	"github.com/antnzr/gatemulator/internal/app/gatemulator/repository"
	"github.com/antnzr/gatemulator/internal/pkg/gatemulator/utils"
)

type storeService struct {
	dao repository.DAO
}

func NewStoreService(dao repository.DAO) StoreService {
	return &storeService{dao}
}

func (s *storeService) GetFile(sha1 string) (*dto.MessageFileResponse, interface{}) {
	entity := s.dao.MessageFile().GetBySha1(sha1)
	if entity == nil {
		return nil, nil
	}

	file := &dto.MessageFileResponse{
		ID:        entity.ID,
		MessageId: entity.MessageId.String(),
		SHA1:      entity.SHA1,
		MimeType:  entity.MimeType,
	}

	if strings.HasPrefix(file.MimeType, "image/") {
		content := utils.GenerateRandomImage(100, 100)
		return file, content
	}

	return file, nil
}
