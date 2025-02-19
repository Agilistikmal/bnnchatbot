package handlers

import (
	"context"
	"errors"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"gorm.io/gorm"
)

func (h *Handler) MessageEvent(event any) {
	e, isOk := event.(*events.Message)
	if !isOk {
		return
	}

	content := e.Message.GetConversation()
	lastResponse := h.LastResponse[e.Info.Sender.ToNonAD()]
	log.Info(content)

	if lastResponse == "" {
		err := h.SendTypingIndicator(e.Info.Sender.ToNonAD())
		if err != nil {
			log.Error(err)
		}

		menu, err := h.MenuService.FindMenuBySlug("welcome")
		if err != nil {
			log.Fatal(err)
		}

		responseContent := menu.String()
		_, err = h.Client.SendMessage(context.Background(), e.Info.Sender.ToNonAD(), &waE2E.Message{
			Conversation: &responseContent,
		})
		if err != nil {
			log.Println(err)
		}

		h.LastResponse[e.Info.Sender.ToNonAD()] = responseContent
	} else {
		if strings.Contains(lastResponse, "*Menu*") {
			err := h.SendTypingIndicator(e.Info.Sender.ToNonAD())
			if err != nil {
				log.Error(err)
			}

			optionNumber, err := strconv.Atoi(content)
			if err != nil {
				responseContent := "Maaf, opsi tersebut tidak tersedia. Silahkan coba lagi atau menunggu jawaban dari tim kami."
				h.Client.SendMessage(context.Background(), e.Info.Sender.ToNonAD(), &waE2E.Message{
					Conversation: &responseContent,
				})
			}

			menuID, err := h.GetResponseMenuID(lastResponse)
			if err != nil {
				responseContent := "Maaf, terjadi kesalahan. Saya akan hubungkan ke tim kami."
				h.Client.SendMessage(context.Background(), e.Info.Sender.ToNonAD(), &waE2E.Message{
					Conversation: &responseContent,
				})
				log.Error(err)
			}

			selectedMenu, err := h.MenuService.FindOptionMenu(menuID, optionNumber)
			if err != nil {
				responseContent := "Maaf, terjadi kesalahan. Saya akan hubungkan ke tim kami."
				if errors.Is(err, gorm.ErrRecordNotFound) {
					responseContent = "Maaf, opsi tersebut tidak tersedia. Silahkan coba lagi atau menunggu jawaban dari tim kami."
				}

				latitude := -7.8093277
				longitude := 110.3666287
				h.Client.SendMessage(context.Background(), e.Info.Sender.ToNonAD(), &waE2E.Message{
					Conversation: &responseContent,
					LocationMessage: &waE2E.LocationMessage{
						DegreesLatitude:  &latitude,
						DegreesLongitude: &longitude,
					},
				})
				log.Error(err)
			}

			responseContent := selectedMenu.SubMenu.String()
			h.Client.SendMessage(context.Background(), e.Info.Sender.ToNonAD(), &waE2E.Message{
				Conversation: &responseContent,
			})

			h.LastResponse[e.Info.Sender.ToNonAD()] = responseContent
		}
	}
}
