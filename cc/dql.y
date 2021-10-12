%{

package cc

import (
  "github.com/berquerant/dql/ast"
  "github.com/berquerant/dql/token"
)
%}

%union{
  token token.Token
  flag bool
  exprs *ast.Exprs
  simpleExpr ast.SimpleExpr
  lit ast.Lit
  expr ast.Expr
  boolPrimary ast.BoolPrimary
  bitExpr ast.BitExpr
  predicate ast.Predicate
  intLit *ast.IntLit
  limitSection *ast.LimitSection
  orderByTerm *ast.OrderByTerm
  orderByTerms *ast.OrderByTerms
  orderBySection *ast.OrderBySection
  havingSection *ast.HavingSection
  groupByTerm *ast.GroupByTerm
  groupByTerms *ast.GroupByTerms
  groupBySection *ast.GroupBySection
  whereCondition *ast.WhereCondition
  whereSection *ast.WhereSection
  selectOption *ast.SelectOption
  ident *ast.Ident
  selectTarget *ast.SelectTarget
  selectTerm *ast.SelectTerm
  selectTerms *ast.SelectTerms
  selectSection *ast.SelectSection
  statement *ast.Statement
}

%type <statement> statement top
%type <selectSection> select_section
%type <selectTerms> select_terms
%type <selectTerm> select_term
%type <selectTarget> select_target
%type <ident> select_as_term
%type <selectOption> select_option
%type <whereSection> where_section
%type <whereCondition> where_condition
%type <groupBySection> group_by_section
%type <groupByTerms> group_by_terms
%type <groupByTerm> group_by_term
%type <havingSection> having_section
%type <orderBySection> order_by_section
%type <orderByTerms> order_by_terms
%type <orderByTerm> order_by_term
%type <limitSection> limit_section
%type <intLit> limit_offset
%type <expr> expr
%type <flag> not_option order_by_term_option
%type <boolPrimary> bool_primary
%type <predicate> predicate
%type <bitExpr> bit_expr
%type <lit> literal
%type <simpleExpr> function_call simple_expr
%type <exprs> arg_list exprs

%type <token> comparison_operator
%type <token> bit_operator
%type <token> arithmetic_operator
%type <token> prefix_operator

%token <token> SELECT  /* select */
%token <token> DISTINCT  /* distinct */
%token <token> WHERE  /* where */
%token <token> HAVING  /* having */
%token <token> GROUP  /* group */
%token <token> BY  /* by */
%token <token> ORDER  /* order */
%token <token> LIMIT  /* limit */

%token <token> IDENT  /* identifier */
%token <token> INT  /* integer */
%token <token> FLOAT  /* floating point */
%token <token> STRING  /* string */

%token <token> AS  /* as */
%token <token> ASC  /* asc */
%token <token> DESC  /* desc */
%token <token> LIKE  /* like */

%token <token> IN  /* in */

%token <token> COMMA  /* , */
%token <token> SCOLON  /* ; */
%token <token> LPAR  /* ( */
%token <token> RPAR  /* ) */

%token <token> PLUS  /* + */
%token <token> MINUS  /* - */
%token <token> AST  /* * */
%token <token> SLASH  /* / */

%token <token> NOT  /* not */
%token <token> AND  /* and */
%token <token> OR  /* or */
%token <token> XOR  /* xor */

%token <token> EQ  /* = */
%token <token> NE  /* <> */
%token <token> GT  /* > */
%token <token> GQ  /* >= */
%token <token> LT  /* < */
%token <token> LQ  /* <= */

%token <token> BETWEEN  /* between */
%token <token> OFFSET  /* offset */

%token <token> AMP  /* & */
%token <token> PIPE  /* | */
%token <token> HAT  /* ^ */
%token <token> TILDE  /* ~ */

%left OR
%left AND
%left XOR
%left NOT
%left EQ NE GT GQ LT LQ
%left PLUS MINUS
%left AST SLASH
$left PIPE
%left AMP
%left HAT

