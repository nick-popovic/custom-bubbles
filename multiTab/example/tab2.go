package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	multitab "github.com/nick-popovic/custom-bubbles/multiTab"
)

//////////////////////
/* TAB 2 Definition */
//////////////////////

type tab2 struct {
	title  string
	body   string
	active bool // New field
}

func NewTab2() multitab.TabModel {
	return &tab2{
		title: "Tab2",
		body:  "Content of Tab2",
	}
}

func (t *tab2) Init() tea.Cmd {
	return nil
}

func (t *tab2) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

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

func (t *tab2) View() string {
	return fmt.Sprintf("%s\n", t.body)
}

func (t *tab2) Title() string {
	return t.title
}

// Add this method to tab1, tab2, and tab3
func (t *tab2) SetActive(isActive bool) {
	t.active = isActive
}
