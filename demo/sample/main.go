package main

import (
	"fmt"
	"strings"
)

func main() {
	sentence := "  the quick brown fox  "
	words := strings.Fields(sentence)
	fmt.Println(len(words))
}
