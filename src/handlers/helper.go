package handlers

import (
	"fmt"
	"strings"

	"github.com/agilistikmal/bnnchat/src/lib"
	"go.mau.fi/whatsmeow/types"
)

func (h *Handler) SendTypingIndicator(jid types.JID) error {
	return h.Client.SendChatPresence(jid, types.ChatPresenceComposing, types.ChatPresenceMediaText)
}

func (h *Handler) GetResponseMenuID(text string) (int, error) {
	splitted := strings.Split(text, "////")
	if len(splitted) < 1 {
		return 0, fmt.Errorf("ID Not found")
	}

	id := lib.DecodeBase62(strings.ReplaceAll(splitted[1], " ", ""))
	return id, nil
}
