package style

import (
	"github.com/charmbracelet/lipgloss"
)

// XXX: For now, this is in its own package so that it can be shared between
// different packages without incurring an illegal import cycle.

// Styles defines styles for the TUI.
type Styles struct {
	App                        lipgloss.Style
	List                       lipgloss.Style
	Query                      lipgloss.Style
	QueryPromptFocused         lipgloss.Style
	SearchListItem             lipgloss.Style
	SearchListItemSelected     lipgloss.Style
	CalculatorListItem         lipgloss.Style
	CalculatorListItemSelected lipgloss.Style
	NoResultItem               lipgloss.Style
	MutedText                  lipgloss.Style
	SelectedMutedText          lipgloss.Style
	MutedTextUnterlined        lipgloss.Style
	Text                       lipgloss.Style
	SelectedText               lipgloss.Style
	SuccessText                lipgloss.Style
	ErrorText                  lipgloss.Style
	CalculatorText             lipgloss.Style
	MutedCalculatorText        lipgloss.Style
}

// DefaultStyles returns default styles for the TUI.
func DefaultStyles() (s Styles) {
	blue := lipgloss.Color("33")
	lightBlue := lipgloss.Color("45")
	yellow := lipgloss.Color("214")
	lightYellow := lipgloss.Color("221")

	s.App = lipgloss.NewStyle().
		MarginTop(1).
		MarginBottom(1)

	s.List = lipgloss.NewStyle()
	s.CalculatorListItem = lipgloss.NewStyle().
		MarginLeft(2)

	s.CalculatorListItemSelected = s.CalculatorListItem.Copy().
		MarginLeft(0).
		PaddingLeft(1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(blue).
		BorderLeft(true).
		Foreground(blue)

	s.SearchListItem = s.CalculatorListItem.Copy()

	s.SearchListItemSelected = s.CalculatorListItemSelected.Copy().
		BorderForeground(yellow)

	s.NoResultItem = s.CalculatorListItem.Copy().
		Foreground(lipgloss.Color("#777"))

	s.MutedText = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#777"))

	s.SelectedMutedText = s.MutedText.Copy().
		Foreground(lightYellow)

	s.MutedTextUnterlined = s.MutedText.Copy().
		Underline(true)

	s.Text = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))

	s.SelectedText = s.Text.Copy().
		Foreground(yellow).
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
		Foreground(blue)

	s.MutedCalculatorText = s.Text.Copy().
		Foreground(lightBlue)

	return s
}
