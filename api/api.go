package api

import (
	_ "github.com/SaidovZohid/minio-medias/api/docs"
	v1 "github.com/SaidovZohid/minio-medias/api/v1"
	"github.com/SaidovZohid/minio-medias/config"
	"github.com/SaidovZohid/minio-medias/pkg/logging"
	"github.com/SaidovZohid/minio-medias/storage/minio"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

type RouterOptions struct {
	Cfg          *config.Config
	MinioStorage minio.MinioStorageI
	Logger       logging.Logger
}

// @title           Swagger for blog api
// @version         1.0
// @description     This is upload and get products api.
// @BasePath  /v1
func New(opt *RouterOptions) *gin.Engine {
	router := gin.Default()

	handlerV1 := v1.New(&v1.HandlerV1Options{
		Cfg:          opt.Cfg,
		MinioStorage: opt.MinioStorage,
		Logger:       opt.Logger,
	})

	// router.Static("/files", "http://172.17.0.2:9090/browser/98b87c92-91d8-11ed-898d-14cb19858780/")

	apiV1 := router.Group("/v1")

	apiV1.POST("/create-file", handlerV1.CreateFile)
	// apiV1.GET("/get-file/:fileId", handlerV1.GetFile)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
