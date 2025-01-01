package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	multitab "github.com/nick-popovic/custom-bubbles/multiTab"
)

//////////////////////
/* TAB 2 Definition */
//////////////////////

type tab3 struct {
	title  string
	body   string
	active bool // New field
}

func NewTab3() multitab.TabModel {
	return &tab3{
		title: "Tab3",
		body:  "Content of Tab3",
	}
}

func (t *tab3) Init() tea.Cmd {
	return nil
}

func (t *tab3) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if t.active {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "i":
				t.body += " < "
			}
		}
	}

	return t, nil
}

func (t *tab3) View() string {
	return fmt.Sprintf("%s\n", t.body)
}

func (t *tab3) Title() string {
	return t.title
}

// Add this method to tab1, tab2, and tab3
func (t *tab3) SetActive(isActive bool) {
	t.active = isActive
}
