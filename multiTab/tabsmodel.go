package multitab

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TabModel represents the interface for a tab component in the UI.
// It defines the necessary methods that a tab should implement to
// initialize, update, render its view, and provide its title.
type TabModel interface {
	// Init initializes the tab and returns an initial command.
	Init() tea.Cmd

	// Update handles incoming messages and updates the tab's state.
	// It returns the updated model and an optional command.
	Update(tea.Msg) (tea.Model, tea.Cmd)

	// View renders the tab's view as a string.
	View() string

	// Title returns the title of the tab.
	Title() string

	// SetActive handles the activation state of the tab.
	SetActive(bool)
}

// TabsModel represents the model for managing a collection of tabs in a UI component.
// It includes the list of tabs, the index of the currently active tab, and a flag to
// indicate whether tab navigation should wrap around when reaching the end of the list.
//
// Fields:
// - tabs: A slice of TabModel implementations representing the individual tabs.
// - active: An integer representing the index of the currently active tab.
// - wraparound: A boolean indicating whether tab navigation should wrap around.
//
// Styles:
// - tabStyle: The style applied to inactive tabs.
// - activeTabStyle: The style applied to the active tab.
type TabsModel struct {
	tabs       []TabModel // A slice of TabModel representing the individual tabs
	active     int        // An integer representing the index of the currently active tab
	wraparound bool       // A boolean indicating whether tab navigation should wrap around

	tabStyle       lipgloss.Style // The style applied to inactive tabs
	activeTabStyle lipgloss.Style // The style applied to the active tab
}

// returns a new instance of TabsModel with the provided tabs and wraparound flag.
func NewTabsModel(tabs []TabModel, wraparound bool) *TabsModel {
	return &TabsModel{
		tabs:       tabs,
		active:     0,
		wraparound: wraparound,

		activeTabStyle: lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(lipgloss.Color("#000")).
			Background(lipgloss.Color("#FFF")).
			Bold(true),
		tabStyle: lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(lipgloss.Color("#888")),
	}
}

func (m *TabsModel) Init() tea.Cmd {
	// Initialize all tab in list
	for _, tab := range m.tabs {
		tab.Init()
	}

	return nil
}

func (m *TabsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left":
			if m.active > 0 {
				m.tabs[m.active].SetActive(false) // Deactivate current
				m.active--
				m.tabs[m.active].SetActive(true) // Activate new
			} else if m.wraparound {
				m.tabs[m.active].SetActive(false)
				m.active = len(m.tabs) - 1
				m.tabs[m.active].SetActive(true)
			}
		case "right":
			if m.active < len(m.tabs)-1 {
				m.tabs[m.active].SetActive(false)
				m.active++
				m.tabs[m.active].SetActive(true)
			} else if m.wraparound {
				m.tabs[m.active].SetActive(false)
				m.active = 0
				m.tabs[m.active].SetActive(true)
			}
		}
	}

	// Update components
	_, cmd = m.tabs[m.active].Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *TabsModel) View() string {
	var b strings.Builder

	// Print the tabs bar & highlight the active tab
	for i, tab := range m.tabs {
		if i == m.active {
			fmt.Fprintf(&b, "%s ", m.activeTabStyle.Render(tab.Title()))
		} else {
			fmt.Fprintf(&b, "%s ", m.tabStyle.Render(tab.Title()))
		}
	}

	// display the active tab's view
	fmt.Fprintf(&b, "\n\n%s\n", m.tabs[m.active].View())
	return b.String()
}
