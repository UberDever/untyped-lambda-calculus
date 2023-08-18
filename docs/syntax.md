# Untyped lambda calculus grammar

The following is the grammar for untyped lambda calculus in strict form 
(with only unary functions and mandatory parens around applications)

### Lexer

identifier ::= `[a-zA-Z+-*/=<>?!_][a-zA-Z0-9+-*/=<>?!_]*`

lambda ::= '\' | 'λ'

### Parser

term ::= 
    identifier
    | '(' application ')'
    | '('? abstraction ')'?
    | 'let' identifier '=' term 'in' term

application ::= term term

abstraction ::= lambda identifier '.' term
