package models

import "fmt"

type Menu struct {
	ID      string       `gorm:"primaryKey" json:"id,omitempty"`
	Header  string       `json:"header,omitempty"`
	Content string       `json:"content,omitempty"`
	Options []MenuOption `json:"options,omitempty"`
	Footer  string       `json:"footer,omitempty"`
}

type MenuOption struct {
	ID       string `gorm:"primaryKey" json:"id,omitempty"`
	MenuID   string `json:"menu_id,omitempty"`
	Position int    `json:"position,omitempty"`
}

func (m *Menu) String() string {
	optionsText := ""
	for _, option := range m.Options {
		optionsText += fmt.Sprintf("%s\n", option.ID)
	}

	return fmt.Sprintf("%s\n%s\n%s\n%s", m.Header, optionsText, m.Content, m.Footer)
}
