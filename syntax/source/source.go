package source

import (
	"fmt"
	"lambda/domain"

	"golang.org/x/exp/utf8string"
)

type SourceCode struct {
	filename string
	text     utf8string.String
	tokens   []domain.Token
}

func NewSourceCode(filename string, text utf8string.String, tokens []domain.Token) SourceCode {
	return SourceCode{
		filename: filename,
		text:     text,
		tokens:   tokens,
	}
}

func (s SourceCode) Location(id domain.TokenId) (line, col int) {
	t := s.Token(id)
	line = t.Line
	col = t.Col
	return
}

func (s SourceCode) Lexeme(id domain.TokenId) string {
	t := s.Token(id)
	return s.text.Slice(int(t.Start), int(t.End))
}

func (s SourceCode) Filename() string {
	return s.filename
}

func (s SourceCode) Token(id domain.TokenId) domain.Token {
	return s.tokens[id]
}

func (s SourceCode) TokenCount() int {
	return len(s.tokens)
}

func (s SourceCode) TraceToken(tag domain.TokenId, lexeme string, line int, col int) string {
	str := fmt.Sprintf("\ttag = %d\n", tag)
	if lexeme != "" {
		str += fmt.Sprintf("\tlexeme = %#v\n", lexeme)
	}
	if line != -1 && col != -1 {
		str += "\tloc = " + fmt.Sprintf("%d", line) + ":" + fmt.Sprintf("%d", col) + "\n"
	}
	return str
}
