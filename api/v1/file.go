package v1

import (
	"context"
	"mime/multipart"
	"net/http"

	"github.com/SaidovZohid/minio-medias/api/models"
	"github.com/SaidovZohid/minio-medias/storage/minio"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type File struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

// func (h *handlerV1) GetFile(ctx *gin.Context) {
// 	h.logger.Info("Get File")

// 	h.logger.Debug("get note_uuid from url")
// 	noteuuid := ctx.Query("note_uuid")
// 	if noteuuid != "" {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("note_uuid is required")))
// 		return
// 	}
// 	fileId := ctx.Param("fileId")
// 	f, err := h.minioStorage.GetFile(context.Background(), noteuuid, fileId)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, models.ResponseOK{
// 		Message: string(f.Bytes),
// 	})
// }

// @Router /create-file [post]
// @Summary File upload
// @Description File upload
// @Tags file-upload
// @Accept json
// @Produce json
// @Param file formData file true "File"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) CreateFile(ctx *gin.Context) {
	h.logger.Info("Create File")
	var file File

	err := ctx.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = ctx.ShouldBind(&file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	fileReader, err := file.File.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	dto := minio.CreateFileDTO{
		Name:   file.File.Filename,
		Size:   file.File.Size,
		Reader: fileReader,
	}

	dto.NormalizeName()
	f, err := minio.NewFile(dto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	id, err := uuid.NewUUID()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	err = h.minioStorage.CreateFile(context.Background(), id.String(), f)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseOK{
		Message: "success",
	})
}