%%

top:
  statement SCOLON
  {
    yylex.(Lexer).SetResult($1)
    $$ = $1
  }

statement:
  select_section
  where_section
  group_by_section
  having_section
  order_by_section
  limit_section {
    $$ = &ast.Statement{
      SelectSection: $1,
      WhereSection: $2,
      GroupBySection: $3,
      HavingSection: $4,
      OrderBySection: $5,
      LimitSection: $6,
    }
  }

select_section:
  SELECT
  select_option
  select_terms {
    $$ = &ast.SelectSection{
      Option: $2,
      Terms: $3,
    }
  }

select_terms:
  select_term {
    $$ = &ast.SelectTerms{Terms: []*ast.SelectTerm{$1}}
  }
  | select_terms COMMA select_term {
    v := append($1.Terms, $3)
    $$ = &ast.SelectTerms{Terms: v}
  }

select_term:
  select_target select_as_term {
    $$ = &ast.SelectTerm{
      Target: $1,
      As: $2,
    }
  }

select_target:
  expr {
    $$ = &ast.SelectTarget{Expr: $1}
  }

select_as_term:
  {
    $$ = nil
  }
  | AS IDENT {
    $$ = &ast.Ident{Value: $2.Value()}
  }

select_option:
  {
    $$ = nil
  }
  | DISTINCT {
    $$ = &ast.SelectOption{IsDistinct: true}
  }

where_section:
  {
    $$ = nil
  }
  | WHERE where_condition {
    $$ = &ast.WhereSection{Condition: $2}
  }

where_condition:
  expr {
    $$ = &ast.WhereCondition{Expr: $1}
  }

group_by_section:
  {
    $$ = nil
  }
  | GROUP BY group_by_terms {
    $$ = &ast.GroupBySection{Terms: $3}
  }

group_by_terms:
  group_by_term {
    $$ = &ast.GroupByTerms{Terms: []*ast.GroupByTerm{$1}}
  }
  | group_by_terms COMMA group_by_term {
    v := append($1.Terms, $3)
    $$ = &ast.GroupByTerms{Terms: v}
  }

group_by_term:
  expr {
    $$ = &ast.GroupByTerm{Expr: $1}
  }

having_section:
  {
    $$ = nil
  }
  | HAVING where_condition {
    $$ = &ast.HavingSection{Condition: $2}
  }

order_by_section:
  {
    $$ = nil
  }
  | ORDER BY order_by_terms {
    $$ = &ast.OrderBySection{Terms: $3}
  }

order_by_terms:
  order_by_term {
    $$ = &ast.OrderByTerms{Terms: []*ast.OrderByTerm{$1}}
  }
  | order_by_terms COMMA order_by_term {
    v := append($1.Terms, $3)
    $$ = &ast.OrderByTerms{Terms: v}
  }

order_by_term:
  expr order_by_term_option {
    opt := &ast.OrderByTermOption{
      IsDesc: $2,
    }
    $$ = &ast.OrderByTerm{
      Expr: $1,
      Option: opt,
    }
  }

order_by_term_option:
  {
    $$ = false
  }
  | ASC {
    $$ = false
  }
  | DESC {
    $$ = true
  }

limit_section:
  {
    $$ = nil
  }
  | LIMIT INT limit_offset {
    l := yylex.(Lexer)
    v := l.ParseInt($2.Value())
    $$ = &ast.LimitSection{
      Limit: &ast.IntLit{Value: v},
      Offset: $3,
    }
  }

limit_offset:
  {
    $$ = nil
  }
  | OFFSET INT {
    l := yylex.(Lexer)
    v := l.ParseInt($2.Value())
    $$ = &ast.IntLit{Value: v}
  }

exprs:
  expr {
    $$ = &ast.Exprs{Exprs: []ast.Expr{$1}}
  }
  | exprs COMMA expr {
    v := append($1.Exprs, $3)
    $$ = &ast.Exprs{Exprs: v}
  }

