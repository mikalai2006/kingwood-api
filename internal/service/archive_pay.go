package service

import (
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
)

type ArchivePayService struct {
	repo     repository.ArchivePay
	Hub      *Hub
	Services *Services
}

func NewArchivePayService(repo repository.ArchivePay, hub *Hub) *ArchivePayService {
	return &ArchivePayService{repo: repo, Hub: hub}
}

func (s *ArchivePayService) FindArchivePay(input *domain.ArchivePayFilter) (domain.Response[domain.ArchivePay], error) {
	return s.repo.FindArchivePay(input)
}

func (s *ArchivePayService) CreateArchivePay(userID string, data *domain.Pay) (*domain.ArchivePay, error) {
	var result *domain.ArchivePay

	result, err := s.repo.CreateArchivePay(userID, data)
	if err != nil {
		return nil, err
	}

	// s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "CREATE", Sender: userID, Recipient: result.WorkerId.Hex(), Content: result, ID: "room1", Service: "ArchivePay"})

	return result, err
}

func (s *ArchivePayService) DeleteArchivePay(id string, userID string) (*domain.ArchivePay, error) {
	var result *domain.ArchivePay
	prevResults, err := s.repo.FindArchivePay(&domain.ArchivePayFilter{ID: []string{id}})
	if err != nil {
		return nil, err
	}
	if len(prevResults.Data) > 0 {
		// prevResult := prevResults.Data[0]

		result, err = s.repo.DeleteArchivePay(id, userID)
		if err != nil {
			return nil, err
		}
		// s.Hub.HandleMessage(domain.MessageSocket{Type: "message", Method: "DELETE", Sender: userID, Recipient: "", Content: result, ID: "room1", Service: "ArchivePay"})
	}

	return result, err
}
