package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/ifaisalabid1/file-upload-service/internal/domain"
	"github.com/ifaisalabid1/file-upload-service/internal/repository"
)

var v = validator.New()

type CreateFileRequest struct {
	Filename    string `json:"filename"     validate:"required,min=1,max=255"`
	ContentType string `json:"content_type" validate:"required"`
	SizeBytes   int64  `json:"size_bytes"   validate:"required,min=1"`
}

type uploadHandler struct {
	repo   repository.FileRepository
	logger *slog.Logger
}

func NewUploadHandler(repo repository.FileRepository, logger *slog.Logger) *uploadHandler {
	return &uploadHandler{repo: repo, logger: logger}
}

func (h *uploadHandler) CreateFile(w http.ResponseWriter, r *http.Request) {
	var req CreateFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid json"}`, http.StatusBadRequest)
		return
	}

	if err := v.Struct(req); err != nil {
		http.Error(w, `{"error": "validation failed"}`, http.StatusUnprocessableEntity)
		return
	}

	now := time.Now()
	file := &domain.File{
		ID:          uuid.New(),
		Filename:    req.Filename,
		ContentType: req.ContentType,
		SizeBytes:   req.SizeBytes,
		Status:      domain.StatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := h.repo.CreateFile(r.Context(), file); err != nil {
		h.logger.Error("create file db error", slog.String("error", err.Error()))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	h.logger.Info("file record created",
		slog.String("file_id", file.ID.String()),
		slog.String("filename", file.Filename),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	data, err := json.Marshal(file)
	if err != nil {
		h.logger.Error("json marshal error", slog.String("error", err.Error()))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInsufficientStorage)
		return
	}

	w.Write(data)
}
