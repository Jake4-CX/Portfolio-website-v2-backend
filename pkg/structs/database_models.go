package structs

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormModel struct {
	ID        uint           `gorm:"primarykey;autoIncrement" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
} //@name models.gormModel

type Users struct {
	GormModel
	UserEmail    string   `json:"userEmail"`
	UserPassword string   `json:"userPassword"`
	UserRole     UserRole `json:"userRole" gorm:"default:USER"`
}

type RefreshTokens struct {
	GormModel
	UserId       uint      `json:"userId"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
}

type VerificationTokens struct {
	GormModel
	UserId           uint             `json:"userId"`
	VerificationUUID uuid.UUID        `json:"verificationUUID"`
	VerificationType VerificationType `json:"verificationType"`
	ExpiresAt        time.Time        `json:"expiresAt"`
}

type Projects struct {
	GormModel
	ProjectName         string                `json:"projectName"`
	ProjectDescription  string                `json:"projectDescription"`
	IsFeatured          bool                  `json:"isFeatured"`
	StartDate           time.Time             `json:"startDate"`
	EndDate             time.Time             `json:"endDate"`
	ProjectImages       []ProjectImages       `json:"projectImages" gorm:"foreignKey:ProjectId"`       // One-to-many relationship
	ProjectTechnologies []ProjectTechnologies `json:"projectTechnologies" gorm:"foreignKey:ProjectId"` // One-to-many relationship
	ProjectURLs         ProjectURLs           `json:"projectURLs" gorm:"foreignKey:ProjectId"`         // One-to-one relationship
}

type ProjectURLs struct {
	GormModel
	ProjectId  uint   `json:"projectId"`
	GitHubURL  string `json:"githubURL"`
	WebsiteURL string `json:"websiteURL"`
	YouTubeURL string `json:"youtubeURL"`
}

type ProjectImages struct {
	GormModel
	ProjectId uint   `json:"projectId"`
	ImageURL  string `json:"imageURL"`
}

type Technologies struct {
	GormModel
	TechnologyName  string         `json:"technologyName"`
	TechnologyType  TechnologyType `json:"technologyType"`
	TechnologyImage string         `json:"technologyImage"`
}

type ProjectTechnologies struct {
	GormModel
	ProjectId    uint `json:"projectId"`
	TechnologyId uint `json:"technologyId"`
}

type UploadCategory string
type TechnologyType string
type VerificationType string
type UserRole string

const (
	PROJECT_IMAGE    UploadCategory = "PROJECT_IMAGE"
	TECHNOLOGY_IMAGE UploadCategory = "TECHNOLOGY_IMAGE"
)

const (
	LANGUAGE  TechnologyType = "LANGUAGE"
	FRAMEWORK TechnologyType = "FRAMEWORK"
	DATABASE  TechnologyType = "DATABASE"
	TOOL      TechnologyType = "TOOL"
	OTHER     TechnologyType = "OTHER"
)

const (
	EMAIL_VERIFICATION VerificationType = "EMAIL_VERIFICATION"
	PHONE_VERIFICATION VerificationType = "PHONE_VERIFICATION"
	RESET_PASSWORD     VerificationType = "RESET_PASSWORD"
)

const (
	ADMIN UserRole = "ADMIN"
	USER  UserRole = "USER"
)
