package service

import (
	"github.com/mikalai2006/kingwood-api/internal/config"
	"github.com/mikalai2006/kingwood-api/internal/domain"
	"github.com/mikalai2006/kingwood-api/internal/repository"
)

type PostService struct {
	repo repository.Post
	i18n config.I18nConfig
}

func NewPostService(repo repository.Post, i18n config.I18nConfig) *PostService {
	return &PostService{repo: repo, i18n: i18n}
}

func (s *PostService) FindPost(params domain.RequestParams) (domain.Response[domain.Post], error) {
	return s.repo.FindPost(params)
}

func (s *PostService) GetAllPost(params domain.RequestParams) (domain.Response[domain.Post], error) {
	return s.repo.GetAllPost(params)
}

func (s *PostService) CreatePost(userID string, Post *domain.Post) (*domain.Post, error) {
	return s.repo.CreatePost(userID, Post)
}

func (s *PostService) UpdatePost(id string, userID string, Post *domain.PostInput) (*domain.Post, error) {
	return s.repo.UpdatePost(id, userID, Post)
}

func (s *PostService) DeletePost(id string) (domain.Post, error) {
	return s.repo.DeletePost(id)
}
