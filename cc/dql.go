// Code generated by goyacc -o cc/dql.go -v cc/dql.output cc/dql.y. DO NOT EDIT.

//line cc/dql.y:2

package cc

import __yyfmt__ "fmt"

//line cc/dql.y:3

import (
	"github.com/berquerant/dql/ast"
	"github.com/berquerant/dql/token"
)

//line cc/dql.y:11
type yySymType struct {
	yys            int
	token          token.Token
	flag           bool
	exprs          *ast.Exprs
	simpleExpr     ast.SimpleExpr
	lit            ast.Lit
	expr           ast.Expr
	boolPrimary    ast.BoolPrimary
	bitExpr        ast.BitExpr
	predicate      ast.Predicate
	intLit         *ast.IntLit
	limitSection   *ast.LimitSection
	orderByTerm    *ast.OrderByTerm
	orderByTerms   *ast.OrderByTerms
	orderBySection *ast.OrderBySection
	havingSection  *ast.HavingSection
	groupByTerm    *ast.GroupByTerm
	groupByTerms   *ast.GroupByTerms
	groupBySection *ast.GroupBySection
	whereCondition *ast.WhereCondition
	whereSection   *ast.WhereSection
	selectOption   *ast.SelectOption
	ident          *ast.Ident
	selectTarget   *ast.SelectTarget
	selectTerm     *ast.SelectTerm
	selectTerms    *ast.SelectTerms
	selectSection  *ast.SelectSection
	statement      *ast.Statement
}

const SELECT = 57346
const DISTINCT = 57347
const WHERE = 57348
const HAVING = 57349
const GROUP = 57350
const BY = 57351
const ORDER = 57352
const LIMIT = 57353
const IDENT = 57354
const INT = 57355
const FLOAT = 57356
const STRING = 57357
const AS = 57358
const ASC = 57359
const DESC = 57360
const LIKE = 57361
const IN = 57362
const COMMA = 57363
const SCOLON = 57364
const LPAR = 57365
const RPAR = 57366
const PLUS = 57367
const MINUS = 57368
const AST = 57369
const SLASH = 57370
const NOT = 57371
const AND = 57372
const OR = 57373
const XOR = 57374
const EQ = 57375
const NE = 57376
const GT = 57377
const GQ = 57378
const LT = 57379
const LQ = 57380
const BETWEEN = 57381
const OFFSET = 57382
const AMP = 57383
const PIPE = 57384
const HAT = 57385
const TILDE = 57386

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"SELECT",
	"DISTINCT",
	"WHERE",
	"HAVING",
	"GROUP",
	"BY",
	"ORDER",
	"LIMIT",
	"IDENT",
	"INT",
	"FLOAT",
	"STRING",
	"AS",
	"ASC",
	"DESC",
	"LIKE",
	"IN",
	"COMMA",
	"SCOLON",
	"LPAR",
	"RPAR",
	"PLUS",
	"MINUS",
	"AST",
	"SLASH",
	"NOT",
	"AND",
	"OR",
	"XOR",
	"EQ",
	"NE",
	"GT",
	"GQ",
	"LT",
	"LQ",
	"BETWEEN",
	"OFFSET",
	"AMP",
	"PIPE",
	"HAT",
	"TILDE",
	"$left",
}

var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line cc/dql.y:468

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 17,
	19, 41,
	20, 41,
	39, 41,
	-2, 54,
}

const yyPrivate = 57344

const yyLast = 157

var yyAct = [...]int{
	13, 98, 16, 82, 18, 31, 70, 17, 105, 33,
	12, 78, 76, 39, 84, 40, 19, 27, 28, 29,
	38, 37, 39, 38, 62, 39, 60, 23, 91, 24,
	25, 77, 94, 14, 106, 59, 5, 71, 72, 73,
	74, 38, 37, 39, 75, 95, 68, 95, 26, 90,
	110, 63, 19, 27, 28, 29, 86, 79, 80, 67,
	83, 112, 65, 23, 33, 24, 25, 96, 88, 85,
	55, 56, 57, 58, 89, 102, 42, 43, 44, 45,
	46, 47, 36, 93, 26, 92, 52, 53, 54, 11,
	99, 71, 83, 9, 35, 101, 103, 100, 19, 27,
	28, 29, 4, 7, 21, 111, 50, 99, 113, 23,
	49, 24, 25, 41, 81, 61, 55, 56, 57, 58,
	51, 55, 56, 57, 58, 20, 22, 108, 109, 15,
	26, 107, 52, 53, 54, 48, 104, 52, 53, 54,
	38, 37, 39, 87, 97, 66, 34, 69, 10, 6,
	8, 64, 32, 30, 3, 1, 2,
}

