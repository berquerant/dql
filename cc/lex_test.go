package cc_test

import (
	"bytes"
	"testing"

	"github.com/berquerant/dql/cc"
	"github.com/berquerant/dql/token"
	"github.com/stretchr/testify/assert"
)

func TestLexer(t *testing.T) {
	for _, tc := range []*struct {
		title string
		input string
		want  []token.Token
	}{
		{
			title: "ident",
			input: "x",
			want:  []token.Token{token.New(cc.IDENT, "x")},
		},
		{
			title: "select",
			input: "select name, size;",
			want: []token.Token{
				token.New(cc.SELECT, "select"),
				token.New(cc.IDENT, "name"),
				token.New(cc.COMMA, ","),
				token.New(cc.IDENT, "size"),
				token.New(cc.SCOLON, ";"),
			},
		},
		{
			title: "where",
			input: "where size > 100",
			want: []token.Token{
				token.New(cc.WHERE, "where"),
				token.New(cc.IDENT, "size"),
				token.New(cc.GT, ">"),
				token.New(cc.INT, "100"),
			},
		},
		{
			title: "where with comparisons",
			input: "where p <> 'still' and x <= 1.2 or y >= -1",
			want: []token.Token{
				token.New(cc.WHERE, "where"),
				token.New(cc.IDENT, "p"),
				token.New(cc.NE, "<>"),
				token.New(cc.STRING, "still"),
				token.New(cc.AND, "and"),
				token.New(cc.IDENT, "x"),
				token.New(cc.LQ, "<="),
				token.New(cc.FLOAT, "1.2"),
				token.New(cc.OR, "or"),
				token.New(cc.IDENT, "y"),
				token.New(cc.GQ, ">="),
				token.New(cc.MINUS, "-"),
				token.New(cc.INT, "1"),
			},
		},
		{
			title: "group by",
			input: "group by title, size/1000 as kb",
			want: []token.Token{
				token.New(cc.GROUP, "group"),
				token.New(cc.BY, "by"),
				token.New(cc.IDENT, "title"),
				token.New(cc.COMMA, ","),
				token.New(cc.IDENT, "size"),
				token.New(cc.SLASH, "/"),
				token.New(cc.INT, "1000"),
				token.New(cc.AS, "as"),
				token.New(cc.IDENT, "kb"),
			},
		},
		{
			title: "having",
			input: "having owner in ('root', 'log')",
			want: []token.Token{
				token.New(cc.HAVING, "having"),
				token.New(cc.IDENT, "owner"),
				token.New(cc.IN, "in"),
				token.New(cc.LPAR, "("),
				token.New(cc.STRING, "root"),
				token.New(cc.COMMA, ","),
				token.New(cc.STRING, "log"),
				token.New(cc.RPAR, ")"),
			},
		},
		{
			title: "order by",
			input: "order by sum(size), name desc",
			want: []token.Token{
				token.New(cc.ORDER, "order"),
				token.New(cc.BY, "by"),
				token.New(cc.IDENT, "sum"),
				token.New(cc.LPAR, "("),
				token.New(cc.IDENT, "size"),
				token.New(cc.RPAR, ")"),
				token.New(cc.COMMA, ","),
				token.New(cc.IDENT, "name"),
				token.New(cc.DESC, "desc"),
			},
		},
		{
			title: "limit",
			input: "limit 10 offset 5",
			want: []token.Token{
				token.New(cc.LIMIT, "limit"),
				token.New(cc.INT, "10"),
				token.New(cc.OFFSET, "offset"),
				token.New(cc.INT, "5"),
			},
		},
		{
			title: "ugly",
			input: "SELECT size as Size,-   size As neG24   , Where  NORM( 1, 3,p)>0.5  ;",
			want: []token.Token{
				token.New(cc.SELECT, "SELECT"),
				token.New(cc.IDENT, "size"),
				token.New(cc.AS, "as"),
				token.New(cc.IDENT, "Size"),
				token.New(cc.COMMA, ","),
				token.New(cc.MINUS, "-"),
				token.New(cc.IDENT, "size"),
				token.New(cc.AS, "As"),
				token.New(cc.IDENT, "neG24"),
				token.New(cc.COMMA, ","),
				token.New(cc.WHERE, "Where"),
				token.New(cc.IDENT, "NORM"),
				token.New(cc.LPAR, "("),
				token.New(cc.INT, "1"),
				token.New(cc.COMMA, ","),
				token.New(cc.INT, "3"),
				token.New(cc.COMMA, ","),
				token.New(cc.IDENT, "p"),
				token.New(cc.RPAR, ")"),
				token.New(cc.GT, ">"),
				token.New(cc.FLOAT, "0.5"),
				token.New(cc.SCOLON, ";"),
			},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			buf := bytes.NewBufferString(tc.input)
			l := cc.NewLexer(buf)
			l.Debug(0)
			got := []token.Token{}
			for {
				tok := l.Scan()
				if tok == cc.EOF {
					break
				}
				got = append(got, token.New(tok, l.Buffer()))
				l.ResetBuffer()
			}
			if l.Err() != nil {
				t.Fatalf("lexer got error %v", l.Err())
				return
			}
			assert.Equal(t, len(tc.want), len(got))
			for i, w := range tc.want {
				assert.Equal(t, w.Type(), got[i].Type(), "at index %d", i)
				assert.Equal(t, w.Value(), got[i].Value(), "at index %d", i)
			}
		})
	}
}
