package candleStick

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	bullishStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00")) // Green
	bearishStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")) // Red
)

type Candlestick struct {
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    int
	Timestamp time.Time
}

type Model struct {
	Candlesticks []Candlestick
	width        int
	offset       int // Track scroll position
}

func GenerateSampleData() []Candlestick {
	startDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	data := []Candlestick{
		{Open: 100, High: 100, Low: 95, Close: 105, Timestamp: startDate},
		{Open: 105, High: 115, Low: 100, Close: 95, Timestamp: startDate.AddDate(0, 0, 1)},
		{Open: 95, High: 100, Low: 90, Close: 98, Timestamp: startDate.AddDate(0, 0, 2)},
		{Open: 98, High: 103, Low: 96, Close: 102, Timestamp: startDate.AddDate(0, 0, 3)},
		{Open: 102, High: 108, Low: 100, Close: 107, Timestamp: startDate.AddDate(0, 0, 4)},
		{Open: 107, High: 112, Low: 105, Close: 110, Timestamp: startDate.AddDate(0, 0, 5)},
		{Open: 110, High: 115, Low: 108, Close: 112, Timestamp: startDate.AddDate(0, 0, 6)},
		{Open: 112, High: 118, Low: 110, Close: 115, Timestamp: startDate.AddDate(0, 0, 7)},
		{Open: 115, High: 120, Low: 113, Close: 118, Timestamp: startDate.AddDate(0, 0, 8)},
		{Open: 118, High: 125, Low: 117, Close: 123, Timestamp: startDate.AddDate(0, 0, 9)},
		{Open: 100, High: 100, Low: 95, Close: 105, Timestamp: startDate.AddDate(0, 0, 10)},
		{Open: 105, High: 115, Low: 100, Close: 95, Timestamp: startDate.AddDate(0, 0, 11)},
		{Open: 95, High: 100, Low: 90, Close: 98, Timestamp: startDate.AddDate(0, 0, 12)},
		{Open: 98, High: 103, Low: 96, Close: 102, Timestamp: startDate.AddDate(0, 0, 13)},
		{Open: 102, High: 108, Low: 100, Close: 107, Timestamp: startDate.AddDate(0, 0, 14)},
		{Open: 107, High: 112, Low: 105, Close: 110, Timestamp: startDate.AddDate(0, 0, 15)},
		{Open: 110, High: 115, Low: 108, Close: 112, Timestamp: startDate.AddDate(0, 0, 16)},
		{Open: 112, High: 118, Low: 110, Close: 115, Timestamp: startDate.AddDate(0, 0, 17)},
		{Open: 115, High: 120, Low: 113, Close: 118, Timestamp: startDate.AddDate(0, 0, 18)},
		{Open: 118, High: 125, Low: 117, Close: 123, Timestamp: startDate.AddDate(0, 0, 19)},
		{Open: 100, High: 100, Low: 95, Close: 105, Timestamp: startDate.AddDate(0, 0, 20)},
		{Open: 105, High: 115, Low: 100, Close: 95, Timestamp: startDate.AddDate(0, 0, 21)},
		{Open: 95, High: 100, Low: 90, Close: 98, Timestamp: startDate.AddDate(0, 0, 22)},
		{Open: 98, High: 103, Low: 96, Close: 102, Timestamp: startDate.AddDate(0, 0, 23)},
		{Open: 102, High: 108, Low: 100, Close: 107, Timestamp: startDate.AddDate(0, 0, 24)},
		{Open: 107, High: 112, Low: 105, Close: 110, Timestamp: startDate.AddDate(0, 0, 25)},
		{Open: 110, High: 115, Low: 108, Close: 112, Timestamp: startDate.AddDate(0, 0, 26)},
		{Open: 112, High: 118, Low: 110, Close: 115, Timestamp: startDate.AddDate(0, 0, 27)},
		{Open: 115, High: 120, Low: 113, Close: 118, Timestamp: startDate.AddDate(0, 0, 28)},
		{Open: 118, High: 125, Low: 117, Close: 123, Timestamp: startDate.AddDate(0, 0, 29)},

		// addd more data to cross months
		{Open: 95, High: 100, Low: 90, Close: 98, Timestamp: startDate.AddDate(0, 0, 30)},
		{Open: 98, High: 103, Low: 96, Close: 102, Timestamp: startDate.AddDate(0, 0, 31)},
		{Open: 102, High: 108, Low: 100, Close: 107, Timestamp: startDate.AddDate(0, 0, 32)},
		{Open: 107, High: 112, Low: 105, Close: 110, Timestamp: startDate.AddDate(0, 0, 33)},
		{Open: 110, High: 115, Low: 108, Close: 112, Timestamp: startDate.AddDate(0, 0, 34)},
		{Open: 112, High: 118, Low: 110, Close: 115, Timestamp: startDate.AddDate(0, 0, 35)},
		{Open: 115, High: 120, Low: 113, Close: 118, Timestamp: startDate.AddDate(0, 0, 36)},
		{Open: 118, High: 125, Low: 117, Close: 123, Timestamp: startDate.AddDate(0, 0, 37)},
	}

	// Update all entries with sequential dates
	for i := range data[2:] {
		data[i+2].Timestamp = startDate.AddDate(0, 0, i+2)
	}

	return data
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, tea.ClearScreen

	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			if m.offset > 0 {
				m.offset--
				return m, tea.ClearScreen
			}
		case tea.MouseWheelDown:
			if m.offset < len(m.Candlesticks)-1 {
				m.offset++
				return m, tea.ClearScreen
			}
		}

	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
		if msg.String() == "u" {
		}

	}
	return m, nil
}

