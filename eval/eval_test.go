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

	eval_tree := Eval(source_code, de_bruijn_tree)

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

// TODO: Develop this further
// func TestEvalApplication(test *testing.T) {
// 	text := `((λx.λy.λz.(y z x)) (x y z))`
// 	expected := `(\y'\z'.(y' z' (x y z)))`
// 	if e := testEvalEquality(text, expected); e != nil {
// 		test.Error(e)
// 	}
// }
