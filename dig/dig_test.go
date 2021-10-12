package dig_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/berquerant/dql/dig"
	"github.com/stretchr/testify/assert"
)

func TestDigger(t *testing.T) {
	root := filepath.Join(os.Getenv("ROOT"), "dig", "testdata")
	withRoot := func(p string) string { return filepath.Join(root, p) }

	for _, tc := range []*struct {
		title   string
		target  string // relative path from testdata/
		handler dig.FileInfoHandler
		want    []string // names, relative path from testdata/
		isErr   bool
	}{
		{
			title:  "not exist",
			target: "not_exist",
			handler: func(_ dig.FileInfo) dig.Instr {
				return dig.InstrContinue
			},
			isErr: true,
		},
		{
			title:  "file",
			target: "a.log",
			handler: func(_ dig.FileInfo) dig.Instr {
				return dig.InstrContinue
			},
			want: []string{
				"a.log",
			},
		},
		{
			title:  "visit all",
			target: ".",
			handler: func(_ dig.FileInfo) dig.Instr {
				return dig.InstrContinue
			},
			want: []string{
				".",
				"a.log",
				"dir",
				"dir/b.log",
				"dir2",
				"dir2/c.log",
				"dir2/d.log",
			},
		},
		{
			title:  "skip",
			target: ".",
			handler: func(info dig.FileInfo) dig.Instr {
				if info.Name() == withRoot("dir") {
					return dig.InstrSkipDir
				}
				return dig.InstrContinue
			},
			want: []string{
				".",
				"a.log",
				"dir",
				"dir2",
				"dir2/c.log",
				"dir2/d.log",
			},
		},
		{
			title:  "cancel",
			target: ".",
			handler: func(info dig.FileInfo) dig.Instr {
				if info.Name() == withRoot("dir") {
					return dig.InstrCancel
				}
				return dig.InstrContinue
			},
			want: []string{
				".",
				"a.log",
				"dir",
			},
		},
	} {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			got := []string{}
			err := dig.New().Dig(withRoot(tc.target), func(info dig.FileInfo) dig.Instr {
				instr := tc.handler(info)
				got = append(got, info.Name())
				return instr
			})
			if tc.isErr {
				if err == nil {
					t.Fatal("want error")
				}
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, len(tc.want), len(got))
			for i, x := range tc.want {
				w := withRoot(x)
				g := got[i]
				assert.Equal(t, w, g)
			}
		})
	}
}
