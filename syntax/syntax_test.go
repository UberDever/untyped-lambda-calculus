package syntax

import (
	"lambda/domain"
	"testing"

	"golang.org/x/exp/utf8string"
)

func TestTokenizer(test *testing.T) {
	text := utf8string.NewString(`
        \x.\y.(foo bar) baz
    `)
	expected := [...]struct {
		domain.Token
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
			expected[i].Token != t.Tag {
			test.Fatalf("Expected [%d %s] got [%d %s]",
				expected[i].Token, expected[i].string,
				t.Tag, asStr,
			)
		}
	}
}
