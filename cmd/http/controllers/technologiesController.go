package controllers

import (
	"github.com/Jake4-CX/portfolio-website-v2-backend/pkg/initializers"
	"github.com/Jake4-CX/portfolio-website-v2-backend/pkg/structs"
	"github.com/Jake4-CX/portfolio-website-v2-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

func CreateTechnology(c *gin.Context) {

	var newTechnology struct {
		TechnologyName  string                 `json:"technologyName" binding:"required"`
		TechnologyType  structs.TechnologyType `json:"technologyType" binding:"required"`
		TechnologyImage string                 `json:"technologyImage" binding:"required"`
	}

	if err := c.ShouldBindJSON(&newTechnology); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Validate technologyImage URL
	if err := utils.ValidateS3URL(initializers.S3Session, newTechnology.TechnologyImage); err != nil {
		c.JSON(400, gin.H{"error": "Invalid technologyImage URL", "fullError": err.Error()})
		return
	}

	// Check if a technology already exists with the provided name
	var existingTechnology structs.Technologies
	result := initializers.DB.First(&existingTechnology, "technology_name = ?", newTechnology.TechnologyName)

	if result.Error == nil {
		c.JSON(400, gin.H{"error": "A technology already exists with this name"})
		return
	}

	technology := structs.Technologies{
		TechnologyName:  newTechnology.TechnologyName,
		TechnologyType:  newTechnology.TechnologyType,
		TechnologyImage: newTechnology.TechnologyImage,
	}

	result = initializers.DB.Create(&technology)

	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Error creating technology"})
		return
	}

	// Get value of the inserted technology

	c.JSON(200, gin.H{"message": "Technology created successfully", "technology": technology})
}

func UpdateTechnology(c *gin.Context) {
	
	technologyID := c.Param("technologyID")

	var updatedTechnology struct {
		TechnologyName  string                 `json:"technologyName" binding:"required"`
		TechnologyType  structs.TechnologyType `json:"technologyType" binding:"required"`
		TechnologyImage string                 `json:"technologyImage" binding:"required"`
	}

	if err := c.ShouldBindJSON(&updatedTechnology); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Check if a technology exists with the provided ID
	var existingTechnology structs.Technologies
	if err := initializers.DB.First(&existingTechnology, "id = ?", technologyID).Error; err != nil {
		c.JSON(400, gin.H{"error": "No technology found with this ID"})
		return
	}

	// Validate technologyImage URL
	if err := utils.ValidateS3URL(initializers.S3Session, updatedTechnology.TechnologyImage); err != nil {
		c.JSON(400, gin.H{"error": "Invalid technologyImage URL", "fullError": err.Error()})
		return
	}

	// Check if a technology already exists with the provided name (excluding the current technology with the provided ID)
	if err := initializers.DB.First(&existingTechnology, "technology_name = ? AND id != ?", updatedTechnology.TechnologyName, technologyID).Error; err == nil {
		c.JSON(400, gin.H{"error": "A technology already exists with this name"})
		return
	}

	technology := structs.Technologies{
		TechnologyName:  updatedTechnology.TechnologyName,
		TechnologyType:  updatedTechnology.TechnologyType,
		TechnologyImage: updatedTechnology.TechnologyImage,
	}
	// Update the technology with the provided ID
	if err := initializers.DB.Model(&existingTechnology).Updates(technology).Error; err != nil {
		c.JSON(500, gin.H{"error": "Error updating technology"})
		return
	}

	// Get value of the updated technology

	c.JSON(200, gin.H{"message": "Technology updated successfully", "technology": technology})
}

func DeleteTechnology(c *gin.Context) {

	technologyID := c.Param("technologyID")

	var technology structs.Technologies
	result := initializers.DB.First(&technology, "id = ?", technologyID)

	if result.Error != nil {
		c.JSON(400, gin.H{"error": "No technology found with this ID"})
		return
	}

	initializers.DB.Delete(&technology)

	c.JSON(200, gin.H{"message": "Technology deleted successfully"})
}

func GetTechnologies(c *gin.Context) {

	var technologies []structs.Technologies
	result := initializers.DB.Find(&technologies)

	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Error fetching technologies"})
		return
	}

	c.JSON(200, gin.H{"technologies": technologies})
}

func GetTechnology(c *gin.Context) {

	technologyID := c.Param("technologyID")

	var technology structs.Technologies
	result := initializers.DB.First(&technology, "id = ?", technologyID)

	if result.Error != nil {
		c.JSON(400, gin.H{"error": "No technology found with this ID"})
		return
	}

	c.JSON(200, gin.H{"technology": technology})
}