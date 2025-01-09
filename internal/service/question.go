package service

import (
	"github.com/mikalai2006/kingwood-api/graph/model"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
)

type QuestionService struct {
	repo repository.Question
	Hub  *Hub
}

func NewQuestionService(repo repository.Question, hub *Hub) *QuestionService {
	return &QuestionService{repo: repo, Hub: hub}
}

func (s *QuestionService) FindQuestion(params *model.QuestionFilter) (domain.Response[model.Question], error) {
	return s.repo.FindQuestion(params)
}

func (s *QuestionService) CreateQuestion(userID string, tag *model.QuestionInput) (*model.Question, error) {
	result, err := s.repo.CreateQuestion(userID, tag)

	s.Hub.HandleMessage(domain.Message{Type: "message", Method: "CREATE", Sender: userID, Recipient: "user2", Content: result, ID: "room1", Service: "question"})

	return result, err
}

func (s *QuestionService) UpdateQuestion(id string, userID string, data *model.QuestionInput) (*model.Question, error) {
	result, err := s.repo.UpdateQuestion(id, userID, data)

	s.Hub.HandleMessage(domain.Message{Type: "message", Method: "PATCH", Sender: userID, Recipient: "user2", Content: result, ID: "room1", Service: "question"})

	return result, err
}

func (s *QuestionService) DeleteQuestion(id string, userID string) (model.Question, error) {
	result, err := s.repo.DeleteQuestion(id)

	s.Hub.HandleMessage(domain.Message{Type: "message", Method: "DELETE", Sender: userID, Recipient: "user2", Content: result, ID: "room1", Service: "question"})

	return result, err
}
