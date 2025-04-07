package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list     list.Model
	choices  []string
	cursor   int
	selected map[int]struct{}
}

func initialModel() model {
	return model{
		choices:  []string{"Restart all", "Restart Frontend", "Restart Backend"},
		selected: make(map[int]struct{}),
	}

}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if msg.String() == "0" {
			m.list.Select(0)
		}

		if msg.String() == "1" {
			m.list.Select(1)
		}

		if msg.String() == "2" {
			m.list.Select(2)
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func main() {
	items := []list.Item{
		item{title: "0. Restart All", desc: "Restart all services"},
		item{title: "1. Restart Frontend", desc: "Reruns: npm dev"},
		item{title: "2. Restart Backend", desc: "Reruns: docker compose up"},
		item{title: "A. Generate DB Schema", desc: "Runs: npm run generate"},
	}

	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = lipgloss.NewStyle().Foreground(lipgloss.Color("#a089ff")).MarginLeft(2)
	d.Styles.SelectedDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("#a089ff")).MarginLeft(2)

	myList := list.New(items, d, 0, 0)

	m := model{list: myList}
	m.list.Title = "Dev Menu"

	app := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := app.Run(); err != nil {
		fmt.Printf("oh no, an error! %v", err)
		os.Exit(1)
	}

}
