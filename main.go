package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
)

// TODO: fix bug where the cursor in viewStateFeed
// can only browse the first two items of the array

// Current rendered view
type viewState int

const (
	viewStateMenu viewState = iota
	viewStateFeed
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

type rootModel struct {
	state     viewState
	menuModel menuModel
	feedModel feedModel
}

func rootModelInit(feeds []RssFeed) rootModel {
	return rootModel{
		state:     viewStateMenu,
		menuModel: menuModelInit(feeds),
		feedModel: feedModel{},
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
		switch m.state {
		case viewStateMenu:
			return m.updateMenu(msg)
		case viewStateFeed:
			return m.updateFeed(msg)
		}
	}

	return m, nil
}

func (m rootModel) updateMenu(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit
	case "n":
		// TODO: add new feed - input: nome e endpoint
	case "up", "k":
		if m.menuModel.cursor > 0 {
			m.menuModel.cursor--
		}
	case "down", "j":
		if m.menuModel.cursor < len(m.menuModel.feeds)-1 {
			m.menuModel.cursor++
		}
	case "space", "enter", "e":
		if len(m.menuModel.feeds) == 0 {
			break
		}
		selected := m.menuModel.feeds[m.menuModel.cursor]
		m.feedModel = feedModelInit(selected.Name, selected.Items)

		m.state = viewStateFeed
	}

	return m, nil
}

func (m rootModel) updateFeed(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc", "backspace":
		m.state = viewStateMenu
	case "up", "k":
		if m.feedModel.cursor > 0 {
			m.feedModel.cursor--
			prev_page_treshold := m.feedModel.items_per_page*m.feedModel.current_page - m.feedModel.items_per_page
			if m.feedModel.cursor < prev_page_treshold {
				m.feedModel.current_page--
			}
		}
	case "down", "j":
		if m.feedModel.cursor < len(m.feedModel.items)-1 {
			m.feedModel.cursor++
			next_page_treshold := m.feedModel.items_per_page*m.feedModel.current_page - 1
			if m.feedModel.cursor > next_page_treshold {
				m.feedModel.current_page++
			}
		}
	}
	return m, nil
}

func (m rootModel) View() tea.View {
	switch m.state {
	case viewStateFeed:
		return m.feedViewRender()
	default:
		return m.menuViewRender()
	}
}

func (m rootModel) menuViewRender() tea.View {
	var buffer strings.Builder

	buffer.WriteString("\n")
	buffer.WriteString("Feeds\n")
	buffer.WriteString("-----\n")

	if len(m.menuModel.feeds) == 0 {
		buffer.WriteString("Nenhum feed encontrado\n")
	}

	for i, feed := range m.menuModel.feeds {
		cursor := " "
		if m.menuModel.cursor == i {
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

	return tea.NewView(buffer.String())
}

func (m rootModel) feedViewRender() tea.View {
	var buffer strings.Builder

	buffer.WriteString("\n")
	fmt.Fprintf(&buffer, "%s\n", m.feedModel.feedName)
	buffer.WriteString(strings.Repeat("-", len(m.feedModel.feedName)))
	buffer.WriteString("\n")

	if len(m.feedModel.items) == 0 {
		buffer.WriteString("Nenhum item encontrado\n")
	}

	current_page := m.feedModel.current_page // começa em 1
	items_per_page := m.feedModel.items_per_page

	// i começa no início da página atual
	i := current_page*items_per_page - items_per_page
	for range items_per_page {
		item := m.feedModel.items[i]
		if m.feedModel.cursor == i {
			buffer.WriteString("--")
		}
		buffer.WriteString("\n")
		fmt.Fprintf(&buffer, "%s\n", item.Title)
		fmt.Fprintf(&buffer, "%s %s\n", item.PubDate, item.Author)
		if m.feedModel.cursor == i {
			buffer.WriteString("--")
		}
		buffer.WriteString("\n")
		i++
	}

	buffer.WriteString("\n")
	buffer.WriteString("esc|q: voltar\n")

	return tea.NewView(buffer.String())
}

func main() {
	// TODO: salvar feeds em um arquivo (cache), botão para refetch.
	// caso o arquivo não exista, faz o fetch automaticamente.
	// podemos checar a data de criação do arquivo. caso o arquivo
	// não tenha sido criado hoje, fazemos o refetch automaticamente

	// TODO: feature - adicionar novo feed através do menu

	var feeds []RssFeed

	folhaMundo, err := RssFeedInit(
		"folha mundo",
		"https://feeds.folha.uol.com.br/mundo/rss091.xml")

	if err != nil {
		fmt.Println(err)
	}

	feeds = append(feeds, *folhaMundo)

	ui := tea.NewProgram(rootModelInit(feeds))
	if _, err := ui.Run(); err != nil {
		fmt.Printf("Failed to run the bubbletea program: %s", err)
		os.Exit(1)
	}
}
