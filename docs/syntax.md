# Untyped lambda calculus grammar

The following is the grammar for untyped lambda calculus in strict form:
1. Only unary functions
2. Parens are mandatory

### Lexer

identifier ::= [a-zA-Z+\-*/=<>?!_.][a-zA-Z0-9+\-*/=<>?!_.]*
dot ::= '.'
lambda ::= '\'
left_paren ::= '('
right_paren ::= ')'

### Parser

term ::= 
    identifier
    | left_paren term term right_paren
    | left_paren lambda identifier dot term right_paren
