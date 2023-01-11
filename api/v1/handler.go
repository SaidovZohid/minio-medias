package v1

import (
	"github.com/SaidovZohid/minio-medias/api/models"
	"github.com/SaidovZohid/minio-medias/config"
	"github.com/SaidovZohid/minio-medias/pkg/logging"
	"github.com/SaidovZohid/minio-medias/storage/minio"
)

type handlerV1 struct {
	cfg          *config.Config
	minioStorage minio.MinioStorageI
	logger       logging.Logger
}

type HandlerV1Options struct {
	Cfg          *config.Config
	MinioStorage minio.MinioStorageI
	Logger       logging.Logger
}

func New(options *HandlerV1Options) *handlerV1 {
	return &handlerV1{
		cfg:          options.Cfg,
		minioStorage: options.MinioStorage,
		logger:       options.Logger,
	}
}

func errorResponse(err error) *models.ErrorResponse {
	return &models.ErrorResponse{
		Error: err.Error(),
	}
}