func (m Model) View() string {
	if len(m.Candlesticks) == 0 {
		return "No data to display\n"
	}

	// Find min, max values and calculate label width
	minVal, maxVal := m.Candlesticks[0].Low, m.Candlesticks[0].High
	for _, c := range m.Candlesticks {
		if c.Low < minVal {
			minVal = c.Low
		}
		if c.High > maxVal {
			maxVal = c.High
		}
	}

	// Calculate price label width
	labelWidth := len(fmt.Sprintf("%.1f", maxVal)) + 1
	padding := fmt.Sprintf("%%%d.1f", labelWidth)

	var output string
	height := 15
	valueRange := maxVal - minVal

	// Draw chart rows
	for i := 0; i < height; i++ {
		price := maxVal - (valueRange * float64(i) / float64(height))
		output += fmt.Sprintf(padding+" │", price)

		for _, candle := range m.Candlesticks {
			normalized := func(val float64) int {
				return int(float64(height-1) * (maxVal - val) / valueRange)
			}

			highPos := normalized(candle.High)
			lowPos := normalized(candle.Low)
			openPos := normalized(candle.Open)

			switch {
			case i == openPos:
				if candle.Open > candle.Close {
					output += "  " + bearishStyle.Render("▼") + "  "
				} else {
					output += "  " + bullishStyle.Render("▲") + "  "
				}
			case i > highPos && i < lowPos:
				output += "  │  "
			default:
				output += "     "
			}
		}
		output += "\n"
	}

	// Draw bottom axis with correct spacing
	output += strings.Repeat(" ", labelWidth) + " └"
	for i := 0; i < len(m.Candlesticks); i++ {
		output += "─────"
	}
	output += "\n"

	// Add day numbers
	output += strings.Repeat(" ", labelWidth) + "  "
	for _, candle := range m.Candlesticks {
		output += fmt.Sprintf(" %02d  ", candle.Timestamp.Day())
	}
	output += "\n"

	// Add month labels where month changes
	output += strings.Repeat(" ", labelWidth) + "  "
	for i, candle := range m.Candlesticks {
		if i == 0 || candle.Timestamp.Month() != m.Candlesticks[i-1].Timestamp.Month() {
			month := candle.Timestamp.Format("Jan")
			output += fmt.Sprintf(" %s  ", month)
			i += len(month) // Skip spaces for month name
		} else {
			output += "     " // Match candlestick width
		}
	}
	output += "\n"

	return output
}
