package docs

import (
	"context"
)

type Service interface {
	CreateDocuments(ctx context.Context, documents []Document) error
	GetDocuments(ctx context.Context) ([]Document, error)
}

type service struct {
	newsRepo Repository
}

func NewService(newsRepo Repository) Service {
	return &service{
		newsRepo: newsRepo,
	}
}

func (s *service) CreateDocuments(ctx context.Context, documents []Document) error {
	return s.newsRepo.CreateDocuments(ctx, documents)
}

func (s *service) GetDocuments(ctx context.Context) ([]Document, error) {
	return s.newsRepo.GetDocuments(ctx)
}