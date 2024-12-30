package main //change accordingly

import (
	"fmt"
	"os"

	ui "github.com/nick-popovic/custom-bubbles/chatGPT"

	tea "github.com/charmbracelet/bubbletea"
)

// model is the top-level model for the application.
type model struct {
	gpt ui.ChatGPTdialogueWindow
}

func New() model {
	return model{
		gpt: ui.NewGPTWindow(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		}
	}

	// Update components
	_, cmd = m.gpt.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return m.gpt.View()
}

func main() {
	p := tea.NewProgram(
		New(),
		tea.WithAltScreen(), // Use alternate screen buffer
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
