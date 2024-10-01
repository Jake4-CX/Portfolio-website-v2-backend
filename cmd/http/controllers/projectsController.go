package controllers

import (
	"time"

	"github.com/Jake4-CX/portfolio-website-v2-backend/pkg/initializers"
	"github.com/Jake4-CX/portfolio-website-v2-backend/pkg/structs"
	"github.com/Jake4-CX/portfolio-website-v2-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func CreateProject(c *gin.Context) {
	var newProject struct {
		ProjectName         string `json:"projectName" binding:"required"`
		ProjectDescription  string `json:"projectDescription" binding:"required"`
		IsFeatured          bool   `json:"isFeatured"`
		StartDate           int64  `json:"startDate" binding:"required"`
		EndDate             int64  `json:"endDate" binding:"required"`
		IsEnabled           bool   `json:"isEnabled"`
		ProjectTechnologies []uint `json:"projectTechnologies" binding:"required"`
		ProjectURLs         struct {
			GitHubURL  string `json:"githubURL"`
			WebsiteURL string `json:"websiteURL"`
			YouTubeURL string `json:"youtubeURL"`
		} `json:"projectURLs" binding:"required"`
	}

	if err := c.ShouldBindJSON(&newProject); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Turn the dates into time.Time objects
	startDate := time.Unix(newProject.StartDate/1000, 0)
	endDate := time.Unix(newProject.EndDate/1000, 0)

	// Check if the project already exists with the same name

	var existingProject structs.Projects
	result := initializers.DB.Where("project_name = ?", newProject.ProjectName).First(&existingProject)

	if result.Error == nil {
		c.JSON(400, gin.H{"error": "Project already exists with the same name"})
		return
	}

	// Create the new project
	project := structs.Projects{
		ProjectName:        newProject.ProjectName,
		ProjectDescription: newProject.ProjectDescription,
		IsFeatured:         newProject.IsFeatured,
		StartDate:          startDate,
		EndDate:            endDate,
		IsEnabled:          newProject.IsEnabled,
	}

	// Save the project to get the ID

	result = initializers.DB.Create(&project)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Error creating project"})
		return
	}

	// Create and associate the ProjectURLs with the project
	projectURLs := structs.ProjectURLs{
		ProjectId:  project.ID, // Link the ProjectURLs to the created project
		GitHubURL:  newProject.ProjectURLs.GitHubURL,
		WebsiteURL: newProject.ProjectURLs.WebsiteURL,
		YouTubeURL: newProject.ProjectURLs.YouTubeURL,
	}

	// Save the ProjectURLs to the database
	result = initializers.DB.Create(&projectURLs)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Error creating project URLs"})
		return
	}

	// Check if the technologies exist, and create the projectTechnologies
	for _, technologyID := range newProject.ProjectTechnologies {
		var technology structs.Technologies
		result = initializers.DB.Where("id = ?", technologyID).First(&technology)

		if result.Error != nil {
			c.JSON(400, gin.H{"error": "Technology does not exist"})
			return
		}

		projectTechnology := structs.ProjectTechnologies{
			ProjectId:    project.ID, // Use the generated project ID
			TechnologyId: technology.ID,
		}

		result = initializers.DB.Create(&projectTechnology)
		if result.Error != nil {
			c.JSON(500, gin.H{"error": "Error creating project technology"})
			return
		}
	}

	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Error creating project"})
		return
	}

	c.JSON(200, gin.H{"project": project})
}

