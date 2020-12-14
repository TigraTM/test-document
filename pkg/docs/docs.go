package docs

import (
	"context"
)

type Document struct {
	Name   string `json:"name"`
	Date   string `json:"date"`
	Number int    `json:"number"`
	Sum    string `json:"sum"`
}

type Repository interface {
	CreateDocuments(ctx context.Context, documents []Document) error
	GetDocuments(ctx context.Context) ([]Document, error)
}
