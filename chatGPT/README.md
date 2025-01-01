# ChatGPT Bubble

<img width="1200" src="./assets/demo.gif" />

A ChatGPT interface Bubble using the Bubble Tea framework. This bubble provides an interactive chat window with markdown rendering support, and real-time streaming responses from OpenAI's GPT models.

## Features
- Real-time streaming of ChatGPT responses
- Markdown rendering for better readability
- Token usage tracking (`hardcoded`)
- Keyboard-based navigation
- Support for GPT-3.5-turbo model

## Known Bugs
- Follow-up prompts don't execute on a new line in the chat window
- Token tracking display not updating in real-time
	- `hardcoded` in the status bar below - still trying to figure out how the token limits work ...

# Simple Usage in Go

## Keybindings
- `Enter` - Send message/prompt to ChatGPT
- `Esc` - Toggle input focus
- `Ctrl+c` - Quit application
- Mouse wheel - Scroll through chat history

The component also supports markdown rendering of responses including:
- Code blocks with syntax highlighting
- Lists and tables
- Headers and emphasis

```go
package main

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
```