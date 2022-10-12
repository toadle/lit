package list

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"

	"lit/internal/tui/style"
)

type FilterState int

const (
	Unfiltered FilterState = iota
	Filtered
)

func (f FilterState) String() string {
	return [...]string{
		"unfiltered",
		"filtered",
	}[f]
}

type Model struct {
	height int
	items  []Item
	styles style.Styles

	cursor           int
	windowBeginIndex int
	windowEndIndex   int

	filterState   FilterState
	filterValue   string
	filteredItems []FilteredItem

	noResultText string
}

func New(items []Item, height int) Model {
	m := Model{
		cursor:           -1,
		styles:           style.DefaultStyles(),
		height:           height,
		windowBeginIndex: 0,
		windowEndIndex:   0,
		items:            items,
		filterState:      Unfiltered,
		noResultText:     "No items found.",
	}
	return m
}

func (m *Model) SetNoResultText(str string) {
	m.noResultText = str
}

func (m *Model) Select(index int) {
	m.cursor = index
}

func (m *Model) Unselect() {
	m.cursor = -1
}

func (m *Model) Index() int {
	return m.cursor
}

func (m *Model) SetItems(i []Item) {
	m.items = i
}
func (m Model) Items() []Item {
	return m.items
}

func (m *Model) SetHeight(h int) {
	m.height = h
	if m.windowEndIndex < m.height {
		m.windowEndIndex = h
	}
}
func (m *Model) Height() int {
	return m.height
}

func (m Model) SelectedItem() Item {
	i := m.Index()

	items := m.VisibleItems()
	if i < 0 || len(items) == 0 || len(items) <= i {
		return nil
	}

	return items[i]
}

func (m *Model) CursorUp() {
	if m.cursor > -1 {
		if m.cursor != 0 {
			m.cursor--
		}

		if m.cursor < m.windowBeginIndex {
			m.windowBeginIndex--
			m.windowEndIndex--
		}
	}
}

func (m *Model) CursorDown() {
	if m.cursor > -1 {
		if m.cursor == m.windowEndIndex-1 && m.windowEndIndex < len(m.items)-1 {
			m.windowBeginIndex++
			m.windowEndIndex++
		}
		m.cursor++
	}
}

func (m *Model) SetFilterValue(term string) {
	m.filterValue = term
}

func (m *Model) UnfilterItems() {
	m.filteredItems = []FilteredItem{}
	m.filterState = Unfiltered
}

func (m *Model) FilterItems() {
	if m.filterValue == "" {
		return
	}
	m.filterState = Filtered

	targets := lo.Map[Item, string](m.items, func(i Item, _ int) string {
		return i.FilterValue()
	})

	rankedItems := FuzzyFilter(m.filterValue, targets)

	m.filteredItems = lo.Map[Rank, FilteredItem](rankedItems, func(r Rank, _ int) FilteredItem {
		return FilteredItem{
			item:    m.items[r.Index],
			matches: r.MatchedIndexes,
		}
	})
	m.filterState = Filtered
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var teaCmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.CursorUp()
		case "down":
			m.CursorDown()
		}
	}

	return m, tea.Batch(teaCmds...)

}

func (m Model) View() string {
	var sections []string
	sections = append(sections, m.populatedView())

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m Model) VisibleItems() []Item {
	if m.filterState != Unfiltered {
		return lo.Map[FilteredItem, Item](m.filteredItems, func(i FilteredItem, _ int) Item {
			return i.item
		})
	}
	return m.items
}

func (m Model) MatchesForItem(index int) []int {
	if m.filteredItems == nil || index >= len(m.filteredItems) {
		return nil
	}
	return m.filteredItems[index].matches
}

func (m Model) populatedView() string {
	items := m.VisibleItems()

	var b strings.Builder

	// Empty states
	if len(items) == 0 {
		return m.styles.NoResultItem.Render(m.noResultText)
	}

	endIndex := m.windowEndIndex
	if m.windowEndIndex > len(items) {
		endIndex = len(items)
	}

	if len(items) > 0 {
		windowItems := items[m.windowBeginIndex:endIndex]
		for i, item := range windowItems {
			item.Render(&b, m, m.windowBeginIndex+i, item)
			if i < len(windowItems)-1 {
				fmt.Fprint(&b, "\n")
			}
		}
	}

	return b.String()
}
