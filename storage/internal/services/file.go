package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/google/uuid"

	"github.com/ManyLinesEditor/backend/storage/internal/models"
	"github.com/ManyLinesEditor/backend/storage/internal/repository"
	"github.com/ManyLinesEditor/backend/storage/internal/storage"
)

// FileService orchestrates file uploads and downloads.
type FileService struct {
	files   repository.FileRepo
	storage *storage.MinioStorage
}

// NewFileService wires a FileRepo and MinioStorage into a FileService.
func NewFileService(files repository.FileRepo, s *storage.MinioStorage) *FileService {
	return &FileService{files: files, storage: s}
}

func (s *FileService) Upload(ctx context.Context, ownerID uuid.UUID, fh *multipart.FileHeader) (res *models.UploadResult, err error) {
	src, err := fh.Open()
	if err != nil {
		return nil, fmt.Errorf("open upload: %w", err)
	}
	defer func() {
		if cerr := src.Close(); cerr != nil {
			err = errors.Join(err, cerr)
		}
	}()

	id := uuid.New()
	ct := fh.Header.Get("Content-Type")
	if ct == "" {
		ct = "application/octet-stream"
	}
	key := fmt.Sprintf("%s/%s", ownerID, id)

	if err = s.storage.Put(ctx, key, ct, fh.Size, src); err != nil {
		return nil, fmt.Errorf("minio put: %w", err)
	}

	meta := &models.File{
		ID:          id,
		OwnerID:     ownerID,
		Name:        fh.Filename,
		ContentType: ct,
		SizeBytes:   fh.Size,
		Bucket:      s.storage.Bucket(),
		ObjectKey:   key,
	}
	if err = s.files.Create(ctx, meta); err != nil {
		return nil, fmt.Errorf("save metadata: %w", err)
	}

	res = &models.UploadResult{ID: id, Name: fh.Filename}
	return res, err
}

// DownloadResult carries everything needed for the HTTP response.
// Caller is responsible for closing Reader.
type DownloadResult struct {
	Reader      io.ReadCloser
	ContentType string
	Size        int64
	Name        string
}

// Download fetches file metadata from Postgres and opens the MinIO object.
func (s *FileService) Download(ctx context.Context, fileID uuid.UUID) (*DownloadResult, error) {
	meta, err := s.files.FindByID(ctx, fileID)
	if err != nil {
		return nil, err
	}

	obj, err := s.storage.Get(ctx, meta.ObjectKey)
	if err != nil {
		return nil, fmt.Errorf("minio get: %w", err)
	}
	return &DownloadResult{
		Reader:      obj,
		ContentType: meta.ContentType,
		Size:        meta.SizeBytes,
		Name:        meta.Name,
	}, nil
}
