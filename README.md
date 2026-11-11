# goda-cmp-ext-fuzzy.vim

[日本語](README.ja.md)

> [!WARNING]
> This is a beta (v0.0.x). goda-cmp.vim itself is under active development,
> and its bundled extensions are being moved into repositories of their own, so
> the API, the import paths and even the name of this repository may change.

Fuzzy filters for [goda-cmp.vim](https://github.com/Shougo/goda-cmp.vim).

| Name | Kind | What it does |
|---|---|---|
| `matcherfuzzy` | matcher | Keeps the candidates that contain the typed characters in order, not necessarily adjacent: `Prl` finds `Println`. |
| `sorterfuzzy` | sorter | Sorts the candidates by how well they match: consecutive characters, heads of words (after a delimiter, camelCase) and the start of the candidate score higher. |
| `sorterglobal` | post filter | Sorts the candidates of all the sources together by how well they match, with offsets per source and per candidate. |
| `matchernottyped` | matcher | Drops the candidate that is the very word being typed (`around` and `buffer` also gather the word under the cursor). |

Case is smartcase: ignored unless the typed text has an upper case letter, as
`matcherhead` does. The `ignoreCase` option of a source ignores it always.

## Why

The matchers bundled with goda-cmp match the head of a word only
(`matcherhead`, `matcherprefix`). And goda-cmp concatenates the candidates
sorted per source in the order of the sources, so a weak match of the first
source comes before an exact match of the last one. These filters make
goda-cmp behave more like blink.cmp in both respects.

Matching and scoring follow fzf's v1 algorithm in a simple form. Typos are not
tolerated.

## Install

These are Go packages that your goda-cmp configuration imports. Either
require them in its `go.mod`:

```sh
cd ~/.config/nvim/goda-cmp   # your configuration module
go get github.com/delphinus/goda-cmp-ext-fuzzy.vim@latest
go mod tidy
```

or install the repository with your plugin manager and add its directory to
`g:goda#workspace`, so that `goda_cmp#build()` builds with the installed
version (goda-pm.vim does this for you). With lazy.nvim, for example:

```lua
{
  "Shougo/goda-cmp.vim",
  dependencies = { "Shougo/goda.vim", "delphinus/goda-cmp-ext-fuzzy.vim" },
  init = function()
    local root = require("lazy.core.config").options.root
    vim.g["goda#workspace"] = vim.tbl_map(function(name)
      return root .. "/" .. name
    end, { "goda.vim", "goda-cmp.vim", "goda-cmp-ext-fuzzy.vim" })
  end,
}
```

The module still has to be required in `go.mod`; the workspace only replaces
its version.

## Configure

In the `main.go` of your configuration:

```go
import (
	"github.com/Shougo/goda-cmp.vim/cmp"

	"github.com/delphinus/goda-cmp-ext-fuzzy.vim/ext/filter/matcherfuzzy"
	"github.com/delphinus/goda-cmp-ext-fuzzy.vim/ext/filter/matchernottyped"
	"github.com/delphinus/goda-cmp-ext-fuzzy.vim/ext/filter/sorterfuzzy"
	"github.com/delphinus/goda-cmp-ext-fuzzy.vim/ext/filter/sorterglobal"
)

func main() {
	b := cmp.New()
	// ... UIs and sources ...

	b.AddFilter(matcherfuzzy.New())
	b.AddFilter(sorterfuzzy.New())
	b.AddFilter(matchernottyped.New())
	// Offsets per source, added to the score of each candidate.
	b.AddFilter(sorterglobal.New(map[string]int{"lsp": 50}))

	b.PatchGlobal(cmp.Config{
		Sources:     []string{"lsp", "around", "buffer"},
		PostFilters: []string{"sorterglobal"},
		SourceOptions: map[string]cmp.SourceOption{
			"_": {
				Matchers: []string{"matcherfuzzy"},
				Sorters:  []string{"sorterfuzzy"},
			},
			"around": {Matchers: []string{"matcherfuzzy", "matchernottyped"}},
			"buffer": {Matchers: []string{"matcherfuzzy", "matchernottyped"}},
		},
	})

	b.Main()
}
```

`sorterfuzzy` keeps the order of candidates with the same score and length,
so a sorter before it (`sortertext`, for example) breaks the remaining ties:
`Sorters: []string{"sortertext", "sorterfuzzy"}`.

A source can give each candidate an offset of its own: when the `UserData` of
a candidate implements `sorterglobal.ScoreOffsetter`, `sorterglobal` adds its
`ScoreOffset()` too (like `score_offset` of a blink.cmp item).

```go
type data struct{ offset int }

func (d data) ScoreOffset() int { return d.offset }
```

`fuzzy.Match` returns the positions of the matched characters, for a UI that
highlights them.

## License

MIT, see LICENSE.
