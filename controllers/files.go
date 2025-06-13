package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"litstore/api/initializers"
	"litstore/api/models"
	"litstore/api/upload"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadImages(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Files are too large",
		})
		return
	}

	files := c.Request.MultipartForm.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "No files uploaded",
		})
		return
	}

	bucket := os.Getenv("R2_BUCKET_NAME")
	endpoint := os.Getenv("R2_URL")
	uploader := upload.NewR2Uploader(initializers.R2Client, bucket, endpoint)

	var uploadedImages []map[string]any

	for _, fileHeader := range files {
		src, err := fileHeader.Open()
		if err != nil {
			continue // skip file on error
		}

		content, err := io.ReadAll(src)
		src.Close()
		if err != nil {
			continue // skip on read error
		}

		// Compute SHA256 hash
		hash := sha256.Sum256(content)
		hashHex := hex.EncodeToString(hash[:])

		// Check DB for existing image
		var existingImage models.Image
		err = initializers.DB.Where("hash = ?", hashHex).First(&existingImage).Error
		if err == nil {
			// Image exists, reuse metadata
			uploadedImages = append(uploadedImages, map[string]any{
				"image_id":       existingImage.ID,
				"url":            existingImage.URL,
				"already_exists": true,
			})
			continue
		}

		// Generate filename and upload
		id := uuid.New()
		ext := filepath.Ext(fileHeader.Filename)
		filename := fmt.Sprintf("%s%s", id.String(), ext)

		url, err := uploader.Upload(c.Request.Context(), filename, fileHeader)
		if err != nil {
			continue // skip on upload error
		}

		// Store in DB
		image := models.Image{
			Hash:        hashHex,
			URL:         url,
			Size:        int64(len(content)),
			MimeType:    fileHeader.Header.Get("Content-Type"),
			Description: "", // optionally accept description from form or separate call
		}

		if err := initializers.DB.Create(&image).Error; err != nil {
			continue
		}

		uploadedImages = append(uploadedImages, map[string]any{
			"image_id":       image.ID,
			"url":            image.URL,
			"already_exists": false,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"uploaded": uploadedImages,
	})
}
