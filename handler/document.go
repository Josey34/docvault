package handler

import (
	"docvault/dto"
	"docvault/usecase"
	"fmt"
	"io"
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

func (h *DocumentHandler) List(c *gin.Context) {
	docs, err := h.usecase.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var responses []*dto.DocumentResponse
	for _, doc := range docs {
		responses = append(responses, dto.FromEntity(doc))
	}

	c.JSON(http.StatusOK, responses)
}

func (h *DocumentHandler) GetMetadata(c *gin.Context) {
	id := c.Param("id")

	doc, err := h.usecase.GetMetadata(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := dto.FromEntity(doc)

	c.JSON(http.StatusOK, response)
}

func (h *DocumentHandler) Download(c *gin.Context) {
	id := c.Param("id")

	doc, err := h.usecase.GetMetadata(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	fileStream, err := h.usecase.Download(c.Request.Context(), doc.FileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer fileStream.Close()

	c.Header("Content-Type", doc.ContentType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", doc.FileName))
	c.Status(http.StatusOK)
	io.Copy(c.Writer, fileStream)
}

func (h *DocumentHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "document deleted"})
}
