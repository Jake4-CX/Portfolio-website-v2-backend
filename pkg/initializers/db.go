package initializers

import (
	"os"

	"github.com/Jake4-CX/portfolio-website-v2-backend/pkg/structs"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitializeDB() {
	var err error
	DB, err = gorm.Open(mysql.Open(os.Getenv("DB_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database")
	}

	// Automigrate all models
	err = DB.AutoMigrate(
		&structs.Users{},
		&structs.RefreshTokens{},
		&structs.VerificationTokens{},
		&structs.Projects{},
		&structs.Technologies{},
		&structs.ProjectTechnologies{},
	)

	if err != nil {
		log.Fatal("Error automigrating models")
	}

	log.Info("Database connection established")
}