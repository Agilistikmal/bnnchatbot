package handler

import (
	"log"

	"go.mau.fi/whatsmeow/types/events"
)

func MessageEvent(e *events.Message) {
	content := e.Message.GetConversation()
	log.Println(content)
}