expr:
  expr OR expr {
    $$ = &ast.OrExpr{Left: $1, Right: $3}
  }
  | expr AND expr {
    $$ = &ast.AndExpr{Left: $1, Right: $3}
  }
  | expr XOR expr {
    $$ = &ast.XorExpr{Left: $1, Right: $3}
  }
  | NOT expr {
    $$ = &ast.NotExpr{Expr: $2}
  }
  | bool_primary {
    $$ = $1
  }

not_option:
  {
    $$ = false
  }
  | NOT {
    $$ = true
  }

bool_primary:
  bool_primary comparison_operator predicate {
    l := yylex.(Lexer)
    op := l.AsComparisonType($2.Type())
    $$ = &ast.BoolPrimaryComparison{
      Op: op,
      Left: $1,
      Right: $3,
    }
  }
  | predicate {
    $$ = &ast.BoolPrimaryPredicate{Pred: $1}
  }

comparison_operator:
  EQ | NE | GT | GQ | LT | LQ

predicate:
  bit_expr not_option IN LPAR exprs RPAR {
    $$ = &ast.PredicateIn{
      IsNot: $2,
      Target: $1,
      List: $5,
    }
  }
  | bit_expr not_option BETWEEN bit_expr AND predicate {
    $$ = &ast.PredicateBetween{
      IsNot: $2,
      Target: $1,
      Left: $4,
      Right: $6,
    }
  }
  | bit_expr not_option LIKE simple_expr {
    $$ = &ast.PredicateLike{
      IsNot: $2,
      Target: $1,
      Pattern: $4,
    }
  }
  | bit_expr {
    $$ = &ast.PredicateBitExpr{Expr: $1}
  }

bit_expr:
  bit_expr bit_operator bit_expr {
    l := yylex.(Lexer)
    op := l.AsBitOperatorType($2.Type())
    $$ = &ast.BitExprBitOp{Op: op, Left: $1, Right: $3}
  }
  | bit_expr arithmetic_operator bit_expr {
    l := yylex.(Lexer)
    op := l.AsArithmeticOperatorType($2.Type())
    $$ = &ast.BitExprArtOp{Op: op, Left: $1, Right: $3}
  }
  | simple_expr {
    $$ = &ast.BitExprSimpleExpr{Expr: $1}
  }

arithmetic_operator:
  PLUS | MINUS | AST | SLASH

bit_operator:
  AMP | PIPE | HAT

simple_expr:
  IDENT {
    $$ = &ast.Ident{Value: $1.Value()}
  }
  | function_call {
    $$ = $1
  }
  | prefix_operator simple_expr {
    l := yylex.(Lexer)
    op := l.AsPrefixOperatorType($1.Type())
    $$ = &ast.SimpleExprPrefixOp{Op: op, Expr: $2}
  }
  | literal {
    $$ = &ast.SimpleExprLit{Lit: $1}
  }
  | LPAR expr RPAR {
    $$ = &ast.SimpleExprExpr{Expr: $2}
  }

prefix_operator:
  PLUS | MINUS | TILDE | NOT

literal:
  INT {
    l := yylex.(Lexer)
    v := l.ParseInt($1.Value())
    $$ = &ast.IntLit{Value: v}
  }
  | FLOAT {
    l := yylex.(Lexer)
    v := l.ParseFloat($1.Value())
    $$ = &ast.FloatLit{Value: v}
  }
  | STRING {
    v := $1.Value()
    $$ = &ast.StringLit{Value: v}
  }

function_call:
  IDENT LPAR arg_list RPAR {
    name := &ast.Ident{Value: $1.Value()}
    $$ = &ast.FunctionCall{
      FunctionName: name,
      Arguments: $3,
    }
  }

arg_list:
  {
    $$ = &ast.Exprs{}
  }
  | exprs {
    $$ = $1
  }

%%