var yyPact = [...]int{
	98, -1000, 14, 97, 88, -1000, 81, 4, 4, -1000,
	87, 73, -1000, 11, 40, 43, -1000, 91, -1000, 12,
	-1000, 86, -1000, 4, -1000, -1000, -1000, -1000, -1000, -1000,
	30, -1000, 46, 11, 49, 4, 4, 4, 4, 4,
	-1000, 86, -1000, -1000, -1000, -1000, -1000, -1000, -8, 86,
	86, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 4,
	-1000, -1000, -10, 4, -1000, 44, 57, 65, -1000, 28,
	-1000, 11, -7, -19, -1000, -1000, 5, 86, 86, 96,
	96, 8, 24, 11, -1000, -1000, -1000, -1000, 54, 4,
	4, 4, 45, -1000, -1000, 4, -32, 13, -1000, 110,
	-1000, 26, 86, 11, -1000, 48, 4, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000,
}

var yyPgo = [...]int{
	0, 156, 155, 154, 153, 5, 152, 151, 150, 149,
	10, 148, 147, 6, 146, 145, 144, 1, 143, 136,
	0, 135, 131, 129, 2, 7, 126, 125, 4, 114,
	3, 113, 110, 106, 104,
}

var yyR1 = [...]int{
	0, 2, 1, 3, 4, 4, 5, 6, 7, 7,
	8, 8, 9, 9, 10, 11, 11, 12, 12, 13,
	14, 14, 15, 15, 16, 16, 17, 22, 22, 22,
	18, 18, 19, 19, 30, 30, 20, 20, 20, 20,
	20, 21, 21, 23, 23, 31, 31, 31, 31, 31,
	31, 24, 24, 24, 24, 25, 25, 25, 33, 33,
	33, 33, 32, 32, 32, 28, 28, 28, 28, 28,
	34, 34, 34, 34, 26, 26, 26, 27, 29, 29,
}

var yyR2 = [...]int{
	0, 2, 6, 3, 1, 3, 2, 1, 0, 2,
	0, 1, 0, 2, 1, 0, 3, 1, 3, 1,
	0, 2, 0, 3, 1, 3, 2, 0, 1, 1,
	0, 3, 0, 2, 1, 3, 3, 3, 3, 2,
	1, 0, 1, 3, 1, 1, 1, 1, 1, 1,
	1, 6, 6, 4, 1, 3, 3, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 2, 1, 3,
	1, 1, 1, 1, 1, 1, 1, 4, 0, 1,
}

var yyChk = [...]int{
	-1000, -2, -1, -3, 4, 22, -9, 6, -8, 5,
	-11, 8, -10, -20, 29, -23, -24, -25, -28, 12,
	-27, -34, -26, 23, 25, 26, 44, 13, 14, 15,
	-4, -5, -6, -20, -14, 7, 9, 31, 30, 32,
	-20, -31, 33, 34, 35, 36, 37, 38, -21, -32,
	-33, 29, 41, 42, 43, 25, 26, 27, 28, 23,
	-28, 29, -20, 21, -7, 16, -15, 10, -10, -12,
	-13, -20, -20, -20, -20, -24, 20, 39, 19, -25,
	-25, -29, -30, -20, 24, -5, 12, -18, 11, 9,
	21, 23, -25, -28, 24, 21, 13, -16, -17, -20,
	-13, -30, 30, -20, -19, 40, 21, -22, 17, 18,
	24, -24, 13, -17,
}

