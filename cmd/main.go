package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Jake4-CX/portfolio-website-v2-backend/cmd/http/controllers"
	middlewares "github.com/Jake4-CX/portfolio-website-v2-backend/cmd/http/middleware"
	"github.com/Jake4-CX/portfolio-website-v2-backend/pkg/initializers"
	"github.com/Jake4-CX/portfolio-website-v2-backend/pkg/structs"
	"github.com/gin-gonic/gin"
)

func main() {

	initializers.LoadEnvVariables()
	initializers.InitializeDB()
	initializers.InitializeS3()

	router := gin.Default()

	router.Use(GinMiddleware(("*")))

	router.POST("/auth/register", controllers.CreateUser)
	router.POST("/auth/login", controllers.LoginUser)
	router.POST("/auth/validate", controllers.ValidateUserAccessToken)
	router.POST("/auth/refresh", controllers.RefreshAccessToken)

	// Technologies
	router.GET("/technologies", controllers.GetTechnologies)
	router.GET("/technologies/:technologyID", controllers.GetTechnology)

	// Projects
	router.GET("/projects", controllers.GetProjects)
	router.GET("/projects/:projectID", controllers.GetProject)

	// GitHub
	router.GET("/github/commits", controllers.GetCommitHistory)

	// Contact
	router.POST("/contact", controllers.ContactEmail)

	authorized := router.Group("/")

	authorized.Use(middlewares.RoleMiddleware(structs.ADMIN))
	{
		// Technologies
		authorized.POST("/technologies", controllers.CreateTechnology)
		authorized.PUT("/technologies/:technologyID", controllers.UpdateTechnology)
		authorized.DELETE("/technologies/:technologyID", controllers.DeleteTechnology)

		// Projects
		authorized.POST("/projects", controllers.CreateProject)
		authorized.PUT("/projects/:projectID", controllers.UpdateProject)
		authorized.PUT("/projects/:projectID/images", controllers.AssignProjectImages)
		authorized.DELETE("/projects/:projectID", controllers.DeleteProject)

		// Storage
		authorized.POST("/storage/create-presigned-url", controllers.CreatePresignedURL)
	}

	log.Fatal(router.Run("0.0.0.0:" + os.Getenv("REST_PORT")))

}

func GinMiddleware(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Request.Header.Del("Origin")

		c.Next()
	}
}
