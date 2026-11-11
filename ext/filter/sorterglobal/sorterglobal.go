// Package sorterglobal provides the post filter "sorterglobal", which sorts
// the candidates of all the sources together by how well they match.
package sorterglobal

import (
	"context"
	"slices"
	"strings"

	"github.com/Shougo/goda-cmp.vim/cmp"

	"github.com/delphinus/goda-cmp-ext-fuzzy.vim/fuzzy"
)

type global struct{ offsets map[string]int }

// ScoreOffsetter is implemented by the UserData of a candidate with a score
// offset of its own (like score_offset of a blink.cmp item), which
// "sorterglobal" adds.
type ScoreOffsetter interface{ ScoreOffset() int }

// New returns the sorter "sorterglobal", to be put in PostFilters.
// goda-cmp concatenates the candidates sorted per source in the order of the
// sources; this sorts all of them together by how well they match, as
// blink.cmp does. offsets are added per source (like score_offset of a
// blink.cmp provider), and so is the offset of a candidate whose UserData is a
// ScoreOffsetter. Candidates with the same score keep their order (that of
// the sources).
func New(offsets map[string]int) cmp.Filter { return global{offsets} }

func (global) Name() string { return "sorterglobal" }

func (g global) Filter(_ context.Context, in cmp.FilterInput) ([]cmp.Item, error) {
	typed := []rune(in.Complete)
	f := fuzzy.FoldInput(in)
	scores := make([]int, len(in.Items))
	for i, it := range in.Items {
		scores[i] = globalScore(it, typed, f) + g.offsets[it.Source]
		if o, ok := it.UserData.(ScoreOffsetter); ok {
			scores[i] += o.ScoreOffset()
		}
	}
	idx := make([]int, len(in.Items))
	for i := range idx {
		idx[i] = i
	}
	slices.SortStableFunc(idx, func(a, b int) int { return scores[b] - scores[a] })
	out := make([]cmp.Item, len(idx))
	for i, j := range idx {
		out[i] = in.Items[j]
	}
	return out, nil
}

// globalScore is the score of a candidate to compare it across sources, or 0
// when it does not match.
//
// When sources start at different columns, goda-cmp aligns them to the
// leftmost one by prepending the text in between to Word, but not to
// FilterText. So the head of the typed text prepended to Word (typed[:k]) is
// skipped before matching, trying the smallest k first.
func globalScore(it cmp.Item, typed []rune, fold bool) int {
	if it.FilterText == "" {
		s, _ := fuzzy.Score([]rune(it.Word), typed, fold)
		return s
	}
	text := []rune(it.FilterText)
	for k := 0; k <= len(typed); k++ {
		if !strings.HasPrefix(it.Word, string(typed[:k])) {
			break
		}
		if s, ok := fuzzy.Score(text, typed[k:], fold); ok {
			return s
		}
	}
	return 0
}
