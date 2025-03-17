package service

import (
	"fmt"
	"os"

	"github.com/mikalai2006/kingwood-api/internal/config"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
)

type ArchiveMessageService struct {
	repo        repository.ArchiveMessage
	Hub         *Hub
	imageConfig config.IImageConfig
	Services    *Services
}

func NewArchiveMessageService(repo repository.ArchiveMessage, Hub *Hub, imageConfig config.IImageConfig) *ArchiveMessageService {
	return &ArchiveMessageService{repo: repo, Hub: Hub, imageConfig: imageConfig}
}

func (s *ArchiveMessageService) CreateArchiveMessage(userID string, data *domain.Message) (*domain.ArchiveMessage, error) {
	result, err := s.repo.CreateArchiveMessage(userID, data)

	return result, err
}

func (s *ArchiveMessageService) FindArchiveMessage(params *domain.ArchiveMessageFilter) (domain.Response[domain.ArchiveMessage], error) {
	return s.repo.FindArchiveMessage(params)
}

func (s *ArchiveMessageService) DeleteArchiveMessage(id string) (domain.ArchiveMessage, error) {
	result, err := s.repo.DeleteArchiveMessage(id)

	// Delete images for message.
	for i := range result.Images {
		objImage := result.Images[i]
		pathDir := fmt.Sprintf("public/%s", objImage.Service)

		path := fmt.Sprintf("%s/%s/%s%s", pathDir, objImage.ServiceID, objImage.Path, objImage.Ext)
		os.Remove(path)

		for j := range s.imageConfig.Sizes {
			path := fmt.Sprintf("%s/%s/%s-%s%s", pathDir, objImage.ServiceID, s.imageConfig.Sizes[j].Prefix, objImage.Path, objImage.Ext)
			os.Remove(path)
		}
	}

	return result, err
}
