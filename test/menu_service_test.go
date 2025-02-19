package test

import (
	"log"
	"testing"

	"github.com/agilistikmal/bnnchat/src/config"
	"github.com/agilistikmal/bnnchat/src/database"
	"github.com/agilistikmal/bnnchat/src/services"
	"github.com/stretchr/testify/assert"
)

func TestFindMenus(t *testing.T) {
	config.NewConfig()
	db := database.NewDatabase()
	s := services.NewMenuService(db)

	menus, err := s.FindMenus()
	assert.Nil(t, err)

	assert.NotNil(t, menus)
	log.Println(menus[0].Options[0].Menu.Header)
}