var yyDef = [...]int{
	0, -2, 0, 12, 10, 1, 15, 0, 0, 11,
	20, 0, 13, 14, 73, 40, 44, -2, 57, 65,
	66, 0, 68, 0, 70, 71, 72, 74, 75, 76,
	3, 4, 8, 7, 22, 0, 0, 0, 0, 0,
	39, 0, 45, 46, 47, 48, 49, 50, 0, 0,
	0, 42, 62, 63, 64, 58, 59, 60, 61, 78,
	67, 73, 0, 0, 6, 0, 30, 0, 21, 16,
	17, 19, 36, 37, 38, 43, 0, 0, 0, 55,
	56, 0, 79, 34, 69, 5, 9, 2, 0, 0,
	0, 0, 0, 53, 77, 0, 32, 23, 24, 27,
	18, 0, 0, 35, 31, 0, 0, 26, 28, 29,
	51, 52, 33, 25,
}

var yyTok1 = [...]int{
	1,
}

var yyTok2 = [...]int{
	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45,
}

var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ??, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-2 : yypt+1]
//line cc/dql.y:139
		{
			yylex.(Lexer).SetResult(yyDollar[1].statement)
			yyVAL.statement = yyDollar[1].statement
		}
	case 2:
		yyDollar = yyS[yypt-6 : yypt+1]
//line cc/dql.y:150
		{
			yyVAL.statement = &ast.Statement{
				SelectSection:  yyDollar[1].selectSection,
				WhereSection:   yyDollar[2].whereSection,
				GroupBySection: yyDollar[3].groupBySection,
				HavingSection:  yyDollar[4].havingSection,
				OrderBySection: yyDollar[5].orderBySection,
				LimitSection:   yyDollar[6].limitSection,
			}
		}
	case 3:
		yyDollar = yyS[yypt-3 : yypt+1]
//line cc/dql.y:164
		{
			yyVAL.selectSection = &ast.SelectSection{
				Option: yyDollar[2].selectOption,
				Terms:  yyDollar[3].selectTerms,
			}
		}
	case 4:
		yyDollar = yyS[yypt-1 : yypt+1]
//line cc/dql.y:172
		{
			yyVAL.selectTerms = &ast.SelectTerms{Terms: []*ast.SelectTerm{yyDollar[1].selectTerm}}
		}
	case 5:
		yyDollar = yyS[yypt-3 : yypt+1]
