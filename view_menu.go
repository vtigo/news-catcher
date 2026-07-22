package main

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

// menu view data
type menuModel struct {
	feeds  []RssFeed
	cursor int
}

func menuModelInit(feeds []RssFeed) menuModel {
	return menuModel{
		feeds:  feeds,
		cursor: 0,
	}
}

func (m menuModel) Update(msg tea.KeyPressMsg) (menuModel, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit
	case "n":
		return m, func() tea.Msg { return openAddSourceMsg{} }
	case "f":
		// TODO: fetch all feeds
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.feeds)-1 {
			m.cursor++
		}
	case "space", "enter", "e":
		if len(m.feeds) == 0 {
			break
		}
		selected := m.feeds[m.cursor]
		return m, func() tea.Msg {
			return openFeedMsg{name: selected.Name, items: selected.Items}
		}
	}

	return m, nil
}

func (m menuModel) View() string {
	var buffer strings.Builder

	buffer.WriteString("\n")
	buffer.WriteString("Feeds\n")
	buffer.WriteString("-----\n")

	if len(m.feeds) == 0 {
		buffer.WriteString("Nenhum feed encontrado\n")
	}

	for i, feed := range m.feeds {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		fetchedMsg := "cached"
		if feed.Fetched {
			fetchedMsg = "fetched"
		}
		fmt.Fprintf(&buffer, "%s %s (%d) (%s)\n", cursor, feed.Name, len(feed.Items), fetchedMsg)
	}

	buffer.WriteString("\n")
	buffer.WriteString("enter: abrir feed\n")
	buffer.WriteString("n: novo feed\n")
	buffer.WriteString("q: sair\n")

	return buffer.String()
}
