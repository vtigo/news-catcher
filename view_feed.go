package main

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

// feed view data
type feedModel struct {
	feedName       string
	items          []RssItem
	cursor         int
	items_per_page int
	current_page   int // começa em 1
}

func feedModelInit(name string, items []RssItem) feedModel {
	return feedModel{
		feedName:       name,
		items:          items,
		cursor:         0,
		items_per_page: 5,
		current_page:   1,
	}
}

func (m feedModel) Update(msg tea.KeyPressMsg) (feedModel, tea.Cmd) {
	switch msg.String() {
	case "q", "esc", "backspace":
		return m, func() tea.Msg { return backToMenuMsg{} }
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
			prev_page_treshold := m.items_per_page*m.current_page - m.items_per_page
			if m.cursor < prev_page_treshold {
				m.current_page--
			}
		}
	case "down", "j":
		if m.cursor < len(m.items)-1 {
			m.cursor++
			next_page_treshold := m.items_per_page*m.current_page - 1
			if m.cursor > next_page_treshold {
				m.current_page++
			}
		}
	}
	return m, nil
}

func (m feedModel) View() string {
	var buffer strings.Builder

	buffer.WriteString("\n")
	fmt.Fprintf(&buffer, "%s\n", m.feedName)
	buffer.WriteString(strings.Repeat("-", len(m.feedName)))
	buffer.WriteString("\n")

	if len(m.items) == 0 {
		buffer.WriteString("Nenhum item encontrado\n")
	}

	current_page := m.current_page // começa em 1
	items_per_page := m.items_per_page

	// i começa no início da página atual
	i := current_page*items_per_page - items_per_page
	for range items_per_page {
		item := m.items[i]
		if m.cursor == i {
			buffer.WriteString("--")
		}
		buffer.WriteString("\n")
		fmt.Fprintf(&buffer, "%s\n", item.Title)
		fmt.Fprintf(&buffer, "%s %s\n", item.PubDate, item.Author)
		if m.cursor == i {
			buffer.WriteString("--")
		}
		buffer.WriteString("\n")
		i++
	}

	buffer.WriteString("\n")
	buffer.WriteString("esc|q: voltar\n")

	return buffer.String()
}
