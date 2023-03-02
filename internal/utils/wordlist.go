// @todo add custom wordlist support
package utils

import (
	"embed"
	"strings"
)

// The default wordlist is embedded into the binary
//
//go:embed assets/wordlist.txt
var defaultWordlist embed.FS

var Wordlist *[]string

// Load the default wordlist at init, to read it only once
func init() {
	Wordlist = loadDefaultWordlist()
}

func loadDefaultWordlist() *[]string {
	data, _ := defaultWordlist.ReadFile("assets/wordlist.txt")
	list := strings.Split(string(data), "\n")
	if len(list) == 0 {
		panic("Wordlist is empty")
	}
	if list[len(list)-1] == "" {
		list = list[:len(list)-1]
	}
	return &list
}
