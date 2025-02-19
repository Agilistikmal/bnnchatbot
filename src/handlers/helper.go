package handlers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"go.mau.fi/whatsmeow/types"
)

func (h *Handler) SendTypingIndicator(jid types.JID) error {
	return h.Client.SendChatPresence(jid, types.ChatPresenceComposing, types.ChatPresenceMediaText)
}

func (h *Handler) GetResponseMenuID(text string) (int, error) {
	splitted := strings.Split(text, "//")
	re := regexp.MustCompile("[0-9]+")
	result := re.FindAllString(splitted[0], -1)
	if len(result) < 1 {
		return 0, fmt.Errorf("ID Not found")
	}

	id, err := strconv.Atoi(result[0])
	return id, err
}
