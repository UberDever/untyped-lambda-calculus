package parser

import (
	"errors"
	"fmt"
	"lambda/ast"
	"lambda/domain"
	"lambda/middle"
	"lambda/util"
	"strings"
	"testing"

	"golang.org/x/exp/utf8string"
)

func TestTokenizer(test *testing.T) {
	text := utf8string.NewString(`
        λx.λy.(foo bar) baz
    `)
	expected := [...]struct {
		domain.TokenId
		string
	}{
		{domain.TokenLambda, `λ`},
		{domain.TokenIdentifier, "x"},
		{domain.TokenDot, `.`},
		{domain.TokenLambda, `λ`},
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
	namedTree := parser.Parse(&source_code)
	deBruijn := middle.ToDeBruijn(namedTree)
	if !logger.IsEmpty() {
		return report_errors(&logger)
	}

	got := deBruijn.Print()
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
	text := `λx.x`
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
    ((λальфа.(альфа бета)) гамма)
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
        (((λx.x) free_var) (λfree_var.free_var)) 
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
        ((λx.λy.λz.(x (y z))) ((λi.i) something))
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
        ((λx.λx.x) (λy.y))
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
        (λz.((x y) z)) 
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
        ((λu.λv.(u x)) y)
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
        λx.λy.λs.λz.((x s) ((y s) z))
    `
	expected := `
        (λ (λ (λ (λ ((3 1) ((2 1) 0))))))
    `
	if e := testAstEquality(text, expected); e != nil {
		test.Error(e)
	}
}

func TestAstSimpleLet(test *testing.T) {
	text := `
        let u = y in λv.(u x)
    `
	// ((λu.λv.(u x)) y)
	expected := `
        ((λ (λ (1 2))) 1)
    `
	if e := testAstEquality(text, expected); e != nil {
		test.Error(e)
	}
}

func TestAstMultipleLet(test *testing.T) {
	text := `
        let a = -7 in
        let b = 69 in
        let c = 42 in
        ((* c) ((+ a) b))
    `
	expected := `
        ((λ ((λ ((λ ((3 0)((4 2) 1))) 4)) 4)) 4)
    `
	if e := testAstEquality(text, expected); e != nil {
		test.Error(e)
	}
}

func TestAstUnclosedParen(test *testing.T) {
	text := ` (λx.x `
	expected := `λ 0`
	if e := testAstEquality(text, expected); e == nil {
		test.Error("Expected error with unclosed paren")
	}
}

func TestAstApplicationWithoutParens(test *testing.T) {
	text := `x y`
	expected := `(0 1)`
	if e := testAstEquality(text, expected); e == nil {
		test.Error("Expected error with unparentesized application")
	}
}
