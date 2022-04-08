package style

import (
  "github.com/charmbracelet/lipgloss"
)

// XXX: For now, this is in its own package so that it can be shared between
// different packages without incurring an illegal import cycle.

// Styles defines styles for the TUI.
type Styles struct {
  App						lipgloss.Style
  List						lipgloss.Style
  Query						lipgloss.Style
  ResultListItem			lipgloss.Style
  PinnedListItem			lipgloss.Style
  MutedText					lipgloss.Style
  SelectedMutedText			lipgloss.Style
  MutedTextUnterlined		lipgloss.Style
  Text						lipgloss.Style
  SelectedText				lipgloss.Style
  SuccessText				lipgloss.Style
  ErrorText					lipgloss.Style
}

// DefaultStyles returns default styles for the TUI.
func DefaultStyles() (s Styles) {
	s.App = lipgloss.NewStyle()
	s.List = lipgloss.NewStyle()
	s.PinnedListItem = lipgloss.NewStyle().
						MarginLeft(2)

	s.ResultListItem = s.PinnedListItem.Copy()

	s.MutedText = lipgloss.NewStyle().
					Foreground(lipgloss.Color("#555555"))

	s.SelectedMutedText = s.MutedText.Copy().
					Foreground(lipgloss.Color("#FF00FF"))

	s.MutedTextUnterlined = s.MutedText.Copy().
					Underline(true)

	s.Text = lipgloss.NewStyle().
					Foreground(lipgloss.Color("#FFFFFF"))

	s.SelectedText = s.Text.Copy().
					Foreground(lipgloss.Color("#FF00FF")).
					Bold(true)

	s.SuccessText = lipgloss.NewStyle().
					Foreground(lipgloss.Color("#1BEE91"))
	s.ErrorText = lipgloss.NewStyle().
					Foreground(lipgloss.Color("#FF4C94"))

	s.Query = lipgloss.NewStyle().
					Background(lipgloss.Color("#1C1C1C")).
					Margin(1, 0)

  return s
}
