package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/ManyLinesEditor/backend/storage/internal/middleware"
	"github.com/ManyLinesEditor/backend/storage/internal/models"
	"github.com/ManyLinesEditor/backend/storage/internal/services"
)

// FileHandler provides HTTP handlers for file upload and download.
type FileHandler struct {
	files *services.FileService
}

// NewFileHandler creates a FileHandler with the provided FileService.
func NewFileHandler(files *services.FileService) *FileHandler {
	return &FileHandler{files: files}
}

// Upload handles file upload.
//
// @Summary     Upload file
// @Description Uploads a file via multipart/form-data
// @Tags        files
// @Accept      multipart/form-data
// @Produce     json
// @Param       file formData file true "File to upload"
// @Success     201 {object} models.UploadResult
// @Failure     400 {object} models.ErrorResponse
// @Failure     500 {object} models.ErrorResponse
// @Security    BearerAuth
// @Router      /files [post]
func (h *FileHandler) Upload(c *gin.Context) {
	ownerID := middleware.OwnerID(c)

	fh, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "field 'file' is required"})
		return
	}

	result, err := h.files.Upload(c.Request.Context(), ownerID, fh)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// Download streams a file to the client.
//
// @Summary     Download file
// @Description Downloads a file by its UUID. Returns the file as an attachment.
// @Tags        files
// @Produce     application/octet-stream
// @Param       id path string true "File UUID" format(uuid)
// @Success     200 {file}   string
// @Failure     400 {object} models.ErrorResponse "invalid file id"
// @Failure     404 {object} models.ErrorResponse "file not found"
// @Failure     500 {object} models.ErrorResponse
// @Security    BearerAuth
// @Router      /files/{id} [get]
func (h *FileHandler) Download(c *gin.Context) {
	fileID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid file id"})
		return
	}

	result, err := h.files.Download(c.Request.Context(), fileID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "file not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	defer func() {
		if err := result.Reader.Close(); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		}
	}()

	c.Header("Content-Disposition",
		fmt.Sprintf(`attachment; filename="%s"`, result.Name))
	c.DataFromReader(http.StatusOK, result.Size, result.ContentType, result.Reader, nil)
}
