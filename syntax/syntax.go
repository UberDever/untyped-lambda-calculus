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
		start, end := pos, pos+length
		tokens = append(tokens, token{t, start, end, line, col})
		pos = end
		col = end
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

	identifier_rune := func(r rune) bool {
		return r != domain.TokenDotRune &&
			r != domain.TokenLambdaRune &&
			r != domain.TokenLeftParenRune &&
			r != domain.TokenRightParenRune &&
			!unicode.IsSpace(r)
	}
	identifier_length := func() int {
		start, end := pos, pos
		for end < src.text.RuneCount() {
			c := src.text.At(end)
			if !identifier_rune(c) {
				return end - start
			}
			end++
		}
		return end - start
	}

	for {
		skip_spaces()
		if pos >= src.text.RuneCount() {
			break
		}

		switch src.text.At(pos) {
		case domain.TokenDotRune:
			add_token(domain.TokenDot, 1)
		case domain.TokenLambdaRune:
			add_token(domain.TokenLambda, 1)
		case domain.TokenLeftParenRune:
			add_token(domain.TokenLeftParen, 1)
		case domain.TokenRightParenRune:
			add_token(domain.TokenRightParen, 1)
		default:
			length := identifier_length()
			if length > 0 {
				add_token(domain.TokenIdentifier, length)
			}
		}
	}

	tokens = append(tokens, token{domain.TokenEof, -1, -1, -1, -1})
	src.tokens = tokens
}
