// Package matcherfuzzy provides the matcher "matcherfuzzy", which keeps the
// candidates containing the typed characters in order (not necessarily
// adjacent). The matchers bundled with goda-cmp match the head of a word only
// (matcherhead, matcherprefix); this lets "Prl" find Println, as blink.cmp does.
package matcherfuzzy

import (
	"context"

	"github.com/Shougo/goda-cmp.vim/cmp"

	"github.com/delphinus/goda-cmp-ext-fuzzy.vim/fuzzy"
)

type matcher struct{}

// New returns the matcher "matcherfuzzy".
func New() cmp.Filter { return matcher{} }

func (matcher) Name() string { return "matcherfuzzy" }

func (matcher) Filter(_ context.Context, in cmp.FilterInput) ([]cmp.Item, error) {
	pat, f := []rune(in.Complete), fuzzy.FoldInput(in)
	out := make([]cmp.Item, 0, len(in.Items))
	for _, it := range in.Items {
		if _, ok := fuzzy.Score([]rune(it.Filter()), pat, f); ok {
			out = append(out, it)
		}
	}
	return out, nil
}
