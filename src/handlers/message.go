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

	if e.Info.Sender.IsBot() {
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
		if content == "0" {
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
			return
		}

		if strings.Contains(strings.ToLower(content), "hubungi") {
			err := h.SendTypingIndicator(e.Info.Sender.ToNonAD())
			if err != nil {
				log.Error(err)
			}

			responseContent := "Baik, saya akan hubungkan ke tim kami. Mohon menunggu..."
			_, err = h.Client.SendMessage(context.Background(), e.Info.Sender.ToNonAD(), &waE2E.Message{
				Conversation: &responseContent,
			})

			if err != nil {
				log.Println(err)

			}

			h.LastResponse[e.Info.Sender.ToNonAD()] = responseContent
			return
		}

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
				return
			}

			menuID, err := h.GetResponseMenuID(lastResponse)
			if err != nil {
				responseContent := "Maaf, terjadi kesalahan. Saya akan hubungkan ke tim kami."
				h.Client.SendMessage(context.Background(), e.Info.Sender.ToNonAD(), &waE2E.Message{
					Conversation: &responseContent,
				})
				log.Error(err)
				h.LastResponse[e.Info.Sender.ToNonAD()] = "HUBUNGI_TIM"
				return
			}

			selectedMenu, err := h.MenuService.FindOptionMenu(menuID, optionNumber)
			if err != nil {
				responseContent := "Maaf, terjadi kesalahan. Saya akan hubungkan ke tim kami."
				if errors.Is(err, gorm.ErrRecordNotFound) {
					responseContent = "Maaf, opsi tersebut tidak tersedia. Silahkan coba lagi atau menunggu jawaban dari tim kami."
				}

				h.Client.SendMessage(context.Background(), e.Info.Sender.ToNonAD(), &waE2E.Message{
					Conversation: &responseContent,
				})
				log.Error(err)
				h.LastResponse[e.Info.Sender.ToNonAD()] = "HUBUNGI_TIM"
				return
			}

			menu, err := h.MenuService.FindMenuBySlug(selectedMenu.SubMenu.Slug)
			if err != nil {
				responseContent := "Maaf, terjadi kesalahan. Saya akan hubungkan ke tim kami."
				if errors.Is(err, gorm.ErrRecordNotFound) {
					responseContent = "Maaf, opsi tersebut tidak tersedia. Silahkan coba lagi atau menunggu jawaban dari tim kami."
				}

				h.Client.SendMessage(context.Background(), e.Info.Sender.ToNonAD(), &waE2E.Message{
					Conversation: &responseContent,
				})
				log.Error(err)
				h.LastResponse[e.Info.Sender.ToNonAD()] = "HUBUNGI_TIM"
				return
			}

			responseContent := menu.String()
			h.Client.SendMessage(context.Background(), e.Info.Sender.ToNonAD(), &waE2E.Message{
				Conversation: &responseContent,
			})

			if selectedMenu.SubMenu.Slug == "alamat" {
				latitude := -7.8091363641595715
				longitude := 110.36941817367556
				caption := "BNNP DIY"
				h.Client.SendMessage(context.Background(), e.Info.Sender.ToNonAD(), &waE2E.Message{
					LocationMessage: &waE2E.LocationMessage{
						DegreesLatitude:  &latitude,
						DegreesLongitude: &longitude,
						Name:             &caption,
					},
				})
			}

			h.LastResponse[e.Info.Sender.ToNonAD()] = responseContent
		}
	}
}
