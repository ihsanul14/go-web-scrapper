package database

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-web-scrapper/entity/model"
)

var Db *gorm.DB

func Connect() (*gorm.DB, error) {
	dialector := os.Getenv("POSTGRES_DIALECTOR")
	Db, err := gorm.Open(postgres.Open(dialector), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err = Db.AutoMigrate(&model.Data{}); err != nil {
		return nil, err
	}

	return Db, err
}
