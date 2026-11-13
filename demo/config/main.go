// Command demo is the goda-cmp configuration used by demo/record.sh: the
// fuzzy filters over gopls (the nvimlsp source of goda-cmp-ext-luasrc.vim)
// and the buffer source.
package main

import (
	"github.com/Shougo/goda-cmp-ext/ext/source/buffer"
	"github.com/Shougo/goda-cmp.vim/cmp"
	"github.com/delphinus/goda-cmp-ext-luasrc.vim/ext/source/luasrc"
	"github.com/delphinus/goda-cmp-ext-luasrc.vim/ext/ui/luasrcpum"

	"github.com/delphinus/goda-cmp-ext-fuzzy.vim/ext/filter/matcherfuzzy"
	"github.com/delphinus/goda-cmp-ext-fuzzy.vim/ext/filter/matchernottyped"
	"github.com/delphinus/goda-cmp-ext-fuzzy.vim/ext/filter/sorterfuzzy"
	"github.com/delphinus/goda-cmp-ext-fuzzy.vim/ext/filter/sorterglobal"
)

func main() {
	bridge := luasrc.NewBridge()

	b := cmp.New()
	b.AddUI(luasrcpum.New())
	b.AddSource(bridge.Source("nvimlsp", luasrc.NvimLSPModule))
	b.AddSource(buffer.New(buffer.Params{}))
	b.AddFilter(matcherfuzzy.New())
	b.AddFilter(sorterfuzzy.New())
	b.AddFilter(matchernottyped.New())
	b.AddFilter(sorterglobal.New(nil))
	b.PatchGlobal(cmp.Config{
		UI:          "luasrcpum",
		Sources:     []string{"nvimlsp", "buffer"},
		PostFilters: []string{"sorterglobal"},
		SourceOptions: map[string]cmp.SourceOption{
			"_": {
				Matchers: []string{"matcherfuzzy"},
				Sorters:  []string{"sorterfuzzy"},
			},
			"nvimlsp": {
				Mark:                   "LSP",
				ForceCompletionPattern: `\.\w*`,
			},
			"buffer": {
				Mark:     "B",
				Matchers: []string{"matcherfuzzy", "matchernottyped"},
			},
		},
	})
	b.Main(bridge)
}
