package utils

import (
	"testing"
)

func TestLoadDefaultWordlist(t *testing.T) {
	list := loadDefaultWordlist()
	if len(*list) == 0 {
		t.Error("Wordlist is empty")
	}
	if (*list)[0] != "user" {
		t.Error("First word in wordlist is not 'user'")
	}
	if (*list)[len(*list)-1] == "" {
		t.Error("Last word in wordlist is empty")
	}
}
