package handler

import (
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

type Handler struct {
	Client       *whatsmeow.Client
	LastResponse map[types.JID]string
}

func NewHandler(client *whatsmeow.Client) *Handler {
	return &Handler{
		Client:       client,
		LastResponse: make(map[types.JID]string),
	}
}
