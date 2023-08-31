package parser

import (
	"lambda/syntax/source"
	"lambda/util"
	"testing"

	"golang.org/x/exp/utf8string"
)

func TestTokenizer(test *testing.T) {
	text := utf8string.NewString(`
        λx.λy.(foo bar) baz
    `)
	expected := [...]struct {
		source.TokenId
		string
	}{
		{source.TokenLambda, `λ`},
		{source.TokenIdentifier, "x"},
		{source.TokenDot, `.`},
		{source.TokenLambda, `λ`},
		{source.TokenIdentifier, "y"},
		{source.TokenDot, `.`},
		{source.TokenLeftParen, `(`},
		{source.TokenIdentifier, "foo"},
		{source.TokenIdentifier, "bar"},
		{source.TokenRightParen, `)`},
		{source.TokenIdentifier, "baz"},
	}

	logger := util.NewLogger()

	tokenizer := NewTokenizer(&logger)
	source_code := tokenizer.Tokenize("test", *text)

	// strip eof in iteration (note TokenCount - 1)
	for i := 0; i < source_code.TokenCount()-1; i++ {
		t := source_code.Token(source.TokenId(i))
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
