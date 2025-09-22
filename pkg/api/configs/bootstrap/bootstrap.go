package bootstrap

import (
	"fmt"
	"log"
	"os"

	"github.com/juanrojas09/core_domain/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	if os.Getenv("DATABASE_MIGRATE") == "true" {
		err := db.AutoMigrate(
			&domain.Users{},
			&domain.AccountLedger{},
			&domain.Categories{},
			&domain.Activities{},
			&domain.ConversationAccesses{},
			&domain.ConversationsMessages{},
			&domain.Conversations{},
			&domain.PaymentStatus{},
			&domain.PaymentsRefunds{},
			&domain.Payments{},
			&domain.Roles{},
			&domain.ServiceAcceptance{},
			&domain.ServiceEvidence{},
			&domain.ServiceReviews{},
			&domain.ServicesRequests{},
			&domain.Status{},
			&domain.UserPreferences{},
		)
		if err != nil {

			log.Panic("Migration failed:", err)
		}
	}

	if os.Getenv("DATABASE_DEBUG") == "true" {
		db = db.Debug()
	}

	return db
}

func InitLogger() *log.Logger {
	return log.New(os.Stdout, "[LOG]", log.LstdFlags|log.Lshortfile)
}
