package minio

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/SaidovZohid/minio-medias/config"
	"github.com/SaidovZohid/minio-medias/pkg/logging"
	minioClient "github.com/SaidovZohid/minio-medias/pkg/minio"
)

type minioStorage struct {
	client *minioClient.Client
	logger logging.Logger
}

type MinioStorageI interface {
	GetFile(ctx context.Context, bucketname, filename string) (*File, error)
	GetFilesByNoteUUID(ctx context.Context, uuid string) ([]*File, error)
	CreateFile(ctx context.Context, noteUUID string, file *File) error
	DeleteFile(ctx context.Context, noteUUID, filename string) error
}

func NewStorage(cfg *config.Config, logger logging.Logger) (MinioStorageI, error) {
	client, err := minioClient.NewClient(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	return &minioStorage{
		client: client,
		logger: logger,
	}, nil
}

func (m *minioStorage) GetFile(ctx context.Context, bucketName, fileID string) (*File, error) {
	obj, err := m.client.GetFile(context.Background(), bucketName, fileID)
	if err != nil {
		return nil, fmt.Errorf("failed to get file, err: %w", err)
	}

	defer obj.Close()

	objectInfo, err := obj.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file, err: %w", err)
	}

	buffer := make([]byte, objectInfo.Size)
	_, err = obj.Read(buffer)

	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to get objects, err: %w", err)
	}

	f := File{
		ID:    objectInfo.Key,
		Name:  objectInfo.UserMetadata["Name"],
		Size:  objectInfo.Size,
		Bytes: buffer,
	}

	return &f, nil
}

func (m *minioStorage) GetFilesByNoteUUID(ctx context.Context, noteUUID string) ([]*File, error) {
	objects, err := m.client.GetBucketFiles(context.Background(), noteUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get objects. err: %w", err)
	}

	if len(objects) == 0 {
		return nil, fmt.Errorf("nothing in bucket, err: %w", err)
	}

	var files []*File
	for _, obj := range objects {
		objectInfo, err := obj.Stat()
		if err != nil {
			return nil, fmt.Errorf("failed to get file, err: %w", err)
		}

		buffer := make([]byte, objectInfo.Size)
		_, err = obj.Read(buffer)

		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("failed to get objects, err: %w", err)
		}

		f := File{
			ID:    objectInfo.Key,
			Name:  objectInfo.UserMetadata["Name"],
			Size:  objectInfo.Size,
			Bytes: buffer,
		}
		files = append(files, &f)
		obj.Close()
	}

	return files, nil
}

func (m *minioStorage) CreateFile(ctx context.Context, noteUUID string, file *File) error {
	err := m.client.UploadFile(context.Background(), file.ID, file.Name, noteUUID, file.Size, bytes.NewBuffer(file.Bytes))
	if err != nil {
		return err
	}

	return nil
}

func (m *minioStorage) DeleteFile(ctx context.Context, noteUUID, fileId string) error {
	err := m.client.DeleteFile(context.Background(), noteUUID, fileId)
	if err != nil {
		return err
	}

	return nil
}