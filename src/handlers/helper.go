package handlers

import "go.mau.fi/whatsmeow/types"

func (h *Handler) SendTypingIndicator(jid types.JID) error {
	return h.Client.SendChatPresence(jid, types.ChatPresenceComposing, types.ChatPresenceMediaText)
}
