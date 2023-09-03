package source

import (
	"fmt"

	"golang.org/x/exp/utf8string"
)

type SourceCode struct {
	filename string
	text     utf8string.String
	tokens   []Token
}

func NewSourceCode(filename string, text utf8string.String, tokens []Token) SourceCode {
	return SourceCode{
		filename: filename,
		text:     text,
		tokens:   tokens,
	}
}

func (s SourceCode) Location(id TokenId) (line, col int) {
	t := s.Token(id)
	line = t.Line
	col = t.Col
	return
}

func (s SourceCode) Lexeme(id TokenId) string {
	t := s.Token(id)
	return s.text.Slice(int(t.Start), int(t.End))
}

func (s SourceCode) Filename() string {
	return s.filename
}

func (s SourceCode) Token(id TokenId) Token {
	return s.tokens[id]
}

func (s SourceCode) TokenCount() int {
	return len(s.tokens)
}

func (s SourceCode) TraceToken(tag TokenId, lexeme string, line int, col int) string {
	str := fmt.Sprintf("\ttag = %d\n", tag)
	if lexeme != "" {
		str += fmt.Sprintf("\tlexeme = %#v\n", lexeme)
	}
	if line != -1 && col != -1 {
		str += "\tloc = " + fmt.Sprintf("%d", line) + ":" + fmt.Sprintf("%d", col) + "\n"
	}
	return str
}
