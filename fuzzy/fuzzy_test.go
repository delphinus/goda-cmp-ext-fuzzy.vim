package fuzzy

import (
	"slices"
	"testing"
)

func TestMatchPositions(t *testing.T) {
	for _, tc := range []struct {
		text, typed string
		want        []int
	}{
		{"fooBar", "fob", []int{0, 1, 3}},
		// Narrowed backward, so the match starts at the f nearer to b, not the first f.
		{"fx_fooBar", "fob", []int{3, 4, 6}},
		{"abc", "x", nil},
	} {
		_, got, _ := Match([]rune(tc.text), []rune(tc.typed), Fold(tc.typed))
		if !slices.Equal(got, tc.want) {
			t.Errorf("%q %q: got %v, want %v", tc.text, tc.typed, got, tc.want)
		}
	}
}
