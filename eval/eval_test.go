package eval

import (
	"errors"
	"fmt"
	"lambda/ast"
	"lambda/domain"
	"lambda/syntax"
	"lambda/util"
	"strings"
	"testing"

	"golang.org/x/exp/utf8string"
)

func parse_tree(text string) (tree ast.Sexpr, err error) {
	report_errors := func(logger *domain.Logger) error {
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
	source_code := syntax.NewSourceCode("test", *source)
	logger := domain.NewLogger()

	tokenizer := syntax.NewTokenizer(&logger)
	tokenizer.Tokenize(&source_code)
	if !logger.IsEmpty() {
		err = report_errors(&logger)
		return
	}

	parser := syntax.NewParser(&logger)
	tree = parser.Parse(&source_code)
	if !logger.IsEmpty() {
		err = report_errors(&logger)
		return
	}
	return
}

func testEvalEquality(text, expected string) error {
	tree, err := parse_tree(text)
	if err != nil {
		return err
	}

	ctx := NewEvalContext()
	evaluated := ctx.Eval(tree)
	got := ToString(evaluated, false)
	if ast.Minified(got) != ast.Minified(expected) {
		lhs := ast.Pretty(got)
		rhs := ast.Pretty(expected)
		trace := util.ConcatVertically(lhs, rhs)
		return fmt.Errorf("AST are not equal\n%s", trace)
	}
	return nil
}

func TestEvalPrimitive(test *testing.T) {
	text := `x`
	expected := `x`
	if e := testEvalEquality(text, expected); e != nil {
		test.Error(e)
	}
}

func TestEvalAbstraction(test *testing.T) {
	text := `\x.x`
	expected := `(\ x x)`
	if e := testEvalEquality(text, expected); e != nil {
		test.Error(e)
	}
}

func TestEvalBoundVariables(test *testing.T) {
	text := `\x.\y.\z.((((f g) (h x)) y) z)`
	expected_bound := []string{"x", "y", "z"}
	expected_free := []string{"f", "g", "h"}

	tree, err := parse_tree(text)
	if err != nil {
		test.Error(err)
	}
	ctx := NewEvalContext()
	_ = ctx.Eval(tree)
	for _, name := range expected_bound {
		if !ctx.GetBound().Has(name) {
			test.Errorf("Name %s should be bound in %s", name, ToString(tree, true))
		}
	}
	for _, name := range expected_free {
		if !ctx.GetFree().Has(name) {
			test.Errorf("Name %s should be free in %s", name, ToString(tree, true))
		}
	}
}

func TestEvalUnreducable(test *testing.T) {
	text := `((f g) h)`
	expected := `((f g) h)`
	if e := testEvalEquality(text, expected); e != nil {
		test.Error(e)
	}
}

// TODO: Develop this further
// func TestEvalApplication(test *testing.T) {
// 	text := `((\x.\y.\z.(y z x)) (x y z))`
// 	expected := `(\y'\z'.(y' z' (x y z)))`
// 	if e := testEvalEquality(text, expected); e != nil {
// 		test.Error(e)
// 	}
// }
