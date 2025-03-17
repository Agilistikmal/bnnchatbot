package handlers

import (
	"github.com/agilistikmal/bnnchat/src/services"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"gorm.io/gorm"
)

type Handler struct {
	Client       *whatsmeow.Client
	LastResponse map[types.JID]string
	DB           *gorm.DB
	MenuService  *services.MenuService
}

func NewHandler(client *whatsmeow.Client, db *gorm.DB, menuService *services.MenuService) *Handler {
	return &Handler{
		Client:       client,
		LastResponse: make(map[types.JID]string),
		DB:           db,
		MenuService:  menuService,
	}
}
