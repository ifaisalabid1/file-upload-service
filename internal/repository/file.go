package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ifaisalabid1/file-upload-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FileRepository interface {
	CreateFile(ctx context.Context, file *domain.File) error
	GetFileByID(ctx context.Context, id uuid.UUID) (*domain.File, error)
}

type postgresFileRepository struct {
	db *pgxpool.Pool
}

func NewPostgresFileRepository(db *pgxpool.Pool) FileRepository {
	return &postgresFileRepository{db: db}
}

func (r *postgresFileRepository) CreateFile(ctx context.Context, file *domain.File) error {
	query := `
			INSERT INTO files (id, filename, content_type, size_bytes, status)
			VALUES ($1, $2, $3, $4, $5);
	`
	_, err := r.db.Exec(ctx, query,
		file.ID,
		file.Filename,
		file.ContentType,
		file.SizeBytes,
		file.Status,
	)

	if err != nil {
		return fmt.Errorf("insert file: %w", err)
	}
	return nil
}

func (r *postgresFileRepository) GetFileByID(ctx context.Context, id uuid.UUID) (*domain.File, error) {
	query := `
			SELECT id, filename, content_type, size_bytes, status,
			created_at, updated_at FROM files
			WHERE id = $1;
	`

	var file domain.File
	err := r.db.QueryRow(ctx, query, id).Scan(
		&file.ID,
		&file.Filename,
		&file.ContentType,
		&file.SizeBytes,
		&file.Status,
		&file.CreatedAt,
		&file.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("not found")
		}
		return nil, fmt.Errorf("query file: %w", err)
	}
	return &file, nil
}
