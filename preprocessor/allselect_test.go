package preprocessor_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/berquerant/dql/ast"
	"github.com/berquerant/dql/cc"
	"github.com/berquerant/dql/preprocessor"
	"github.com/stretchr/testify/assert"
)

func TestSelectAll(t *testing.T) {
	const allSymbol = "all"
	for _, tc := range []*struct {
		title   string
		input   string
		want    string
		isError bool
	}{
		{
			title: "no all",
			input: "select name;",
			want:  "select name;",
		},
		{
			title: "all",
			input: fmt.Sprintf("select %s;", allSymbol),
			want:  "select name, size, mode, mod_time, is_dir;",
		},
		{
			title: "with all",
			input: fmt.Sprintf("select name,%s,size;", allSymbol),
			want:  "select name, name, size, mode, mod_time, is_dir, size;",
		},
		{
			title:   "as all",
			input:   fmt.Sprintf("select %s as bronze;", allSymbol),
			isError: true,
		},
	} {
		tc := tc
		t.Run(tc.title, func(t *testing.T) {
			lexer := cc.NewLexer(strings.NewReader(tc.input))
			_ = cc.Parse(lexer)
			assert.Nil(t, lexer.Err())
			tree := lexer.Result().(*ast.Statement)
			s := preprocessor.NewSelectAll(allSymbol)
			if tc.isError {
				assert.NotNil(t, s.PreProcess(tree))
				return
			}
			assert.Nil(t, s.PreProcess(tree))
			assert.Equal(t, tc.want, tree.String())
		})
	}
}
