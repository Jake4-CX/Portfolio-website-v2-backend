package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Jake4-CX/portfolio-website-v2-backend/cmd/http/controllers"
	"github.com/Jake4-CX/portfolio-website-v2-backend/pkg/initializers"
	"github.com/gin-gonic/gin"
)

func main() {

	initializers.LoadEnvVariables()
	initializers.InitializeDB()

	router := gin.Default()

	router.Use(GinMiddleware(("*")))

	router.POST("/auth/register", controllers.CreateUser)
	router.POST("/auth/login", controllers.LoginUser)
	router.POST("/auth/validate", controllers.ValidateUserAccessToken)
	router.POST("/auth/refresh", controllers.RefreshAccessToken)

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