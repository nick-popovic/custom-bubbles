package chatGPT

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/mistakenelf/teacup/statusbar"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/ssestream"
	"github.com/pkoukk/tiktoken-go"
)

// MaxTokens defines the maximum context window size for GPT-3.5-turbo
// This limit includes both input and output tokens.
var (
	MaxTokens = 4096 // GPT-3.5-turbo default context window
)

// countTokens calculates the number of tokens in the given text using the
// cl100k_base encoding used by GPT-3.5-turbo and GPT-4.
// Returns the token count and any error encountered during encoding.
func countTokens(text string) (int, error) {
	encoding := "cl100k_base" // encoding for GPT-3.5-turbo and GPT-4
	tkm, err := tiktoken.GetEncoding(encoding)
	if err != nil {
		return 0, err
	}
	tokens := tkm.Encode(text, nil, nil)

	return len(tokens), nil
}

// SessionStats tracks the conversation metrics including total tokens used,
// number of messages exchanged, and remaining tokens in the context window.
type SessionStats struct {
	totalTokens     int // Total tokens used in the conversation
	messageCount    int // Number of messages exchanged
	remainingTokens int // Remaining tokens in the context window
}

// update recalculates the session statistics after processing new tokens.
// It updates the total tokens used, increments message count, and
// adjusts the remaining tokens available in the context window.
func (s *SessionStats) update(newTokens int) {
	s.totalTokens += newTokens
	s.messageCount++
	s.remainingTokens = MaxTokens - s.totalTokens
}

// display prints the current session statistics to standard output,
// showing the number of messages exchanged, total tokens used,
// and remaining tokens in the context window.
func (s *SessionStats) display() {
	fmt.Printf("\n=== Session Stats ===\n")
	fmt.Printf("Messages: %d\n", s.messageCount)
	fmt.Printf("Tokens Used: %d\n", s.totalTokens)
	fmt.Printf("Tokens Remaining: %d\n", s.remainingTokens)
	fmt.Printf("Context Usage: %.1f%%\n", float64(s.totalTokens)/float64(MaxTokens)*100)
	fmt.Println("===================\n")
}

type execMsg struct {
	msg string
}

type ChatGPTdialogueWindow struct {
	client openai.Client
	ctx    context.Context

	// conversation history
	messages          []openai.ChatCompletionMessageParamUnion
	assistantResponse *strings.Builder // Change to pointer

	stream *ssestream.Stream[openai.ChatCompletionChunk]

	viewport  viewport.Model
	textInput textinput.Model
	statusbar statusbar.Model
}

func NewGPTWindow() ChatGPTdialogueWindow {
	client := openai.NewClient()
	ctx := context.Background()
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage("You are a helpful assistant."),
	}

	assistantResponse := &strings.Builder{} // Initialize as pointer

	ti := textinput.New()
	ti.Placeholder = "Type here..."
	ti.Focus()

	return ChatGPTdialogueWindow{
		client:            *client,
		ctx:               ctx,
		messages:          messages,
		assistantResponse: assistantResponse,
		viewport:          viewport.Model{},
		textInput:         ti,
		statusbar: statusbar.New(
			statusbar.ColorConfig{
				Foreground: lipgloss.AdaptiveColor{Dark: "#ffffff", Light: "#ffffff"},
				Background: lipgloss.AdaptiveColor{Light: "#F25D94", Dark: "#F25D94"},
			},
			statusbar.ColorConfig{
				Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
				Background: lipgloss.AdaptiveColor{Light: "#3c3836", Dark: "#3c3836"},
			},
			statusbar.ColorConfig{
				Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
				Background: lipgloss.AdaptiveColor{Light: "#A550DF", Dark: "#A550DF"},
			},
			statusbar.ColorConfig{
				Foreground: lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"},
				Background: lipgloss.AdaptiveColor{Light: "#6124DF", Dark: "#6124DF"},
			},
		),
	}
}

func (m *ChatGPTdialogueWindow) Init() tea.Cmd {
	m.textInput.Focus()
	return textinput.Blink
}

func (m *ChatGPTdialogueWindow) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	// Handle window resize events
	case tea.WindowSizeMsg:

		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 2

		m.statusbar.SetSize(msg.Width)
		m.statusbar.SetContent("Feature Idea:", "It would be nice if you had some status tokens here ...", "Tokens Used", "% Used")

	// catch the next() calls and receive the message
	case execMsg:
		m.assistantResponse.WriteString(msg.msg)

		// unrendered text
		//m.viewport.SetContent(m.assistantResponse.String())

		// rendered text
		// rendered text
		r, _ := glamour.NewTermRenderer(
			glamour.WithStandardStyle("dark"),
			glamour.WithWordWrap(m.viewport.Width),
		)
		out, _ := r.Render(m.assistantResponse.String() + "\n")
		m.viewport.SetContent(out)
		m.viewport.ViewDown() // scroll to the bottom as new messages arrive

		if m.stream.Next() {
			evt := m.stream.Current()
			if len(evt.Choices) > 0 {

				//kick off the compilation of chatgpt message stream
				return m, func() tea.Msg {
					return execMsg{
						msg: evt.Choices[0].Delta.Content,
					}
				}
			}

		} else {
			// Add assistant's response to history
			m.messages = append(m.messages, openai.AssistantMessage(m.assistantResponse.String()))

			// unrendered text
			// m.viewport.SetContent(m.assistantResponse.String() + "\n")

			// rendered text
			r, _ := glamour.NewTermRenderer(
				glamour.WithStandardStyle("dark"),
				glamour.WithWordWrap(m.viewport.Width),
			)
			out, _ := r.Render(m.assistantResponse.String() + "\n")
			m.viewport.SetContent(out)
			m.viewport.ViewDown() // scroll to the bottom as new messages arrive
		}

	case tea.KeyMsg:
		switch msg.String() {

		// toggle text input focus
		case "esc":
			if m.textInput.Focused() {
				m.textInput.Blur()
			} else {
				m.textInput.Focus()
			}

		case "enter":

			// skip if the text input is empty or if text input is blurred
			if m.textInput.Value() == "" || !m.textInput.Focused() {
				return m, nil
			}

			// set the current message to the text input value
			m.messages = append(m.messages, openai.UserMessage(m.textInput.Value()))

			//empty out the text input
			m.textInput.SetValue("")

			m.stream = m.client.Chat.Completions.NewStreaming(m.ctx, openai.ChatCompletionNewParams{
				Messages: openai.F(m.messages),
				Seed:     openai.Int(1),
				Model:    openai.F(openai.ChatModelGPT3_5Turbo),
			})

			//if stream.Next() is true then send a custom command with content payload: evt.Choices[0].Delta.Content
			if m.stream.Next() {
				evt := m.stream.Current()
				if len(evt.Choices) > 0 {

					//kick off the compilation of chatgpt message stream
					return m, func() tea.Msg {
						return execMsg{
							msg: evt.Choices[0].Delta.Content,
						}
					}
				}

			} else {
				// Add assistant's response to history
				m.messages = append(m.messages, openai.AssistantMessage(m.assistantResponse.String()))
			}
		}
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *ChatGPTdialogueWindow) View() string {

	output := ""
	output += "\n" + m.viewport.View() + "\n" + m.textInput.View() + "\n" + m.statusbar.View()
	return output
}
