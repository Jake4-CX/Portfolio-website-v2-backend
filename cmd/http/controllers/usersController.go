package controllers

import (
	"time"

	"github.com/Jake4-CX/portfolio-website-v2-backend/pkg/initializers"
	"github.com/Jake4-CX/portfolio-website-v2-backend/pkg/structs"
	"github.com/Jake4-CX/portfolio-website-v2-backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {

	var newUser struct {
		UserEmail    string `json:"userEmail" binding:"required"`
		UserPassword string `json:"userPassword" binding:"required"`
	}

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Check if an account already exists with the provided email
	var existingUser structs.Users
	result := initializers.DB.First(&existingUser, "user_email = ?", newUser.UserEmail)

	if result.Error == nil {
		c.JSON(400, gin.H{"error": "An account already exists with this email"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newUser.UserPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error hashing password"})
		return
	}

	user := structs.Users{
		UserEmail:    newUser.UserEmail,
		UserPassword: string(hash),
	}

	result = initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Error creating user"})
		return
	}

	c.JSON(200, gin.H{"message": "User created successfully"})

}

func LoginUser(c *gin.Context) {

	var login struct {
		UserEmail    string `json:"userEmail" binding:"required"`
		UserPassword string `json:"userPassword" binding:"required"`
	}

	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var user structs.Users
	result := initializers.DB.First(&user, "user_email = ?", login.UserEmail)

	if result.Error != nil {
		c.JSON(400, gin.H{"error": "No account found with this email"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(login.UserPassword))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid password"})
		return
	}

	// Genterate JWT token (accessToken & refreshToken)
	accessToken, refreshToken, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error generating tokens"})
		return
	}

	// Save the refresh token in the database
	err = utils.SaveRefreshToken(user.ID, refreshToken)

	if err != nil {
		c.JSON(500, gin.H{"error": "Error saving refresh token"})
		return
	}

	var loginResponse structs.LoginResponseModel

	loginResponse.User = user
	loginResponse.Token.AccessToken = accessToken
	loginResponse.Token.RefreshToken = refreshToken

	loginResponse.User.UserPassword = ""

	c.JSON(200, gin.H{"message": "Login successful", "data": loginResponse})

}

func ValidateUserAccessToken(c *gin.Context) {

	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(400, gin.H{"error": "Authorization header missing"})
		return
	}

	// Remove "Bearer " from the token string
	tokenString = tokenString[len("Bearer "):]

	// Check if the access token is valid
	userId, err := utils.ValidateToken(tokenString)

	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid access token"})
		return
	}

	var user structs.Users
	result := initializers.DB.First(&user, userId)

	if result.Error != nil {
		c.JSON(400, gin.H{"error": "No account found with this email"})
		return
	}

	user.UserPassword = ""

	c.JSON(200, gin.H{"message": "Access token is valid", "data": user})
}

func RefreshAccessToken(c *gin.Context) {

	var refreshToken struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}

	if err := c.ShouldBindJSON(&refreshToken); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Check if the refresh token is valid
	var refreshTokenRecord structs.RefreshTokens
	result := initializers.DB.First(&refreshTokenRecord, "refresh_token = ?", refreshToken.RefreshToken)

	if result.Error != nil {
		c.JSON(401, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Check if the refresh token is expired
	if time.Now().After(refreshTokenRecord.ExpiresAt) {
		c.JSON(400, gin.H{"error": "Refresh token has expired"})
		return
	}

	// Delete the refresh token from the database
	result = initializers.DB.Delete(&refreshTokenRecord)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Error deleting refresh token"})
		return
	}

	// Generate new access token
	accessToken, newRefreshToken, err := utils.GenerateToken(refreshTokenRecord.UserId)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error generating access token"})
		return
	}

	// Save the refresh token in the database
	err = utils.SaveRefreshToken(refreshTokenRecord.UserId, newRefreshToken)

	if err != nil {
		c.JSON(500, gin.H{"error": "Error saving refresh token"})
		return
	}

	var refreshResponse structs.TokensModel

	refreshResponse.AccessToken = accessToken
	refreshResponse.RefreshToken = newRefreshToken

	c.JSON(200, gin.H{"message": "Access token refreshed", "data": refreshResponse})

}
