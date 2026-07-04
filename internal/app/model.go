package app

import (
	"fmt"
	"reqtea/internal/request"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type screen int

const (
	screenMenu screen = iota
	screenRequest
	screenCollections
	screenHistory
	screenSettings
)

type item string

func (i item) Title() string       { return string(i) }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return string(i) }

type Model struct {
	currentScreen screen
	menu          list.Model
	request       request.Model
}

func New() Model {
	items := []list.Item{
		item("New Request"),
		item("Collections"),
		item("History"),
		item("Settings"),
		item("Exit"),
	}

	menu := list.New(items, list.NewDefaultDelegate(), 0, 0)
	menu.Title = "ReqTea"
	menu.KeyMap.Quit.SetEnabled(false)

	return Model{
		currentScreen: screenMenu,
		menu:          menu,
		request:       request.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.menu.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "q":
			if m.currentScreen == screenMenu {
				return m, tea.Quit
			}

		case "esc":
			if m.currentScreen != screenMenu {
				m.currentScreen = screenMenu
				return m, nil
			}

		case "enter":
			if m.currentScreen == screenMenu {
				selected := m.menu.SelectedItem()
				if selected == nil {
					return m, nil
				}

				switch selected.(item) {
				case "New Request":
					m.currentScreen = screenRequest
				case "Collections":
					m.currentScreen = screenCollections
				case "History":
					m.currentScreen = screenHistory
				case "Settings":
					m.currentScreen = screenSettings
				case "Exit":
					return m, tea.Quit
				}

				return m, nil
			}
		}
	}

	switch m.currentScreen {
	case screenMenu:
		m.menu, cmd = m.menu.Update(msg)
		return m, cmd

	case screenRequest:
		m.request, cmd = m.request.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	switch m.currentScreen {
	case screenMenu:
		return m.menu.View()

	case screenRequest:
		return m.request.View()
		//return "New Request Screen\n\nPress esc to go back."

	case screenCollections:
		return "Collections Screen\n\nPress esc to go back."

	case screenHistory:
		return "History Screen\n\nPress esc to go back."

	case screenSettings:
		return "Settings Screen\n\nPress esc to go back."

	default:
		return fmt.Sprintf("Unknown screen: %v", m.currentScreen)
	}
}
