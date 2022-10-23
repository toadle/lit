package list

import (
	"sort"

	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"
)

type Rank struct {
	// The index of the item in the original input.
	Index int
	// Indices of the actual word that were matched against the filter term.
	MatchedIndexes []int
}

func FuzzyFilter(term string, targets []string) []Rank {
	var ranks fuzzy.Matches = fuzzy.Find(term, targets)
	sort.Stable(ranks)

	return lo.Map(ranks, func(m fuzzy.Match, _ int) Rank {
		return Rank{
			Index:          m.Index,
			MatchedIndexes: m.MatchedIndexes,
		}
	})
}
