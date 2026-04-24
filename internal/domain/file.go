package domain

import (
	"time"

	"github.com/google/uuid"
)

type FileStatus string

const (
	StatusPending   FileStatus = "pending"
	StatusUploading FileStatus = "uploading"
	StatusCompleted FileStatus = "completed"
	StatusFailed    FileStatus = "failed"
)

type File struct {
	ID          uuid.UUID
	Filename    string
	ContentType string
	SizeBytes   int64
	Status      FileStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
