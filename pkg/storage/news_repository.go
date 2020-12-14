package storage

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"

	"test-document/pkg/docs"
)

type docsRepo struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) docs.Repository {
	return &docsRepo{
		db: db,
	}
}

func (r *docsRepo) CreateDocuments(ctx context.Context, documents []docs.Document) error {
	g, ctx := errgroup.WithContext(ctx)

	tx := r.db.MustBegin()
	for _, document := range documents {
		document := document
		g.Go(func() error {
			err := r.CreateDocument(ctx, tx, document)
			if err != nil {
				return err
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		if errRl := tx.Rollback(); errRl != nil {
			return fmt.Errorf("wgWait: %w, rollback: %v", err, errRl)
		}
		return fmt.Errorf("%w", err)
	}

	if err := tx.Commit(); err != nil {
		if errRl := tx.Rollback(); errRl != nil {
			return fmt.Errorf("commit %w, rollback: %v", err, errRl)
		}
		return fmt.Errorf("create documents: %w", err)
	}

	return nil
}

func (r *docsRepo) CreateDocument(ctx context.Context, tx *sqlx.Tx, document docs.Document) error {
	const query = `INSERT INTO documents (name, date, number, sum) VALUES ($1, $2, $3, $4)`

	rows, err := tx.QueryContext(ctx, query, document.Name, document.Date, document.Number, document.Sum)
	if err != nil {
		return fmt.Errorf("create document: %w", err)
	}
	defer rows.Close()

	return nil
}

func (r *docsRepo) GetDocuments(ctx context.Context) ([]docs.Document, error) {
	const query = `SELECT name, date, number, sum FROM documents`

	var documents []docs.Document

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("not found: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var document docs.Document
		if err := rows.Scan(&document.Name, &document.Date, &document.Number, &document.Sum); err != nil {
			return nil, fmt.Errorf("not found: %w", err)
		}

		documents = append(documents, document)
	}

	return documents, nil
}