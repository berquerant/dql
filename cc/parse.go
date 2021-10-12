package cc

func Parse(lexer Lexer) int { return yyParse(lexer) }
