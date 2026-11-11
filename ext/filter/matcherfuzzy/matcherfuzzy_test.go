package matcherfuzzy

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

func TestMatcher(t *testing.T) {
	for _, tc := range []struct {
		typed string
		in    []string
		want  []string
	}{
		// Case sensitive, as the typed text has an upper case letter: the L of PrintLine does not match l.
		{"Prl", []string{"Print", "Println", "Printf", "PrintLine"}, []string{"Println"}},
		{"prl", []string{"Println", "parallel", "pr"}, []string{"Println", "parallel"}},
		// Typing upper case letters makes it case sensitive.
		{"PL", []string{"PrintLine", "println"}, []string{"PrintLine"}},
		{"", []string{"a", "b"}, []string{"a", "b"}},
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
