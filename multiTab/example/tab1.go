package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	multitab "github.com/nick-popovic/custom-bubbles/multiTab"
)

//////////////////////
/* TAB 1 Definition */
//////////////////////

type tab1 struct {
	title  string
	body   string
	active bool // New field
}

func NewTab1() multitab.TabModel {
	return &tab1{
		title: "Tab1",
		body:  "Content of Tab1",
	}
}

func (t *tab1) Init() tea.Cmd {
	return nil
}

func (t *tab1) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

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

func (t *tab1) View() string {
	return fmt.Sprintf("%s\n", t.body)
}

func (t *tab1) Title() string {
	return t.title
}

// Add this method to tab1, tab2, and tab3
func (t *tab1) SetActive(isActive bool) {
	t.active = isActive
}
