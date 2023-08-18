
### Let expression rewrite
```
let a = 4 in
let b = 5 in
(((λx.λy.((+ x) y)) a) b)
=let-rewrite>
(((λa.λb.(((λx.λy.((+ x) y)) a) b)) 4) 5)
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
