package controllers

import (
	"os"
	"regexp"

	"github.com/Jake4-CX/portfolio-website-v2-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

func ContactEmail(c *gin.Context) {
	var request struct {
		Name      string `json:"name" binding:"required"`
		Email     string `json:"email" binding:"required"`
		Message   string `json:"message" binding:"required"`
		Recaptcha string `json:"recaptcha" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Verify recaptcha
	if !utils.VerifyRecaptcha(request.Recaptcha) {
		c.JSON(400, gin.H{"error": "Recaptcha verification failed"})
		return
	}

	// Verify that the provided email is valid (regex)
	if !verifyEmail(request.Email) {
		c.JSON(400, gin.H{"error": "Invalid email address"})
		return
	}

	// Send email (tell the user that the email was sent)
	utils.SendEmail(
		request.Email,
		"Portfolio Contact Form - Email Sent",
		"Hello "+request.Name+",<br><br>"+
			"Thank you for reaching out to me. I will get back to you as soon as possible.<br><br>"+
			"<strong>Your Message:</strong><br>"+request.Message+"<br><br>"+
			"Best Regards,<br>Jack",
	)

	// Send email (send the email to me)
	utils.SendEmail(
		os.Getenv("EMAIL_CONTACT"),
		"Portfolio Contact Form - New Message",
		"Name: "+request.Name+"<br>"+
			"Email: "+request.Email+"<br>"+
			"Message: "+request.Message,
	)

	c.JSON(200, gin.H{"message": "Email sent"})
}

func verifyEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}
