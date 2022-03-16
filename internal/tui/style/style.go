package style

import (
  "github.com/charmbracelet/lipgloss"
)

// XXX: For now, this is in its own package so that it can be shared between
// different packages without incurring an illegal import cycle.

// Styles defines styles for the TUI.
type Styles struct {
  App lipgloss.Style
}

// DefaultStyles returns default styles for the TUI.
func DefaultStyles() *Styles {
  s := new(Styles)

  s.App = lipgloss.NewStyle().
    Margin(1, 2)

  // s.Header = lipgloss.NewStyle().
  //   Foreground(lipgloss.Color("62")).
  //   Align(lipgloss.Right).
  //   Bold(true)

  return s
}