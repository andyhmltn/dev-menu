package main

import (
	"flag"
	"fmt"
	"github.com/andyhmltn/dev-menu/tmux"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"strings"
)

type MenuListItem struct {
	id     string
	paneId string
	title  string
	desc   string
	cmd    string
}

type menuListItemFlag []MenuListItem

func (m *menuListItemFlag) String() string {
	var parts []string

	for _, item := range *m {
		parts = append(parts, fmt.Sprintf("%s:%s:%s:%s:%s", item.id, item.paneId, item.title, item.desc, item.cmd))
	}

	return strings.Join(parts, ",")
}

func (m *menuListItemFlag) GetById(id string) (MenuListItem, error) {
	for _, item := range *m {
		if item.id == id {
			return item, nil
		}
	}

	return MenuListItem{}, fmt.Errorf("Menu item not found with id %s", id)
}

func (m *menuListItemFlag) Set(value string) error {
	items := strings.Split(value, ",")

	for _, item := range items {
		parts := strings.SplitN(item, ":", 5)

		if len(parts) != 5 {
			return fmt.Errorf("Invalid command provided")
		}

		*m = append(*m, MenuListItem{id: parts[0],
			paneId: parts[1],
			title:  parts[2],
			desc:   parts[3],
			cmd:    parts[4],
		})
	}

	return nil
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc, cmd, paneId string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list      list.Model
	menuItems menuListItemFlag
	choices   []string
	cursor    int
	selected  map[int]struct{}
}

func selectMenuItem(id string, menuItems menuListItemFlag) {
	if id == "0" {
		for _, menuItem := range menuItems {
			if menuItem.id != "0" {
				tmux.RunCmdInTmuxPane(menuItem.cmd, menuItem.paneId)
			}
		}

	} else {

		menuItem, err := menuItems.GetById(id)

		if err != nil {
			panic("panic")
		}

		tmux.RunCmdInTmuxPane(menuItem.cmd, menuItem.paneId)
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

		if msg.String() == "enter" {
			selected := m.list.SelectedItem()

			value := strings.Split(selected.FilterValue(), ".")
			id := value[0]

			selectMenuItem(id, m.menuItems)
		}

		// Allow highlighting a menu item by typing the ID
		if msg.String() == "0" {
			m.list.Select(0)
			selectMenuItem("0", m.menuItems)
		}

		for key, menuItem := range m.menuItems {
			if msg.String() == menuItem.id {
				m.list.Select(key + 1)
				selectMenuItem(menuItem.id, m.menuItems)
			}

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
	var menuItems menuListItemFlag = menuListItemFlag{}

	flag.Var(&menuItems, "items", "Command separated list of id:title:command")
	flag.Parse()

	items := []list.Item{
		item{
			title: "0. Restart all",
			desc:  "Restart all services",
		},
		// item{
		// 	title: "G. Re-generate types",
		// 	desc:  "Blah blah",
		// 	cmd:   "XYZ",
		// },
	}

	for _, menuItem := range menuItems {
		title := fmt.Sprintf("%s. %s", menuItem.id, menuItem.title)
		items = append(items, item{title: title, desc: menuItem.desc})
	}

	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = lipgloss.NewStyle().Foreground(lipgloss.Color("#a089ff")).MarginLeft(2)
	d.Styles.SelectedDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("#a089ff")).MarginLeft(2)

	myList := list.New(items, d, 0, 0)

	m := model{list: myList, menuItems: menuItems}
	m.list.Title = "Dev Menu"

	app := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := app.Run(); err != nil {
		fmt.Printf("Encountered error: %v", err)
		os.Exit(1)
	}

}
