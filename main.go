package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

// Current rendered view
type viewState int

const (
	viewStateMenu viewState = iota
	viewStateFeed
	viewStateAddSource
)

// Transition messages: views emit these via a tea.Cmd and rootModel
// handles them to switch the active view.
type openAddSourceMsg struct{}
type openFeedMsg struct {
	name  string
	items []RssItem
}
type sourceAddedMsg struct{ feed RssFeed }
type backToMenuMsg struct{}

type rootModel struct {
	state     viewState
	menuModel menuModel
	feedModel feedModel
	addSource addSourceModel
}

func rootModelInit(feeds []RssFeed) rootModel {
	return rootModel{
		state:     viewStateMenu,
		menuModel: menuModelInit(feeds),
	}
}

func (m rootModel) Init() tea.Cmd {
	return nil
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		var cmd tea.Cmd
		switch m.state {
		case viewStateMenu:
			m.menuModel, cmd = m.menuModel.Update(msg)
		case viewStateFeed:
			m.feedModel, cmd = m.feedModel.Update(msg)
		case viewStateAddSource:
			m.addSource, cmd = m.addSource.Update(msg)
		}
		return m, cmd
	case openAddSourceMsg:
		m.addSource = addSourceModelInit()
		m.state = viewStateAddSource
	case openFeedMsg:
		m.feedModel = feedModelInit(msg.name, msg.items)
		m.state = viewStateFeed
	case sourceAddedMsg:
		m.menuModel.feeds = append(m.menuModel.feeds, msg.feed)
		m.state = viewStateMenu
	case backToMenuMsg:
		m.state = viewStateMenu
	}

	return m, nil
}

func (m rootModel) View() tea.View {
	switch m.state {
	case viewStateFeed:
		return tea.NewView(m.feedModel.View())
	case viewStateAddSource:
		return tea.NewView(m.addSource.View())
	default:
		return tea.NewView(m.menuModel.View())
	}
}

func main() {
	// TODO: botão para refetch. checar a data de criação do arquivo de
	// cache; caso não tenha sido criado hoje, fazer o refetch automaticamente.

	configs, err := LoadFeedConfigs()
	if err != nil {
		fmt.Println(err)
	}

	var feeds []RssFeed
	for _, cfg := range configs {
		feed, err := RssFeedInit(cfg.Name, cfg.Endpoint)
		if err != nil {
			fmt.Println(err)
			continue
		}
		feeds = append(feeds, *feed)
	}

	ui := tea.NewProgram(rootModelInit(feeds))
	if _, err := ui.Run(); err != nil {
		fmt.Printf("Failed to run the bubbletea program: %s", err)
		os.Exit(1)
	}
}
