# Untyped lambda calculus grammar

The following is the grammar for untyped lambda calculus in strict form (with only unary functions and mandatory parens)

### Lexer

identifier ::= `[a-zA-Z+-*/=<>?!_][a-zA-Z0-9+-*/=<>?!_]*`

dot ::= '.'

lambda ::= '\'

left_paren ::= '('

right_paren ::= ')'

### Parser

term ::= 
    identifier
    | left_paren (application | abstraction) right_paren

application ::= term term

abstraction ::= lambda identifier dot term
