package main

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

// add-source form data: form for adding a new RSS source.
// two inputs (name, url) with a focus index.
type addSourceModel struct {
	inputs     []textinput.Model
	focusIndex int
	err        string
}

func addSourceModelInit() addSourceModel {
	name := textinput.New()
	name.Placeholder = "Nome do feed"
	name.Focus()
	url := textinput.New()
	url.Placeholder = "https://..."
	return addSourceModel{inputs: []textinput.Model{name, url}}
}

func (m addSourceModel) Update(msg tea.KeyPressMsg) (addSourceModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		return m, func() tea.Msg { return backToMenuMsg{} }
	case "tab", "shift+tab", "up", "down":
		if msg.String() == "up" || msg.String() == "shift+tab" {
			m.focusIndex--
		} else {
			m.focusIndex++
		}
		m.focusIndex = (m.focusIndex + len(m.inputs)) % len(m.inputs)
		for i := range m.inputs {
			if i == m.focusIndex {
				m.inputs[i].Focus()
			} else {
				m.inputs[i].Blur()
			}
		}
		return m, nil
	case "enter":
		name := strings.TrimSpace(m.inputs[0].Value())
		url := strings.TrimSpace(m.inputs[1].Value())
		if name == "" || url == "" {
			m.err = "nome e url são obrigatórios"
			return m, nil
		}
		feed, err := RssFeedInit(name, url)
		if err != nil {
			m.err = err.Error()
			return m, nil
		}
		if err := AddFeedConfig(name, url); err != nil {
			m.err = err.Error()
			return m, nil
		}
		added := *feed
		return m, func() tea.Msg { return sourceAddedMsg{feed: added} }
	}

	var cmd tea.Cmd
	i := m.focusIndex
	m.inputs[i], cmd = m.inputs[i].Update(msg)
	return m, cmd
}

func (m addSourceModel) View() string {
	var buffer strings.Builder

	buffer.WriteString("\n")
	buffer.WriteString("Novo feed RSS\n")
	buffer.WriteString("-------------\n\n")

	fmt.Fprintf(&buffer, "Nome: %s\n", m.inputs[0].View())
	fmt.Fprintf(&buffer, "URL:  %s\n", m.inputs[1].View())

	if m.err != "" {
		fmt.Fprintf(&buffer, "\n%s\n", m.err)
	}

	buffer.WriteString("\n")
	buffer.WriteString("tab: alternar campo · enter: salvar · esc: cancelar\n")

	return buffer.String()
}
