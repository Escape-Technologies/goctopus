package utils

import (
	"fmt"
	"testing"
)

func TestLoadDefaultWordlist(t *testing.T) {
	list := loadDefaultWordlist()
	if len(*list) == 0 {
		t.Error("Wordlist is empty")
	}
	if (*list)[0] != "the" {
		t.Error("First word in wordlist is not 'the'")
	}
	if (*list)[len(*list)-1] == "" {
		t.Error("Last word in wordlist is empty")
	}
	fmt.Printf("Last word in wordlist: %s", (*list)[len(*list)-1])
}
