package function_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/berquerant/dql/data"
	"github.com/berquerant/dql/function"
	"github.com/berquerant/gogrep"
	"github.com/stretchr/testify/assert"
)

func TestGrep(t *testing.T) {
	f, err := os.CreateTemp("", "greptest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	lines := []string{
		"internet",
		"interpolation",
		"internationalization",
		"information",
		"initiation",
		"imitation",
		"intonation",
		"illustration",
		"irritation",
		"indication",
	}
	for _, line := range lines {
		if _, err := fmt.Fprintln(f, line); err != nil {
			t.Fatal(err)
		}
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	for _, tc := range []*struct {
		title   string
		pattern string
		want    int
	}{
		{
			title:   "no hits",
			pattern: "function",
		},
		{
			title:   "in-tion",
			pattern: "in.+tion",
			want:    6,
		},
		{
			title:   "liza",
			pattern: "liza",
			want:    1,
		},
	} {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			got, err := function.NewGrep(gogrep.New()).Call(data.FromString(tc.pattern), data.FromString(f.Name()))
			assert.Nil(t, err)
			assert.Equal(t, data.TypeInt, got.Type())
			assert.Equal(t, tc.want, got.Int())
		})
	}
}
