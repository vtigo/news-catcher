package main

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
)

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧"}

// menu view data
type menuModel struct {
	feeds  []RssFeed
	cursor int

	fetching     map[string]bool // feed names currently being fetched
	errored      map[string]bool // feed names whose last fetch failed
	spinning     bool            // whether the spinner ticker is running
	spinnerFrame int
}

func menuModelInit(feeds []RssFeed) menuModel {
	return menuModel{
		feeds:    feeds,
		cursor:   0,
		fetching: map[string]bool{},
		errored:  map[string]bool{},
	}
}

// fetchFeedCmd fetches a single feed in its own goroutine (the runtime runs
// each Cmd concurrently) and reports back via feedFetchedMsg. It builds a fresh
// RssFeed and never touches the shared feeds slice, so parallel fetches are safe.
func fetchFeedCmd(name, endpoint string) tea.Cmd {
	return func() tea.Msg {
		f := &RssFeed{Name: name, Endpoint: endpoint}
		if err := f.Fetch(); err != nil {
			return feedFetchedMsg{name: name, err: err}
		}
		if err := f.Store(); err != nil {
			return feedFetchedMsg{name: name, err: err}
		}
		return feedFetchedMsg{name: name, items: f.Items}
	}
}

func spinnerTick() tea.Cmd {
	return tea.Tick(120*time.Millisecond, func(time.Time) tea.Msg {
		return spinnerTickMsg{}
	})
}

// startFetch marks a feed as fetching and returns its fetch command. It also
// returns a spinner tick command the first time the ticker needs to start.
func (m *menuModel) startFetch(f RssFeed) []tea.Cmd {
	if m.fetching[f.Name] {
		return nil
	}
	m.fetching[f.Name] = true
	delete(m.errored, f.Name)
	cmds := []tea.Cmd{fetchFeedCmd(f.Name, f.Endpoint)}
	if !m.spinning {
		m.spinning = true
		cmds = append(cmds, spinnerTick())
	}
	return cmds
}

// applyFetchResult records a completed fetch: on error it flags the feed and
// keeps existing items; otherwise it replaces the feed's items.
func (m menuModel) applyFetchResult(msg feedFetchedMsg) menuModel {
	delete(m.fetching, msg.name)
	if msg.err != nil {
		m.errored[msg.name] = true
		return m
	}
	for i := range m.feeds {
		if m.feeds[i].Name == msg.name {
			m.feeds[i].Items = msg.items
			m.feeds[i].Fetched = true
			break
		}
	}
	return m
}

// onSpinnerTick advances the animation, keeping the ticker alive only while
// fetches are pending.
func (m menuModel) onSpinnerTick() (menuModel, tea.Cmd) {
	if len(m.fetching) == 0 {
		m.spinning = false
		return m, nil
	}
	m.spinnerFrame++
	return m, spinnerTick()
}

func (m menuModel) Update(msg tea.KeyPressMsg) (menuModel, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit
	case "n":
		return m, func() tea.Msg { return openAddSourceMsg{} }
	case "r":
		if len(m.feeds) == 0 {
			break
		}
		return m, tea.Batch(m.startFetch(m.feeds[m.cursor])...)
	case "R":
		var cmds []tea.Cmd
		for _, f := range m.feeds {
			cmds = append(cmds, m.startFetch(f)...)
		}
		return m, tea.Batch(cmds...)
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

		var status string
		switch {
		case m.fetching[feed.Name]:
			status = fmt.Sprintf("%s atualizando…", spinnerFrames[m.spinnerFrame%len(spinnerFrames)])
		case m.errored[feed.Name]:
			status = "(erro)"
		default:
			fetchedMsg := "cached"
			if feed.Fetched {
				fetchedMsg = "fetched"
			}
			status = fmt.Sprintf("(%d) (%s)", len(feed.Items), fetchedMsg)
		}
		fmt.Fprintf(&buffer, "%s %s %s\n", cursor, feed.Name, status)
	}

	buffer.WriteString("\n")
	buffer.WriteString("enter: abrir feed\n")
	buffer.WriteString("n: novo feed\n")
	buffer.WriteString("r: atualizar feed\n")
	buffer.WriteString("R: atualizar todos\n")
	buffer.WriteString("q: sair\n")

	return buffer.String()
}
