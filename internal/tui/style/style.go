package style

import (
  "github.com/charmbracelet/lipgloss"
)

// XXX: For now, this is in its own package so that it can be shared between
// different packages without incurring an illegal import cycle.

// Styles defines styles for the TUI.
type Styles struct {
  App						lipgloss.Style
  Query						lipgloss.Style
  ResultListItem			lipgloss.Style
  PinnedListItem			lipgloss.Style
  SelectedPinnedListItem	lipgloss.Style
  MutedText					lipgloss.Style
  MutedTextUnterlined		lipgloss.Style
  Text						lipgloss.Style
  SuccessText				lipgloss.Style
  ErrorText					lipgloss.Style
}

// DefaultStyles returns default styles for the TUI.
func DefaultStyles() *Styles {
	s := new(Styles)

	s.App = lipgloss.NewStyle()
	s.PinnedListItem = lipgloss.NewStyle().
						MarginLeft(2)

	s.ResultListItem = s.PinnedListItem.Copy()

	s.MutedText = lipgloss.NewStyle().
					Foreground(lipgloss.Color("#555555"))
	s.MutedTextUnterlined = s.MutedText.Copy().
					Underline(true)

	s.Text = lipgloss.NewStyle().
					Foreground(lipgloss.Color("#FFFFFF"))

	s.SuccessText = lipgloss.NewStyle().
					Foreground(lipgloss.AdaptiveColor{Light: "#1BEE91", Dark: "#1BEE91"})
	s.ErrorText = lipgloss.NewStyle().
					Foreground(lipgloss.AdaptiveColor{Light: "#FF4C94", Dark: "#FF4C94"})

	s.Query = lipgloss.NewStyle().
					Background(lipgloss.Color("#1C1C1C")).
					Margin(1, 0)

  return s
}
