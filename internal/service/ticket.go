package service

import (
	"github.com/mikalai2006/kingwood-api/graph/model"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
)

type TicketService struct {
	repo repository.Ticket
}

func NewTicketService(repo repository.Ticket) *TicketService {
	return &TicketService{repo: repo}
}

func (s *TicketService) FindTicket(params domain.RequestParams) (domain.Response[model.Ticket], error) {
	return s.repo.FindTicket(params)
}

func (s *TicketService) GetAllTicket(params domain.RequestParams) (domain.Response[model.Ticket], error) {
	return s.repo.GetAllTicket(params)
}

func (s *TicketService) CreateTicket(userID string, ticket *model.Ticket) (*model.Ticket, error) {
	return s.repo.CreateTicket(userID, ticket)
}

func (s *TicketService) CreateTicketMessage(userID string, message *model.TicketMessage) (*model.TicketMessage, error) {
	return s.repo.CreateTicketMessage(userID, message)
}

func (s *TicketService) DeleteTicket(id string) (model.Ticket, error) {
	return s.repo.DeleteTicket(id)
}
