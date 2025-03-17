package database

import (
	"errors"
	"log"

	"github.com/agilistikmal/bnnchat/src/models"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("data.db"))
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.Menu{}, &models.MenuOption{}, &models.Help{})

	var menu *models.Menu
	err = db.Take(&menu, "slug = ?", "welcome").Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			menu = &models.Menu{
				Slug:    "welcome",
				Header:  "Selamat datang di call center BNN",
				Content: "Silahkan pilih menu layanan dibawah ini",
				Footer:  `atau ketik "hubungi tim" untuk tersambung ke tim kami.`,
			}
			err := db.Create(&menu).Error
			if err != nil {
				logrus.Fatal("Failed to create default menu")
			}
		}
	}

	return db
}
