package controllers

import (
	"fmt"
	"litstore/api/initializers"
	"litstore/api/upload"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}

	bucket := os.Getenv("R2_BUCKET_NAME")
	endpoint := os.Getenv("R2_URL")

	uploader := upload.NewR2Uploader(initializers.R2Client, bucket, endpoint)

	id := uuid.New()

	ext := filepath.Ext(file.Filename)

	filename := fmt.Sprintf("%s%s", id.String(), ext)

	url, err := uploader.Upload(c.Request.Context(), filename, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"url":     url,
	})
}
