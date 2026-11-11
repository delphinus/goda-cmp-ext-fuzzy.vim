package sorterfuzzy

import (
	"context"
	"slices"
	"testing"

	"github.com/Shougo/goda-cmp.vim/cmp"
)

func words(items []cmp.Item) []string {
	out := make([]string, len(items))
	for i, it := range items {
		out[i] = it.Word
	}
	return out
}

func items(ws ...string) []cmp.Item {
	out := make([]cmp.Item, len(ws))
	for i, w := range ws {
		out[i] = cmp.Item{Word: w}
	}
	return out
}

func TestSorter(t *testing.T) {
	for _, tc := range []struct {
		typed string
		in    []string
		want  []string
	}{
		// A consecutive match from the head of a word comes first.
		{"get", []string{"budget", "target", "getLines", "get"}, []string{"get", "getLines", "budget", "target"}},
		// A match at camelCase heads comes before one in the middle of a word.
		{"gl", []string{"angle", "getLines"}, []string{"getLines", "angle"}},
		{"bgl", []string{"nvim_buf_get_lines", "nvim_buf_get_lines_extra", "bagel"}, []string{"nvim_buf_get_lines", "nvim_buf_get_lines_extra", "bagel"}},
		// The same score keeps the order.
		{"a", []string{"ab", "ac"}, []string{"ab", "ac"}},
	} {
		got, err := New().Filter(context.Background(), cmp.FilterInput{Complete: tc.typed, Items: items(tc.in...)})
		if err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(words(got), tc.want) {
			t.Errorf("%q: got %q, want %q", tc.typed, words(got), tc.want)
		}
	}
}
