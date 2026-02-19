package handler

import (
	"docvault/dto"
	"docvault/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DocumentHandler struct {
	usecase *usecase.DocumentUsecase
}

func NewDocumentHandler(usecase *usecase.DocumentUsecase) *DocumentHandler {
	return &DocumentHandler{usecase: usecase}
}

func (h *DocumentHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from request"})
		return
	}

	fileReader, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}
	defer fileReader.Close()

	doc, err := h.usecase.Upload(
		c.Request.Context(),
		file.Filename,
		file.Size,
		file.Header.Get("Content-Type"),
		fileReader,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.FromEntity(doc)
	c.JSON(http.StatusCreated, response)

}
