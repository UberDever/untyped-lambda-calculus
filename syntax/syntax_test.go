package syntax

import (
	"lambda/domain"
	"testing"

	"golang.org/x/exp/utf8string"
)

func TestTokenizer(test *testing.T) {
	text := utf8string.NewString(`
        \x.\y.(f g) h
    `)
	// expected := [...]struct {
	// 	domain.Token
	// 	string
	// }{
	// 	{domain.TokenDot, `.`},
	// }

	logger := domain.NewLogger()

	source_code := NewSourceCode("test", *text)
	tokenizer := NewTokenizer(&logger)
	tokenizer.Tokenize(&source_code)

	for _, t := range source_code.tokens {
		test.Log(t)
		if t.Tag == domain.TokenEof {
			break
		}
		asStr := text.Slice(t.Start, t.End)
		test.Log(asStr)
	}
}
