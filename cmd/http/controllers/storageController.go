package controllers

import (
	"net/http"
	"github.com/Jake4-CX/portfolio-website-v2-backend/pkg/structs"
	"github.com/Jake4-CX/portfolio-website-v2-backend/pkg/initializers"
	"github.com/Jake4-CX/portfolio-website-v2-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

func CreatePresignedURL(c *gin.Context) {
	var request struct {
		UploadCategory structs.UploadCategory `json:"uploadCategory" binding:"required"` // Use structs.UploadCategory
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UploadCategory is required"})
		return
	}

	url, err := utils.GeneratePresignedPost(initializers.S3Session, request.UploadCategory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating presigned URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}