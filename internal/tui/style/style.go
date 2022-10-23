package style

import (
	"github.com/charmbracelet/lipgloss"
)

// XXX: For now, this is in its own package so that it can be shared between
// different packages without incurring an illegal import cycle.

// Styles defines styles for the TUI.
type Styles struct {
	App                 lipgloss.Style
	List                lipgloss.Style
	Query               lipgloss.Style
	QueryPromptFocused  lipgloss.Style
	SearchListItem      lipgloss.Style
	CalculatorListItem  lipgloss.Style
	NoResultItem        lipgloss.Style
	MutedText           lipgloss.Style
	SelectedMutedText   lipgloss.Style
	MutedTextUnterlined lipgloss.Style
	Text                lipgloss.Style
	SelectedText        lipgloss.Style
	SuccessText         lipgloss.Style
	ErrorText           lipgloss.Style
	CalculatorText      lipgloss.Style
	MutedCalculatorText lipgloss.Style
}

// DefaultStyles returns default styles for the TUI.
func DefaultStyles() (s Styles) {
	s.App = lipgloss.NewStyle()
	s.List = lipgloss.NewStyle()
	s.CalculatorListItem = lipgloss.NewStyle().
		MarginLeft(2)

	s.SearchListItem = s.CalculatorListItem.Copy()
	s.NoResultItem = s.CalculatorListItem.Copy().
		Foreground(lipgloss.Color("#777"))

	s.MutedText = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#777"))

	s.SelectedMutedText = s.MutedText.Copy().
		Foreground(lipgloss.Color("#E5C352"))

	s.MutedTextUnterlined = s.MutedText.Copy().
		Underline(true)

	s.Text = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))

	s.SelectedText = s.Text.Copy().
		Foreground(lipgloss.Color("#F5C90C")).
		Bold(true)

	s.SuccessText = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#1BEE91"))
	s.ErrorText = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF4C94"))

	s.Query = lipgloss.NewStyle().
		Margin(1, 0)

	s.QueryPromptFocused = s.Text.Copy().
		Foreground(lipgloss.Color("#E94904")).
		Bold(true)

	s.CalculatorText = s.Text.Copy().
		Foreground(lipgloss.Color("#00BECD"))

	s.MutedCalculatorText = s.Text.Copy().
		Foreground(lipgloss.Color("#B4E0E7"))

	return s
}
