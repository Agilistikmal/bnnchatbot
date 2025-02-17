package handlers

import (
	"context"
	"log"
	"strings"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
)

func (h *Handler) MessageEvent(event any) {
	e, isOk := event.(*events.Message)
	if !isOk {
		return
	}

	content := e.Message.GetConversation()
	log.Println(content)

	if strings.Contains(content, "halo") {
		err := h.SendTypingIndicator(e.Info.Sender.ToNonAD())
		if err != nil {
			log.Println(err)
		}

		responseContent := "Halo"
		log.Println(h.LastResponse[e.Info.Sender.ToNonAD()], e.Info.Sender.ToNonAD())

		// Check Last Response
		if h.LastResponse[e.Info.Sender.ToNonAD()] == "Halo" {
			responseContent = "Saya sudah balas Halo"
		}

		_, err = h.Client.SendMessage(context.Background(), e.Info.Sender.ToNonAD(), &waE2E.Message{
			Conversation: &responseContent,
		})
		if err != nil {
			log.Println(err)
		}

		h.LastResponse[e.Info.Sender.ToNonAD()] = responseContent
	}
}
