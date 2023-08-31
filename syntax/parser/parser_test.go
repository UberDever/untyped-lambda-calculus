package parser

import (
	"lambda/domain"
	"lambda/util"
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

	logger := util.NewLogger()

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
