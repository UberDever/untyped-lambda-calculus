## Features of lambda calculus as programming language

1. Identifiers as atoms
2. Abstraction as a way of directing computation (control flow) and specifying means of abstraction (wow)
3. Application as only way of describing evaluation
4. Bindings in form of `let x = expr1 in expr2`
    - They are nested
    - Actually syntactic sugar for forming redex
5. This calculus is untyped
6. Primary evaluation strategy - to WHNF (Call by name / normal order)
7. AST has 2 forms - normal and de-bruijn. Latter is used as interpretation target.
