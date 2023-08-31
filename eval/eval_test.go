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
	namedTree := parser.Parse(&source_code)

	result := debruijn.ToDeBruijn(&source_code, &namedTree)
	if !logger.IsEmpty() {
		return report_errors(&logger)
	}
	de_bruijn_tree := result.Tree

	got := ast.Print(&source_code, &de_bruijn_tree)
	if sexpr.Minified(got) != sexpr.Minified(expected) {
		lhs := sexpr.Pretty(got)
		rhs := sexpr.Pretty(expected)
		trace := util.ConcatVertically(lhs, rhs)
		return fmt.Errorf("AST are not equal\n%s", trace)
	}
	return nil
}

//
// func testEvalEquality(text, expected string) error {
// 	tree, err := parse_tree(text)
// 	if err != nil {
// 		return err
// 	}
//
// 	ctx := NewEvalContext()
// 	evaluated := ctx.Eval(tree)
// 	got := ToString(evaluated, false)
// 	if ast.Minified(got) != ast.Minified(expected) {
// 		lhs := ast.Pretty(got)
// 		rhs := ast.Pretty(expected)
// 		trace := util.ConcatVertically(lhs, rhs)
// 		return fmt.Errorf("AST are not equal\n%s", trace)
// 	}
// 	return nil
// }
//
// func TestEvalPrimitive(test *testing.T) {
// 	text := `x`
// 	expected := `x`
// 	if e := testEvalEquality(text, expected); e != nil {
// 		test.Error(e)
// 	}
// }
//
// func TestEvalAbstraction(test *testing.T) {
// 	text := `λx.x`
// 	expected := `(\ x x)`
// 	if e := testEvalEquality(text, expected); e != nil {
// 		test.Error(e)
// 	}
// }
//
// func TestEvalBoundVariables(test *testing.T) {
// 	text := `λx.λy.λz.((((f g) (h x)) y) z)`
// 	expected_bound := []string{"x", "y", "z"}
// 	expected_free := []string{"f", "g", "h"}
//
// 	tree, err := parse_tree(text)
// 	if err != nil {
// 		test.Error(err)
// 	}
// 	ctx := NewEvalContext()
// 	_ = ctx.Eval(tree)
// 	for _, name := range expected_bound {
// 		if !ctx.GetBound().Has(name) {
// 			test.Errorf("Name %s should be bound in %s", name, ToString(tree, true))
// 		}
// 	}
// 	for _, name := range expected_free {
// 		if !ctx.GetFree().Has(name) {
// 			test.Errorf("Name %s should be free in %s", name, ToString(tree, true))
// 		}
// 	}
// }
//
// func TestEvalUnreducable(test *testing.T) {
// 	text := `((f g) h)`
// 	expected := `((f g) h)`
// 	if e := testEvalEquality(text, expected); e != nil {
// 		test.Error(e)
// 	}
// }

// func TestEvalWHNF(test *testing.T) {
// 	text := `((λx.λy.(x y)) y)`
// 	expected := ``
// }

// TODO: Develop this further
// func TestEvalApplication(test *testing.T) {
// 	text := `((λx.λy.λz.(y z x)) (x y z))`
// 	expected := `(\y'\z'.(y' z' (x y z)))`
// 	if e := testEvalEquality(text, expected); e != nil {
// 		test.Error(e)
// 	}
// }
