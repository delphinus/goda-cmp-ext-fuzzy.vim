// Package sorterfuzzy provides the sorter "sorterfuzzy", which sorts the
// candidates by how well they match the typed characters (see package fuzzy).
package sorterfuzzy

import (
	"context"
	"slices"

	"github.com/Shougo/goda-cmp.vim/cmp"

	"github.com/delphinus/goda-cmp-ext-fuzzy.vim/fuzzy"
)

type sorter struct{}

// New returns the sorter "sorterfuzzy". Candidates with the same score
// are sorted shorter first, and those of the same length keep their order:
// sort them with another sorter (sortertext, for example) before this one to
// break the remaining ties.
func New() cmp.Filter { return sorter{} }

func (sorter) Name() string { return "sorterfuzzy" }

func (sorter) Filter(_ context.Context, in cmp.FilterInput) ([]cmp.Item, error) {
	if in.Complete == "" {
		return in.Items, nil
	}
	pat, f := []rune(in.Complete), fuzzy.FoldInput(in)
	type scored struct {
		it    cmp.Item
		score int
		n     int
	}
	ss := make([]scored, len(in.Items))
	for i, it := range in.Items {
		text := []rune(it.Filter())
		s, _ := fuzzy.Score(text, pat, f)
		ss[i] = scored{it, s, len(text)}
	}
	slices.SortStableFunc(ss, func(a, b scored) int {
		if a.score != b.score {
			return b.score - a.score
		}
		return a.n - b.n
	})
	out := make([]cmp.Item, len(ss))
	for i, s := range ss {
		out[i] = s.it
	}
	return out, nil
}
