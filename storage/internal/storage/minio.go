// Package storage wraps the MinIO client for object persistence.
package storage

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/ManyLinesEditor/backend/storage/config"
)

// MinioStorage wraps the MinIO client and a fixed bucket name.
type MinioStorage struct {
	client *minio.Client
	bucket string
}

// NewMinioStorage connects to MinIO and ensures the target bucket exists.
func NewMinioStorage(cfg config.MinIOConfig) (*MinioStorage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}
	return &MinioStorage{client: client, bucket: cfg.Bucket}, nil
}

// Put uploads r to the bucket under key with the given content-type and size.
func (s *MinioStorage) Put(ctx context.Context, key, contentType string, size int64, r io.Reader) error {
	_, err := s.client.PutObject(ctx, s.bucket, key, r, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

// Get opens an object for reading; caller must close the returned object.
func (s *MinioStorage) Get(ctx context.Context, key string) (*minio.Object, error) {
	return s.client.GetObject(ctx, s.bucket, key, minio.GetObjectOptions{})
}

// Bucket returns the configured bucket name.
func (s *MinioStorage) Bucket() string { return s.bucket }
