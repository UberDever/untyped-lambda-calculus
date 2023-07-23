package syntax

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

	source_code := NewSourceCode("test", *text)
	tokenizer := NewTokenizer(&logger)
	tokenizer.Tokenize(&source_code)

	// strip eof
	tokens := source_code.tokens[:len(source_code.tokens)-1]

	for i, t := range tokens {
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
	source_code := NewSourceCode("test", *source)
	logger := domain.NewLogger()

	tokenizer := NewTokenizer(&logger)
	tokenizer.Tokenize(&source_code)
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
			abstraction := fmt.Sprintf(`(lambda %s %s)`, arg, body)
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

// TODO: This is strict form, but it also would be good to support
// convenient form (multiple arguments + inferred parens)
func TestAstSimple(test *testing.T) {
	text := `
        ((\x.\y.\z.(x (y z))) ((\i.i) something))
    `
	expected := `
        ((lambda x (lambda y (lambda z (x (y z))))) ((lambda i i) something))
    `
	if e := testAstEquality(text, expected); e != nil {
		test.Error(e)
	}
}
