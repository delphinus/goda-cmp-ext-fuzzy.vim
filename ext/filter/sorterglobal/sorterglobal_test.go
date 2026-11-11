package sorterglobal

import (
	"context"
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

func TestGlobal(t *testing.T) {
	in := []cmp.Item{
		// A source starting after the colon: goda-cmp has prepended it to Word.
		{Word: ":streaming-misc", FilterText: "streaming-misc", Source: "ghq"},
		{Word: "😄", FilterText: ":smile:", Source: "emoji"},
		{Word: ":sm", Source: "around"},
	}
	got, err := New(map[string]int{"emoji": 15}).Filter(context.Background(), cmp.FilterInput{Complete: ":smi", Items: in})
	if err != nil {
		t.Fatal(err)
	}
	if got[0].Source != "emoji" {
		t.Errorf("got %q first, want emoji", got[0].Source)
	}
	if s := globalScore(in[0], []rune(":smi"), true); s <= 0 {
		t.Errorf("the prefixed item should match without its prefix, got %d", s)
	}
}
