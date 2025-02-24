package handlers

import (
	"github.com/agilistikmal/bnnchat/src/services"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

type Handler struct {
	Client       *whatsmeow.Client
	LastResponse map[types.JID]string
	MenuService  *services.MenuService
}

func NewHandler(client *whatsmeow.Client, menuService *services.MenuService) *Handler {
	return &Handler{
		Client:       client,
		LastResponse: make(map[types.JID]string),
		MenuService:  menuService,
	}
}
