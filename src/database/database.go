package database

import (
	"log"

	"github.com/agilistikmal/bnnchat/src/models"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase() *gorm.DB {
	db, err := gorm.Open(postgres.Open(viper.GetString("postgres.dsn")))
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.Question{}, &models.Menu{}, &models.MenuOption{})

	return db
}
