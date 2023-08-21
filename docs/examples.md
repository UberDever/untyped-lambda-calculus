
### Let expression rewrite
```
let a = 4 in
let b = 5 in
(((λx.λy.((+ x) y)) a) b)
=let-rewrite>
(((λa.λb.(((λx.λy.((+ x) y)) a) b)) 4) 5)
=indices>
(((λ λ (((λ λ ((+ 1) 0)) 1) 0)) 0) 1)
=eval>
((+ 4) 5)
```

### Let in lambda body
```
λx. let y = 4 in (x y)
=let-rewrite>
((λy.λx.(x y)) 4)
```

### Let in application
```
(f (let g = 8 in g))
=let-rewrite>
(f ((λg.g) 8))
```

### Let arithemtic
```
    let a = -7 in
    let b = 69 in
    let c = 42 in
    ((* c) ((+ a) b))
=let-rewrite>
    ((λa.((λb.((λc.((* c)((+ a) b))) 42)) 69)) -7)
=indices>
    ((λ ((λ ((λ ((3 0)((4 2) 1))) 4)) 4)) 4)
```
