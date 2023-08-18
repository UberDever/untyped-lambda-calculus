package parser

import (
	"errors"
	"fmt"
	"lambda/ast"
	"lambda/domain"
	"lambda/util"
	"strings"
	"testing"

	"golang.org/x/exp/utf8string"
)

func TestTokenizer(test *testing.T) {
	text := utf8string.NewString(`
        \x.\y.(foo bar) baz
    `)
	expected := [...]struct {
		domain.TokenId
		string
	}{
		{domain.TokenLambda, `\`},
		{domain.TokenIdentifier, "x"},
		{domain.TokenDot, `.`},
		{domain.TokenLambda, `\`},
		{domain.TokenIdentifier, "y"},
		{domain.TokenDot, `.`},
		{domain.TokenLeftParen, `(`},
		{domain.TokenIdentifier, "foo"},
		{domain.TokenIdentifier, "bar"},
		{domain.TokenRightParen, `)`},
		{domain.TokenIdentifier, "baz"},
	}

	logger := domain.NewLogger()

	tokenizer := NewTokenizer(&logger)
	source_code := tokenizer.Tokenize("test", *text)

	// strip eof in iteration (note TokenCount - 1)
	for i := 0; i < source_code.TokenCount()-1; i++ {
		t := source_code.Token(domain.TokenId(i))
		asStr := text.Slice(t.Start, t.End)
		if expected[i].string != asStr ||
			expected[i].TokenId != t.Tag {
			test.Fatalf("Expected [%d %s] got [%d %s]",
				expected[i].TokenId, expected[i].string,
				t.Tag, asStr,
			)
		}
	}
}

func testAstEquality(text, expected string) error {
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
	logger := domain.NewLogger()

	tokenizer := NewTokenizer(&logger)
	source_code := tokenizer.Tokenize("test", *source)
	if !logger.IsEmpty() {
		return report_errors(&logger)
	}

	parser := NewParser(&logger)
	tree := parser.Parse(&source_code)
	if !logger.IsEmpty() {
		return report_errors(&logger)
	}

	got := tree.Print()
	if ast.Minified(got) != ast.Minified(expected) {
		lhs := ast.Pretty(got)
		rhs := ast.Pretty(expected)
		trace := util.ConcatVertically(lhs, rhs)
		return fmt.Errorf("AST are not equal\n%s", trace)
	}
	return nil
}

func TestAstPrimitive(test *testing.T) {
	text := `x`
	expected := `0`
	if e := testAstEquality(text, expected); e != nil {
		test.Error(e)
	}
}

func TestAstAbstraction(test *testing.T) {
	text := `\x.x`
	expected := `(λ 0)`
	if e := testAstEquality(text, expected); e != nil {
		test.Error(e)
	}
}

func TestAstApplication(test *testing.T) {
	text := `((f g) h)`
	expected := `((0 1) 2)`
	if e := testAstEquality(text, expected); e != nil {
		test.Error(e)
	}
}

func TestAstUtf8(test *testing.T) {
	text := `
    ((\альфа.(альфа бета)) гамма)
    `
	expected := `
        ((λ (0 1)) 1)
    `
	if e := testAstEquality(text, expected); e != nil {
		test.Error(e)
	}
}

func TestAstFreeVarEncounteredAfter(test *testing.T) {
	text := `
        (((\x.x) free_var) (\free_var.free_var)) 
    `
	expected := `
        (((λ 0) 0) (λ 0))
    `
	if e := testAstEquality(text, expected); e != nil {
		test.Error(e)
	}
}

func TestAst1(test *testing.T) {
	text := `
        ((\x.\y.\z.(x (y z))) ((\i.i) something))
    `
	expected := `
        ((λ (λ (λ (2 (1 0))))) ((λ 0) 0))
    `
	if e := testAstEquality(text, expected); e != nil {
		test.Error(e)
	}
}

func TestAst2(test *testing.T) {
	text := `
        ((\x.\x.x) (\y.y))
    `
	expected := `
        ((λ (λ 0)) (λ 0))
    `
	if e := testAstEquality(text, expected); e != nil {
		test.Error(e)
	}
}

func TestAst3(test *testing.T) {
	text := `
        (\z.((x y) z)) 
    `
	expected := `
        (λ ((1 2) 0))
    `
	if e := testAstEquality(text, expected); e != nil {
		test.Error(e)
	}
}

func TestAst4(test *testing.T) {
	text := `
        ((\u.\v.(u x)) y)
    `
	expected := `
        ((λ (λ (1 2))) 1)
    `
	if e := testAstEquality(text, expected); e != nil {
		test.Error(e)
	}
}

func TestAst5(test *testing.T) {
	text := `
        \x.\y.\s.\z.((x s) ((y s) z))
    `
	expected := `
        (λ (λ (λ (λ ((3 1) ((2 1) 0))))))
    `
	if e := testAstEquality(text, expected); e != nil {
		test.Error(e)
	}
}
