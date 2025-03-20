package models

import (
	"fmt"

	"github.com/agilistikmal/bnnchat/src/lib"
)

type Menu struct {
	ID      int           `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Slug    string        `gorm:"unique" json:"slug,omitempty"`
	Header  string        `json:"header,omitempty"`
	Content string        `json:"content,omitempty"`
	Options []*MenuOption `gorm:"foreignKey:MenuID;references:ID" json:"options,omitempty"`
	Footer  string        `json:"footer,omitempty"`
}

type MenuOption struct {
	ID        int   `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	MenuID    int   `json:"menu_id,omitempty"`
	SubMenuID int   `json:"sub_menu_id,omitempty"`
	SubMenu   *Menu `json:"sub_menu,omitempty"`
	Position  int   `json:"position,omitempty"`
}

func (m *Menu) String() string {
	optionsText := ""
	if len(m.Options) > 0 {
		optionsText = "\n*Menu*\n"
		for _, option := range m.Options {
			optionsText += fmt.Sprintf("%d) %s\n", option.Position, option.SubMenu.Header)
		}
	}

	if m.Content != "" {
		m.Content = "\n" + m.Content
	}

	if m.Footer != "" {
		m.Footer = "\n" + m.Footer
	}

	return fmt.Sprintf("%s\n%s\n%s\n%s //// %s", m.Header, m.Content, optionsText, m.Footer, lib.EncodeBase62(m.ID))
}
