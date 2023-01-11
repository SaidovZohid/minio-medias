package minio

import (
	"context"
	"fmt"
	"io"

	"github.com/SaidovZohid/minio-medias/config"
	"github.com/SaidovZohid/minio-medias/pkg/logging"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Object struct {
	ID   string
	Size int64
	Tags map[string]string
}

type Client struct {
	logger      logging.Logger
	minioClient *minio.Client
}

func NewClient(cfg *config.Config, logger logging.Logger) (*Client, error) {
	minioClient, err := minio.New(cfg.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.AccessKey, cfg.Minio.SecretKey, ""),
		Secure: false,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create minio client. err: %w", err)
	}

	return &Client{
		logger:      logger,
		minioClient: minioClient,
	}, nil
}

func (c *Client) GetFile(ctx context.Context, bucketName, fileId string) (*minio.Object, error) {
	obj, err := c.minioClient.GetObject(context.Background(), bucketName, fileId, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file with id: %s from minio bucket %s. err: %w", fileId, bucketName, err)
	}

	return obj, nil
}

func (c *Client) GetBucketFiles(ctx context.Context, bucketName string) ([]*minio.Object, error) {
	var files []*minio.Object
	for lobj := range c.minioClient.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{}) {
		if lobj.Err != nil {
			c.logger.Errorf("failed to list object from minio bucket %s. err: %v", bucketName, lobj.Err)
			continue
		}

		object, err := c.minioClient.GetObject(context.Background(), bucketName, lobj.Key, minio.GetObjectOptions{})
		if err != nil {
			c.logger.Errorf("failed to get object key=%s from minio bucket %s. err: %v", lobj.Key, bucketName, lobj.Err)
			continue
		}
		files = append(files, object)
	}

	return files, nil
}

func (c *Client) UploadFile(ctx context.Context, fileId, fileName, bucketName string, fileSize int64, reader io.Reader) error {
	exists, errBucketExists := c.minioClient.BucketExists(ctx, bucketName)
	if errBucketExists != nil || !exists {
		c.logger.Errorf("no bucket %s, creating one...\n", bucketName)
		err := c.minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create new bucket. err: %w", err)
		}
	}

	c.logger.Debugf("put new object %s to bucket %s\n", fileName, bucketName)
	_, err := c.minioClient.PutObject(context.Background(), bucketName, fileId, reader, fileSize,
		minio.PutObjectOptions{
			UserMetadata: map[string]string{
				"Name": fileName,
			},
			ContentType: "application/octet-stream",
		},
	)
	if err != nil {
		return fmt.Errorf("failed to upload file. err: %w", err)
	}

	return nil
}

func (c *Client) DeleteFile(ctx context.Context, noteUUID, filename string) error {
	err := c.minioClient.RemoveObject(context.Background(), noteUUID, filename, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file. err: %w", err)
	}

	return nil 
}