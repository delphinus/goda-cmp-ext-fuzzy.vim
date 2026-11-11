// Package matchernottyped provides the matcher "matchernottyped", which drops the
// candidate that is the very word being typed.
//
// The sources around and buffer also gather the word under the cursor, so
// typing "Prl" offers "Prl" itself. The buffer source of blink.cmp does not,
// and neither does this.
package matchernottyped

import (
	"context"

	"github.com/Shougo/goda-cmp.vim/cmp"
)

type filter struct{}

// New returns the matcher "matchernottyped".
func New() cmp.Filter { return filter{} }

func (filter) Name() string { return "matchernottyped" }

func (filter) Filter(_ context.Context, in cmp.FilterInput) ([]cmp.Item, error) {
	out := make([]cmp.Item, 0, len(in.Items))
	for _, it := range in.Items {
		if it.Word != in.Complete {
			out = append(out, it)
		}
	}
	return out, nil
}
