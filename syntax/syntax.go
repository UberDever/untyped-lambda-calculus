package syntax

import (
	"lambda/domain"
	"unicode"

	"golang.org/x/exp/utf8string"
)

type token struct {
	Tag                   domain.Token
	Start, End, Line, Col int
}

type source_code struct {
	filename string
	text     utf8string.String
	tokens   []token
}

func NewSourceCode(filename string, text utf8string.String) source_code {
	return source_code{
		filename: filename,
		text:     text,
		tokens:   nil,
	}
}

type tokenizer struct {
	logger *domain.Logger
}

func NewTokenizer(logger *domain.Logger) tokenizer {
	return tokenizer{logger: logger}
}

func (tok *tokenizer) Tokenize(src *source_code) {
	tokens := make([]token, 0, 16)
	pos := 0
	line, col := 1, 0

	add_token := func(t domain.Token, length int) {
		pos += length
		col += length
		tokens = append(tokens, token{t, pos, pos + length, line, col})
	}

	skip_spaces := func() {
		for pos < src.text.RuneCount() {
			c := src.text.At(pos)
			if !unicode.IsSpace(c) {
				break
			}
			if c == '\n' {
				line++
				col = 0
			}
			pos++
		}
	}

	for {
		skip_spaces()
		if pos >= src.text.RuneCount() {
			break
		}

		switch src.text.At(pos) {
		case domain.TokenDotString:
			add_token(domain.TokenDot, 1)
		case domain.TokenLambdaString:
			add_token(domain.TokenLambda, 1)
		case domain.TokenLeftParenString:
			add_token(domain.TokenLeftParen, 1)
		case domain.TokenRightParenString:
			add_token(domain.TokenRightParen, 1)
		default:
			index := domain.TokenIdentifierRegex.FindIndex([]byte(src.text.Slice(pos, src.text.RuneCount())))
			if index == nil {
				panic("Something went wrong")
			}
			start, end := index[0], index[1]
			length := end - start + 1
			add_token(domain.TokenIdentifier, length)
		}
	}

	tokens = append(tokens, token{domain.TokenEof, -1, -1, -1, -1})
	src.tokens = tokens
}
