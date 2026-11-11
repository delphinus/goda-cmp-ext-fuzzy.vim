// Package fuzzy is the fuzzy matching shared by the filters of this module
// (ext/filter/matcherfuzzy, ext/filter/sorterfuzzy, ext/filter/sorterglobal).
// UIs can use Match to highlight the matched characters.
//
// Matching and scoring follow fzf's v1 algorithm in a simple form, and unlike
// blink.cmp (frizbee), typos are not tolerated. Case is treated with
// smartcase, as matcherhead does.
package fuzzy

import (
	"strings"
	"unicode"

	"github.com/Shougo/goda-cmp.vim/cmp"
)

// FoldInput reports whether case is ignored for the input: when the ignoreCase
// option of the source is set, or when the typed text has no upper case letter.
func FoldInput(in cmp.FilterInput) bool {
	return in.IgnoreCase || Fold(in.Complete)
}

// Fold is smartcase: it reports whether case is ignored, that is, whether typed
// has no upper case letter.
func Fold(typed string) bool {
	return !strings.ContainsFunc(typed, unicode.IsUpper)
}

const (
	scoreMatch       = 16
	bonusBoundary    = 10 // at the head of a word (the start, after a delimiter, a camelCase upper case letter)
	bonusFirst       = 8  // at the start of the candidate
	bonusConsecutive = 4  // right after the previous matched character
	penaltyGap       = 1  // for each character between two matched characters
)

func isDelimiter(r rune) bool {
	return !unicode.IsLetter(r) && !unicode.IsDigit(r)
}

func boundary(text []rune, i int) bool {
	if i == 0 {
		return true
	}
	prev, cur := text[i-1], text[i]
	return isDelimiter(prev) && !isDelimiter(cur) || unicode.IsLower(prev) && unicode.IsUpper(cur)
}

// Score reports whether text contains the characters of pat in order, and
// with what score.
func Score(text, pat []rune, fold bool) (int, bool) {
	score, _, ok := Match(text, pat, fold)
	return score, ok
}

// Match is Score that also returns the positions (indexes of runes in text) of
// the matched characters, to highlight them in a menu. It matches greedily
// forward to find the end, then backward from there to narrow the start, as
// fzf's v1 does.
func Match(text, pat []rune, fold bool) (int, []int, bool) {
	if len(pat) == 0 {
		return 0, nil, true
	}
	eq := func(a, b rune) bool {
		if fold {
			return unicode.ToLower(a) == unicode.ToLower(b)
		}
		return a == b
	}
	// Forward: where the last character matches.
	pi, end := 0, -1
	for i, r := range text {
		if eq(r, pat[pi]) {
			pi++
			if pi == len(pat) {
				end = i
				break
			}
		}
	}
	if end < 0 {
		return 0, nil, false
	}
	// Backward: move the start of the match to the right.
	pi, start := len(pat)-1, end
	for i := end; i >= 0; i-- {
		if eq(text[i], pat[pi]) {
			pi--
			if pi < 0 {
				start = i
				break
			}
		}
	}
	// Match forward again within the narrowed range, and score it.
	score, pi, last := 0, 0, -1
	pos := make([]int, 0, len(pat))
	for i := start; i <= end && pi < len(pat); i++ {
		if !eq(text[i], pat[pi]) {
			continue
		}
		score += scoreMatch
		if boundary(text, i) {
			score += bonusBoundary
		}
		if i == 0 {
			score += bonusFirst
		}
		if last >= 0 {
			if i == last+1 {
				score += bonusConsecutive
			} else {
				score -= penaltyGap * (i - last - 1)
			}
		}
		last = i
		pos = append(pos, i)
		pi++
	}
	// A later start lowers the score a little (capped, so that a prefix like
	// nvim_buf_ does not sink a candidate much), and so does a longer
	// candidate.
	score -= min(start, 3) + (len(text)-len(pat))/8
	return score, pos, true
}
