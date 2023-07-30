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

	eval_stack := util.NewStack[any]()
	eval := func() {
		n := eval_stack.ForcePop().(int)
		tag := domain.NodeId(n)
		switch tag {
		case domain.NodeIdentifier:
			// nothing
		case domain.NodeApplication:
			lhs := eval_stack.ForcePop()
			rhs := eval_stack.ForcePop()
			application := fmt.Sprintf(`(%s %s)`, lhs, rhs)
			eval_stack.Push(application)
		case domain.NodeAbstraction:
			arg := eval_stack.ForcePop()
			body := eval_stack.ForcePop()
			abstraction := fmt.Sprintf(`(λ %s %s)`, arg, body)
			eval_stack.Push(abstraction)
		default:
			panic("unreachable")
		}
	}
	onEnter := func(s ast.Sexpr) {
		if s.IsAtom() {
			eval_stack.Push(s.Data())
		} else {
			eval()
		}
	}
	ast.TraversePostorder(tree, onEnter)
	eval()

	got := eval_stack.ForcePop().(string)
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
	expected := `x`
	if e := testAstEquality(text, expected); e != nil {
		test.Error(e)
	}
}

func TestAstAbstraction(test *testing.T) {
	text := `\x.x`
	expected := `(λ x x)`
	if e := testAstEquality(text, expected); e != nil {
		test.Error(e)
	}
}

func TestAstApplication(test *testing.T) {
	text := `((f g) h)`
	expected := `((f g) h)`
	if e := testAstEquality(text, expected); e != nil {
		test.Error(e)
	}
}

func TestAstSimple(test *testing.T) {
	text := `
        ((\x.\y.\z.(x (y z))) ((\i.i) something))
    `
	expected := `
        ((λ x (λ y (λ z (x (y z))))) ((λ i i) something))
    `
	if e := testAstEquality(text, expected); e != nil {
		test.Error(e)
	}
}

func TestAstUtf8(test *testing.T) {
	text := `
    ((\альфа.(альфа бета)) гамма)
    `
	expected := `
        ((λ альфа (альфа бета)) гамма)
    `
	if e := testAstEquality(text, expected); e != nil {
		test.Error(e)
	}
}
