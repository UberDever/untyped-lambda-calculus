package eval

import (
	"errors"
	"fmt"
	"lambda/ast/ast"
	"lambda/ast/sexpr"
	debruijn "lambda/middle/de-bruijn"
	"lambda/syntax/parser"
	"lambda/util"
	"strings"
	"testing"

	"golang.org/x/exp/utf8string"
)

func testEvalEquality(text, expected string) error {
	report_errors := func(logger *util.Logger) error {
		builder := strings.Builder{}
		for {
			m, ok := logger.Next()
			if !ok {
				break
			}
			builder.WriteString(m.String())
			builder.WriteByte('\n')
		}
		return errors.New(builder.String())
	}

	source := utf8string.NewString(text)
	logger := util.NewLogger()

	tokenizer := parser.NewTokenizer(&logger)
	source_code := tokenizer.Tokenize("test", *source)
	if !logger.IsEmpty() {
		return report_errors(&logger)
	}

	parser := parser.NewParser(&logger)
	namedTree := parser.Parse(source_code)
	if !logger.IsEmpty() {
		return report_errors(&logger)
	}

	result := debruijn.ToDeBruijn(source_code, namedTree)
	de_bruijn_tree := result.Tree

	eval_tree := Eval(&logger, source_code, de_bruijn_tree)
	// for !logger.IsEmpty() {
	// 	m, _ := logger.Next()
	// 	// fmt.Println(m)
	// }

	got := ast.Print(source_code, eval_tree)
	if sexpr.Minified(got) != sexpr.Minified(expected) {
		lhs := sexpr.Pretty(got)
		rhs := sexpr.Pretty(expected)
		trace := util.ConcatVertically(lhs, rhs)
		return fmt.Errorf("AST are not equal\n%s", trace)
	}
	return nil
}

func TestEvalNonRedex(test *testing.T) {
	{
		text := `x`
		expected := `0`
		if e := testEvalEquality(text, expected); e != nil {
			test.Error(e)
		}
	}
	{
		text := `λx.x`
		expected := `(λ 0)`
		if e := testEvalEquality(text, expected); e != nil {
			test.Error(e)
		}
	}
	{
		text := `(f g)`
		expected := `(0 1)`
		if e := testEvalEquality(text, expected); e != nil {
			test.Error(e)
		}
	}
	{
		text := `(f (g h))`
		expected := `(0 (1 2))`
		if e := testEvalEquality(text, expected); e != nil {
			test.Error(e)
		}
	}
}

func TestEvalSimpleRedex(test *testing.T) {
	text := `((λx.x) y)`
	expected := `0`
	if e := testEvalEquality(text, expected); e != nil {
		test.Error(e)
	}
}

func TestEvalNormalForm(test *testing.T) {
	text := `
        λx1.λx2.λx3.(((y N1) N2) N3)
    `
	expected := `(λ (λ (λ (((3 4) 5) 6))))`
	if e := testEvalEquality(text, expected); e != nil {
		test.Error(e)
	}
}

func TestEvalRedex1(test *testing.T) {
	text := `((λu.λv.(u x)) y)`
	expected := `(λ (2 1))`
	if e := testEvalEquality(text, expected); e != nil {
		test.Error(e)
	}
}

func TestEvalRedex2(test *testing.T) {
	text := `((((λx.x) N1) N2) N3) `
	expected := `((0 1) 2)`
	if e := testEvalEquality(text, expected); e != nil {
		test.Error(e)
	}
}

func TestEvalRedex3(test *testing.T) {
	text := `((λx.x) ((λy.y) ((λz.z) N))) `
	expected := `0`
	if e := testEvalEquality(text, expected); e != nil {
		test.Error(e)
	}
}

func TestEvalSKI(test *testing.T) {
	// Since evaluation goes to WHNF, this SKK example should be applied to something
	// to test it and because SKK == I then (I something) ->β something
	text := `
    let K = λx.λy.x in
    let S = λx.λy.λz.((x z) (y z)) in
    let I = λx.x in
    ((S K) K)
    `
	expected := `2`
	if e := testEvalEquality(text, expected); e != nil {
		test.Error(e)
	}
}

// func TestFactorial(test *testing.T) {
// 	text := `
//     let True = λt.λf.t in
//     let False = λt.λf.f in
//     let If = λb.λx.λy.((b x) y) in
//     let And = λp.λq.((p q) p) in
//     let Or = λp.λq.((p p) q) in
//     let Not = λp.((p False) True) in
//
//     let Pair = λx.λy.λf.((f x) y) in
//     let Fst = λp.(p True) in
//     let Snd = λp.(p False) in
//
//     let 0 = False in
//     let Succ = λn.λs.λz.(s ((n s) z)) in
//     let 1 = (Succ 0) in
//     let 2 = (Succ 1) in
//     let 3 = (Succ 2) in
//     let 4 = (Succ 3) in
//
//     let Plus = λm.λn.λs.λz.((m s) ((n s) z)) in
//     let Mult = λm.λn.λs.(m (n s)) in
//     let Pow = λb.λe.(e b) in
//     let IsZero = λn.((n (λx.False)) True) in
//     let Pred = λn.λf.λx.(((n (λg.λh.(h (g f)))) (λu.x)) (λu.u)) in
//
// 	let Y = λf.((λx.(f (x x))) (λx.(f (x x)))) in
// 	let Fact = λf.λn.(((If (IsZero n)) 1) ((Mult n) (f (Pred n)))) in
// 	let FactRec = (Y Fact) in
//         (FactRec 5)
//     `
// 	expected := `2`
// 	if e := testEvalEquality(text, expected); e != nil {
// 		test.Error(e)
// 	}
// }
