## General TODO list

- [x] Develop grammar
- [x] Make lexer
- [x] Make parser (strict)
- [ ] Make parser (nonstrict) with following properties:
    * Functions can be non-unary (arguments separated by some symbol like ` ` or `,`)
    * Applications have non-mandatory parens, that inferred by associativity rules:
        - Applications are left-associative
        - Abstractions are right-associative
- [x] Use [De bruijn indicies](https://www.researchgate.net/publication/2368794_Reviewing_the_Classical_and_the_de_Bruijn_Notation_for_-calculus_and_Pure_Type_Systems) 
for evaluation to avoid alpla-conversion. Hence, make AST (or equivalent structure) to represent the use of such indicies. 
Also use [slides](https://www.cs.vu.nl/~femke/courses/ep/slides/4x4.pdf)
- [x] For above use [sophisticated algorithm from here](https://www.researchgate.net/publication/2368794_Reviewing_the_Classical_and_the_de_Bruijn_Notation_for_-calculus_and_Pure_Type_Systems) for getting indicies for nodes
      [**Note: used [lecture](https://www.cs.cornell.edu/courses/cs4110/2018fa/lectures/lecture15.pdf)**]
- [ ] Use normal order evaluation with WHNF
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
        let a = 5 in ((λf.f a) (λn.n)) => (λa.((λf.f a) (λn.n)) 5)
    ```
- [ ] Make std lib including all standard abstractions for sane programming from [here](https://www.lektorium.tv/sites/lektorium.tv/files/additional_files/20110227_systems_of_typed_lambda_calculi_moskvin_lecture02.pdf)
    * Booleans
    * Numbers
    * Pairs
    * Operations on numbers
    * Recursion
    * Lists