func UpdateProject(c *gin.Context) {

	projectID := c.Param("projectID")

	var updatedProject struct {
		ProjectName         string `json:"projectName" binding:"required"`
		ProjectDescription  string `json:"projectDescription" binding:"required"`
		IsFeatured          bool   `json:"isFeatured"`
		StartDate           int64  `json:"startDate" binding:"required"`
		EndDate             int64  `json:"endDate" binding:"required"`
		IsEnabled           bool   `json:"isEnabled"`
		ProjectTechnologies []uint `json:"projectTechnologies" binding:"required"`
		ProjectURLs         struct {
			GitHubURL  string `json:"githubURL"`
			WebsiteURL string `json:"websiteURL"`
			YouTubeURL string `json:"youtubeURL"`
		} `json:"projectURLs" binding:"required"`
	}

	if err := c.ShouldBindJSON(&updatedProject); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Convert dates to time.Time objects
	startDate := time.Unix(updatedProject.StartDate/1000, 0)
	endDate := time.Unix(updatedProject.EndDate/1000, 0)

	// Start a transaction
	tx := initializers.DB.Begin()
	if tx.Error != nil {
		log.Error("Failed to start transaction: ", tx.Error)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	// Check if the project exists
	var project structs.Projects
	if err := tx.Where("id = ?", projectID).First(&project).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(400, gin.H{"error": "Project does not exist"})
		} else {
			log.Error("Error finding project: ", err)
			c.JSON(500, gin.H{"error": "Internal server error"})
		}
		tx.Rollback()
		return
	}

	// Remove all existing old projectTechnologies
	if err := tx.Where("project_id = ?", projectID).Delete(&structs.ProjectTechnologies{}).Error; err != nil {
		log.Error("Error deleting old existing projectTechnologies: ", err)
		c.JSON(500, gin.H{"error": "Error removing old project technologies"})
		tx.Rollback()
		return
	}

	// Recreate the projectTechnologies
	for _, technologyID := range updatedProject.ProjectTechnologies {
		var technology structs.Technologies
		if err := tx.Where("id = ?", technologyID).First(&technology).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(400, gin.H{"error": "Technology does not exist"})
			} else {
				log.Error("Error finding technology: ", err)
				c.JSON(500, gin.H{"error": "Internal server error"})
			}
			tx.Rollback()
			return
		}

		projectTechnology := structs.ProjectTechnologies{
			ProjectId:    project.ID,
			TechnologyId: technology.ID,
		}

		if err := tx.Create(&projectTechnology).Error; err != nil {
			log.Error("Error creating project technology: ", err)
			c.JSON(500, gin.H{"error": "Error creating project technology"})
			tx.Rollback()
			return
		}
	}

	// Update project details
	project.ProjectName = updatedProject.ProjectName
	project.ProjectDescription = updatedProject.ProjectDescription
	project.IsFeatured = updatedProject.IsFeatured
	project.StartDate = startDate
	project.EndDate = endDate
	project.IsEnabled = updatedProject.IsEnabled

	// Update ProjectURLs
	projectURLs := structs.ProjectURLs{
		ProjectId:  project.ID,
		GitHubURL:  updatedProject.ProjectURLs.GitHubURL,
		WebsiteURL: updatedProject.ProjectURLs.WebsiteURL,
		YouTubeURL: updatedProject.ProjectURLs.YouTubeURL,
	}

	// Upsert the ProjectURLs (update if exists, create if not)
	if err := tx.Where("project_id = ?", project.ID).Assign(projectURLs).FirstOrCreate(&projectURLs).Error; err != nil {
		log.Error("Error updating project URLs: ", err)
		c.JSON(500, gin.H{"error": "Error updating project URLs"})
		tx.Rollback()
		return
	}

	// Save the updated project
	if err := tx.Save(&project).Error; err != nil {
		log.Error("Error updating project: ", err)
		c.JSON(500, gin.H{"error": "Error updating project"})
		tx.Rollback()
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		log.Error("Error committing transaction: ", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(200, gin.H{"project": project})
}

func AssignProjectImages(c *gin.Context) {

	projectID := c.Param("projectID")

	var newProjectImages struct {
		ImageURLs []string `json:"imageURLs" binding:"required"`
	}

	if err := c.ShouldBindJSON(&newProjectImages); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Start a transaction
	tx := initializers.DB.Begin()
	if tx.Error != nil {
		log.Error("Failed to start transaction: ", tx.Error)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	// Check if the project exists
	var project structs.Projects
	if err := tx.Where("id = ?", projectID).First(&project).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(400, gin.H{"error": "Project does not exist"})
		} else {
			log.Error("Error finding project: ", err)
			c.JSON(500, gin.H{"error": "Internal server error"})
		}
		tx.Rollback()
		return
	}

	// Remove all existing images for the project
	if err := tx.Where("project_id = ?", projectID).Delete(&structs.ProjectImages{}).Error; err != nil {
		log.Error("Error deleting old project images: ", err)
		c.JSON(500, gin.H{"error": "Error removing old project images"})
		tx.Rollback()
		return
	}

	// Create the new images
	for _, imageURL := range newProjectImages.ImageURLs {
		// Validate technologyImage URL
		if err := utils.ValidateS3URL(initializers.S3Session, imageURL); err != nil {
			c.JSON(400, gin.H{"error": "Invalid imageURL", "fullError": err.Error()})
			tx.Rollback()
			return
		}

		projectImage := structs.ProjectImages{
			ProjectId: project.ID,
			ImageURL:  imageURL,
		}

		if err := tx.Create(&projectImage).Error; err != nil {
			log.Error("Error creating project image: ", err)
			c.JSON(500, gin.H{"error": "Error creating project image"})
			tx.Rollback()
			return
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		log.Error("Error committing transaction: ", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(200, gin.H{"message": "Project images assigned successfully"})
}

func DeleteProject(c *gin.Context) {

	projectID := c.Param("projectID")

	// Check if the project exists
	var project structs.Projects
	result := initializers.DB.Where("id = ?", projectID).First(&project)

	if result.Error != nil {
		c.JSON(400, gin.H{"error": "Project does not exist"})
		return
	}

	// Delete the project
	result = initializers.DB.Delete(&project)

	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Error deleting project"})
		return
	}

	c.JSON(200, gin.H{"message": "Project deleted successfully"})
}

func GetProjects(c *gin.Context) {

	var projects []structs.Projects
	result := initializers.DB.Preload("ProjectImages").Preload("ProjectTechnologies").Preload("ProjectURLs").Find(&projects)

	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Error retrieving projects"})
		return
	}

	c.JSON(200, gin.H{"projects": projects})
}

func GetProject(c *gin.Context) {

	projectID := c.Param("projectID")

	var project structs.Projects
	result := initializers.DB.Where("id = ?", projectID).Preload("ProjectImages").Preload("ProjectTechnologies").Preload("ProjectURLs").First(&project)

	if result.Error != nil {
		c.JSON(400, gin.H{"error": "Project does not exist"})
		return
	}

	c.JSON(200, gin.H{"project": project})
}
