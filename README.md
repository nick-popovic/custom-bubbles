# Custom Bubble Library

![Status](https://img.shields.io/badge/Status-Development-yellow)
![Version](https://img.shields.io/badge/Version-0.0.1-blue)

> ⚠️ **Note**: This is an initial development version and not ready for production use. The library is currently being structured for proper Go package distribution.

A collection of reusable Bubbles built with the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework for terminal user interfaces.

## Development Status
- [x] Initial commit with core functionality
- [x] Package structure being finalized
- [x] Documentation in progress
- [x] Testing and examples
- [ ] verified import via `go get`

## Installation

```bash
go get github.com/nick-popovic/custom-bubbles/multiTab
go get github.com/nick-popovic/custom-bubbles/chatGPT
go get github.com/nick-popovic/custom-bubbles/candleStick
```

## Components

### ChatGPT Terminal Interface
A terminal-based ChatGPT interface with markdown rendering and real-time response streaming.
- [Documentation](chatGPT/README.md)
- Features: Markdown rendering, real-time streaming, support for GPT-3.5-turbo model

### Multi-Tab Component
A flexible tab management system for organizing multiple views.
- [Documentation](multiTab/README.md)
- Features: Independent tab state, customizable tab styling, modular design

### CandleStick Chart
A terminal-based financial chart component for displaying OHLC data.
- [Documentation](candleStick/README.md)
- Features: Real-time updates, color-coded indicators