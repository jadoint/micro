package words_test

import (
	"testing"

	"github.com/jadoint/micro/pkg/words"
)

func TestCount(t *testing.T) {
	tables := []struct {
		input string
		want  int
	}{
		{"Lorem Ipsum Dolor Sit Amet", 5},
		{"Lorem Ipsum Dolor", 3},
		{"Je m'appelle", 2},
		{"Enchanté", 1},
	}

	for _, table := range tables {
		got := words.Count(&table.input)
		if got != table.want {
			t.Errorf(`Count failed on "%s", got: %d, want: %d`, table.input, got, table.want)
		}
	}
}