//line cc/dql.y:175
		{
			v := append(yyDollar[1].selectTerms.Terms, yyDollar[3].selectTerm)
			yyVAL.selectTerms = &ast.SelectTerms{Terms: v}
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
//line cc/dql.y:181
		{
			yyVAL.selectTerm = &ast.SelectTerm{
				Target: yyDollar[1].selectTarget,
				As:     yyDollar[2].ident,
			}
		}
	case 7:
		yyDollar = yyS[yypt-1 : yypt+1]
//line cc/dql.y:189
		{
			yyVAL.selectTarget = &ast.SelectTarget{Expr: yyDollar[1].expr}
		}
	case 8:
		yyDollar = yyS[yypt-0 : yypt+1]
//line cc/dql.y:194
		{
			yyVAL.ident = nil
		}
	case 9:
		yyDollar = yyS[yypt-2 : yypt+1]
//line cc/dql.y:197
		{
			yyVAL.ident = &ast.Ident{Value: yyDollar[2].token.Value()}
		}
	case 10:
		yyDollar = yyS[yypt-0 : yypt+1]
//line cc/dql.y:202
		{
			yyVAL.selectOption = nil
		}
	case 11:
		yyDollar = yyS[yypt-1 : yypt+1]
//line cc/dql.y:205
		{
			yyVAL.selectOption = &ast.SelectOption{IsDistinct: true}
		}
	case 12:
		yyDollar = yyS[yypt-0 : yypt+1]
//line cc/dql.y:210
		{
			yyVAL.whereSection = nil
		}
	case 13:
		yyDollar = yyS[yypt-2 : yypt+1]
//line cc/dql.y:213
		{
			yyVAL.whereSection = &ast.WhereSection{Condition: yyDollar[2].whereCondition}
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
//line cc/dql.y:218
		{
			yyVAL.whereCondition = &ast.WhereCondition{Expr: yyDollar[1].expr}
		}
	case 15:
		yyDollar = yyS[yypt-0 : yypt+1]
//line cc/dql.y:223
		{
			yyVAL.groupBySection = nil
		}
	case 16:
		yyDollar = yyS[yypt-3 : yypt+1]
//line cc/dql.y:226
		{
			yyVAL.groupBySection = &ast.GroupBySection{Terms: yyDollar[3].groupByTerms}
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
//line cc/dql.y:231
		{
			yyVAL.groupByTerms = &ast.GroupByTerms{Terms: []*ast.GroupByTerm{yyDollar[1].groupByTerm}}
		}
	case 18:
		yyDollar = yyS[yypt-3 : yypt+1]
//line cc/dql.y:234
		{
			v := append(yyDollar[1].groupByTerms.Terms, yyDollar[3].groupByTerm)
			yyVAL.groupByTerms = &ast.GroupByTerms{Terms: v}
		}
	case 19:
		yyDollar = yyS[yypt-1 : yypt+1]
//line cc/dql.y:240
		{
			yyVAL.groupByTerm = &ast.GroupByTerm{Expr: yyDollar[1].expr}
		}
	case 20:
		yyDollar = yyS[yypt-0 : yypt+1]
//line cc/dql.y:245
		{
			yyVAL.havingSection = nil
		}
	case 21:
		yyDollar = yyS[yypt-2 : yypt+1]
//line cc/dql.y:248
		{
			yyVAL.havingSection = &ast.HavingSection{Condition: yyDollar[2].whereCondition}
		}
	case 22:
		yyDollar = yyS[yypt-0 : yypt+1]
//line cc/dql.y:253
		{
			yyVAL.orderBySection = nil
		}
	case 23:
		yyDollar = yyS[yypt-3 : yypt+1]
//line cc/dql.y:256
		{
			yyVAL.orderBySection = &ast.OrderBySection{Terms: yyDollar[3].orderByTerms}
		}
	case 24:
		yyDollar = yyS[yypt-1 : yypt+1]
//line cc/dql.y:261
		{
			yyVAL.orderByTerms = &ast.OrderByTerms{Terms: []*ast.OrderByTerm{yyDollar[1].orderByTerm}}
		}
	case 25:
		yyDollar = yyS[yypt-3 : yypt+1]
//line cc/dql.y:264
		{
			v := append(yyDollar[1].orderByTerms.Terms, yyDollar[3].orderByTerm)
			yyVAL.orderByTerms = &ast.OrderByTerms{Terms: v}
		}
	case 26:
		yyDollar = yyS[yypt-2 : yypt+1]
//line cc/dql.y:270
		{
			opt := &ast.OrderByTermOption{
				IsDesc: yyDollar[2].flag,
			}
			yyVAL.orderByTerm = &ast.OrderByTerm{
				Expr:   yyDollar[1].expr,
				Option: opt,
			}
		}
	case 27:
		yyDollar = yyS[yypt-0 : yypt+1]
//line cc/dql.y:281
		{
			yyVAL.flag = false
		}
	case 28:
		yyDollar = yyS[yypt-1 : yypt+1]
//line cc/dql.y:284
		{
			yyVAL.flag = false
		}
	case 29:
		yyDollar = yyS[yypt-1 : yypt+1]
//line cc/dql.y:287
		{
			yyVAL.flag = true
		}
	case 30:
		yyDollar = yyS[yypt-0 : yypt+1]
//line cc/dql.y:292
		{
			yyVAL.limitSection = nil
		}
	case 31:
		yyDollar = yyS[yypt-3 : yypt+1]
//line cc/dql.y:295
		{
			l := yylex.(Lexer)
			v := l.ParseInt(yyDollar[2].token.Value())
			yyVAL.limitSection = &ast.LimitSection{
				Limit:  &ast.IntLit{Value: v},
				Offset: yyDollar[3].intLit,
			}
		}
	case 32:
		yyDollar = yyS[yypt-0 : yypt+1]
//line cc/dql.y:305
		{
			yyVAL.intLit = nil
		}
	case 33:
		yyDollar = yyS[yypt-2 : yypt+1]
//line cc/dql.y:308
		{
			l := yylex.(Lexer)
			v := l.ParseInt(yyDollar[2].token.Value())
			yyVAL.intLit = &ast.IntLit{Value: v}
		}
	case 34:
		yyDollar = yyS[yypt-1 : yypt+1]
//line cc/dql.y:315
		{
			yyVAL.exprs = &ast.Exprs{Exprs: []ast.Expr{yyDollar[1].expr}}
		}
	case 35:
		yyDollar = yyS[yypt-3 : yypt+1]
//line cc/dql.y:318
		{
			v := append(yyDollar[1].exprs.Exprs, yyDollar[3].expr)
			yyVAL.exprs = &ast.Exprs{Exprs: v}
		}
	case 36:
		yyDollar = yyS[yypt-3 : yypt+1]
//line cc/dql.y:324
		{
			yyVAL.expr = &ast.OrExpr{Left: yyDollar[1].expr, Right: yyDollar[3].expr}
		}
	case 37:
		yyDollar = yyS[yypt-3 : yypt+1]
//line cc/dql.y:327
		{
			yyVAL.expr = &ast.AndExpr{Left: yyDollar[1].expr, Right: yyDollar[3].expr}
		}
	case 38:
		yyDollar = yyS[yypt-3 : yypt+1]
//line cc/dql.y:330
		{
			yyVAL.expr = &ast.XorExpr{Left: yyDollar[1].expr, Right: yyDollar[3].expr}
		}
	case 39:
		yyDollar = yyS[yypt-2 : yypt+1]
//line cc/dql.y:333
		{
			yyVAL.expr = &ast.NotExpr{Expr: yyDollar[2].expr}
		}
	case 40:
		yyDollar = yyS[yypt-1 : yypt+1]
//line cc/dql.y:336
		{
			yyVAL.expr = yyDollar[1].boolPrimary
		}
	case 41:
		yyDollar = yyS[yypt-0 : yypt+1]
//line cc/dql.y:341
		{
			yyVAL.flag = false
		}
	case 42:
		yyDollar = yyS[yypt-1 : yypt+1]
//line cc/dql.y:344
		{
			yyVAL.flag = true
		}
	case 43:
		yyDollar = yyS[yypt-3 : yypt+1]
//line cc/dql.y:349
		{
			l := yylex.(Lexer)
			op := l.AsComparisonType(yyDollar[2].token.Type())
			yyVAL.boolPrimary = &ast.BoolPrimaryComparison{
				Op:    op,
				Left:  yyDollar[1].boolPrimary,
				Right: yyDollar[3].predicate,
			}
		}
	case 44:
		yyDollar = yyS[yypt-1 : yypt+1]
//line cc/dql.y:358
		{
			yyVAL.boolPrimary = &ast.BoolPrimaryPredicate{Pred: yyDollar[1].predicate}
		}
	case 51:
		yyDollar = yyS[yypt-6 : yypt+1]
//line cc/dql.y:366
		{
			yyVAL.predicate = &ast.PredicateIn{
				IsNot:  yyDollar[2].flag,
				Target: yyDollar[1].bitExpr,
				List:   yyDollar[5].exprs,
			}
		}
	case 52:
		yyDollar = yyS[yypt-6 : yypt+1]
//line cc/dql.y:373
		{
			yyVAL.predicate = &ast.PredicateBetween{
				IsNot:  yyDollar[2].flag,
				Target: yyDollar[1].bitExpr,
				Left:   yyDollar[4].bitExpr,
				Right:  yyDollar[6].predicate,
			}
		}
	case 53:
		yyDollar = yyS[yypt-4 : yypt+1]
//line cc/dql.y:381
		{
			yyVAL.predicate = &ast.PredicateLike{
				IsNot:   yyDollar[2].flag,
				Target:  yyDollar[1].bitExpr,
				Pattern: yyDollar[4].simpleExpr,
			}
		}
	case 54:
		yyDollar = yyS[yypt-1 : yypt+1]
//line cc/dql.y:388
		{
			yyVAL.predicate = &ast.PredicateBitExpr{Expr: yyDollar[1].bitExpr}
		}
	case 55:
		yyDollar = yyS[yypt-3 : yypt+1]
//line cc/dql.y:393
		{
			l := yylex.(Lexer)
			op := l.AsBitOperatorType(yyDollar[2].token.Type())
			yyVAL.bitExpr = &ast.BitExprBitOp{Op: op, Left: yyDollar[1].bitExpr, Right: yyDollar[3].bitExpr}
		}
	case 56:
		yyDollar = yyS[yypt-3 : yypt+1]
//line cc/dql.y:398
		{
			l := yylex.(Lexer)
			op := l.AsArithmeticOperatorType(yyDollar[2].token.Type())
			yyVAL.bitExpr = &ast.BitExprArtOp{Op: op, Left: yyDollar[1].bitExpr, Right: yyDollar[3].bitExpr}
		}
	case 57:
		yyDollar = yyS[yypt-1 : yypt+1]
//line cc/dql.y:403
		{
			yyVAL.bitExpr = &ast.BitExprSimpleExpr{Expr: yyDollar[1].simpleExpr}
		}
	case 65:
		yyDollar = yyS[yypt-1 : yypt+1]
//line cc/dql.y:414
		{
			yyVAL.simpleExpr = &ast.Ident{Value: yyDollar[1].token.Value()}
		}
	case 66:
		yyDollar = yyS[yypt-1 : yypt+1]
//line cc/dql.y:417
		{
			yyVAL.simpleExpr = yyDollar[1].simpleExpr
		}
	case 67:
		yyDollar = yyS[yypt-2 : yypt+1]
//line cc/dql.y:420
		{
			l := yylex.(Lexer)
			op := l.AsPrefixOperatorType(yyDollar[1].token.Type())
			yyVAL.simpleExpr = &ast.SimpleExprPrefixOp{Op: op, Expr: yyDollar[2].simpleExpr}
		}
	case 68:
		yyDollar = yyS[yypt-1 : yypt+1]
//line cc/dql.y:425
		{
			yyVAL.simpleExpr = &ast.SimpleExprLit{Lit: yyDollar[1].lit}
		}
	case 69:
		yyDollar = yyS[yypt-3 : yypt+1]
//line cc/dql.y:428
		{
			yyVAL.simpleExpr = &ast.SimpleExprExpr{Expr: yyDollar[2].expr}
		}
	case 74:
		yyDollar = yyS[yypt-1 : yypt+1]
//line cc/dql.y:436
		{
			l := yylex.(Lexer)
			v := l.ParseInt(yyDollar[1].token.Value())
			yyVAL.lit = &ast.IntLit{Value: v}
		}
	case 75:
		yyDollar = yyS[yypt-1 : yypt+1]
//line cc/dql.y:441
		{
			l := yylex.(Lexer)
			v := l.ParseFloat(yyDollar[1].token.Value())
			yyVAL.lit = &ast.FloatLit{Value: v}
		}
	case 76:
		yyDollar = yyS[yypt-1 : yypt+1]
//line cc/dql.y:446
		{
			v := yyDollar[1].token.Value()
			yyVAL.lit = &ast.StringLit{Value: v}
		}
	case 77:
		yyDollar = yyS[yypt-4 : yypt+1]
//line cc/dql.y:452
		{
			name := &ast.Ident{Value: yyDollar[1].token.Value()}
			yyVAL.simpleExpr = &ast.FunctionCall{
				FunctionName: name,
				Arguments:    yyDollar[3].exprs,
			}
		}
	case 78:
		yyDollar = yyS[yypt-0 : yypt+1]
//line cc/dql.y:461
		{
			yyVAL.exprs = &ast.Exprs{}
		}
	case 79:
		yyDollar = yyS[yypt-1 : yypt+1]
//line cc/dql.y:464
		{
			yyVAL.exprs = yyDollar[1].exprs
		}
	}
	goto yystack /* stack new state and value */
}
