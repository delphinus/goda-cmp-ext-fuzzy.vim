# goda-cmp-ext-fuzzy.vim

[English](README.md)

> [!WARNING]
> β版 (v0.0.x) です。goda-cmp.vim 自体が開発の途中で、同梱の拡張を別のリポジトリへ移す作業も進んでいるため、API・import のパス・このリポジトリの名前が変わることがあります。

[goda-cmp.vim](https://github.com/Shougo/goda-cmp.vim) 用の fuzzy なフィルタ。

| 名前 | 種類 | 内容 |
|---|---|---|
| `matcherfuzzy` | マッチャ | 打った文字を、隣り合っていなくても順に含む候補を残す。`Prl` で `Println` が見つかる。 |
| `sorterfuzzy` | ソータ | 一致の良さで候補を並べる。連続した文字、語の頭 (区切り文字の後、camelCase の大文字)、候補の先頭での一致ほど点が高い。 |
| `sorterglobal` | post filter | 全部のソースの候補をまとめて一致の良さで並べる。ソースごと・候補ごとの加点も足す。 |
| `matchernottyped` | マッチャ | 打っている語そのものの候補を外す (`around` と `buffer` はカーソルのある語も拾うため)。 |

大文字と小文字は `matcherhead` と同じく smartcase で扱う。打った文字に大文字が無ければ区別しない。ソースの `ignoreCase` オプションを指定すると、常に区別しない。

## なぜ作ったか

goda-cmp 同梱のマッチャは語頭の一致 (`matcherhead`、`matcherprefix`) だけを扱う。また goda-cmp は、ソースごとに並べた候補をソースの順に繋ぐだけなので、最初のソースの弱い一致が、最後のソースのちょうどの一致より先に来る。これらのフィルタは、その両方で goda-cmp を blink.cmp に近い動きにする。

一致の判定と点数の付け方は、fzf の v1 のアルゴリズムを簡単にしたもの。打ち間違いは許容しない。

## インストール

これらは goda-cmp の設定が import する Go のパッケージ。設定の `go.mod` で require するか:

```sh
cd ~/.config/nvim/goda-cmp   # 設定のモジュール
go get github.com/delphinus/goda-cmp-ext-fuzzy.vim@latest
go mod tidy
```

プラグインマネージャでリポジトリを入れ、そのディレクトリを `g:goda#workspace` に足す。そうすると `goda_cmp#build()` が、入れたものでビルドする (goda-pm.vim はこれを自動で行う)。lazy.nvim なら例えば:

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

この場合も `go.mod` での require は要る。workspace が差し替えるのは、どのバージョンを使うかだけ。

## 設定

設定の `main.go` で:

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
	// ... UI とソース ...

	b.AddFilter(matcherfuzzy.New())
	b.AddFilter(sorterfuzzy.New())
	b.AddFilter(matchernottyped.New())
	// ソースごとの加点。各候補の点数に足される。
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

`sorterfuzzy` は、点数も長さも同じ候補の順を保つ。そのため、前に置いたソータ (例えば `sortertext`) が残りの同点の並びを決める: `Sorters: []string{"sortertext", "sorterfuzzy"}`。

ソースは候補ごとに加点を付けられる。候補の `UserData` が `sorterglobal.ScoreOffsetter` を実装していれば、`sorterglobal` はその `ScoreOffset()` も足す (blink.cmp の候補の `score_offset` と同じ)。

```go
type data struct{ offset int }

func (d data) ScoreOffset() int { return d.offset }
```

`fuzzy.Match` は一致した文字の位置を返す。一致した文字に色を付ける UI で使える。

## ライセンス

MIT。LICENSE を参照。
