## General TODO list

- [x] Make lexer
- [x] Make parser (strict)
- [ ] Make parser (nonstrict) with following properties:
    * Functions can be non-unary
    * Applications have non-mandatory parens, that inferred by associativity rules:
        - Applications are left-associative
        - Abstractions are right-associative
- [ ] Make incremental contextes for each lambda term during postorder
traversal
    * They should include free/bound variables (for now)
- [ ] Add lexical bindings that simplify syntax
    * Should have the following form:
    ```
        let <id> = <expr> in <expr>
    ```
    * Expand on that form further, allowing for
    ```
        let <id1> = <expr1>[;\n] <id2> = <expr2>... in <expr>
    ```
    * Note that any such binding can be rewritten as application of abstraction to bound value
    ```
        let a = 5 in ((\f.f a) (\n.n)) => (\a.((\f.f a) (\n.n)) 5)
    ```
- [ ] Make std lib including all standard abstractions for sane programming from [here](https://www.lektorium.tv/sites/lektorium.tv/files/additional_files/20110227_systems_of_typed_lambda_calculi_moskvin_lecture02.pdf)
    * Booleans
    * Numbers
    * Pairs
    * Operations on numbers
    * Recursion
    * Lists
