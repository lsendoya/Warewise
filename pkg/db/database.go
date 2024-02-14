package db

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/lsendoya/Warewise/internal/product/domain"
	domain2 "github.com/lsendoya/Warewise/internal/user/domain"
	secret "github.com/lsendoya/Warewise/pkg/aws"
	"github.com/lsendoya/Warewise/pkg/logger"
	"github.com/lsendoya/Warewise/pkg/tools"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func NewDataBase(cfg aws.Config) (*gorm.DB, error) {
	dbCredentials := secret.GetSecretManager(cfg, os.Getenv("SECRET_NAME"))

	var rdsJSON SecretRDSJson
	if err := tools.UnmarshalJSON([]byte(dbCredentials), &rdsJSON); err != nil {
		logger.Fatalf("error reading secret credentials AWS RDS")
	}

	db, errConn := gorm.Open(postgres.Open(makeDSN(rdsJSON)), &gorm.Config{})
	if errConn != nil {
		logger.Errorf("gorm.Open() %v", errConn.Error())
		return nil, fmt.Errorf("gorm.Open(), Failed to connect to database, %w", errConn)
	}

	logger.Info("Connection opened to database")

	return db, nil
}

func makeDSN(s SecretRDSJson) string {
	dbName := os.Getenv("DB_NAME")
	port := 5432
	host := os.Getenv("DB_HOST")
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		host, port, s.Username, s.Password, dbName)

	return dsn
}

func MigrateDB(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	models := []interface{}{
		&domain.Product{},
		&domain2.User{},
	}

	for _, mdl := range models {
		if err := db.AutoMigrate(mdl); err != nil {
			log.Printf("Error migrating model %T: %v", mdl, err)
			return fmt.Errorf("error migrating model %T: %w", mdl, err)
		}
	}

	log.Println("Database migration completed successfully")
	return nil
}
