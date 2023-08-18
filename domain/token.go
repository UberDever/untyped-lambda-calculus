package domain

import (
	"math"
)

type TokenId int

const TokenInvalid TokenId = math.MinInt
const TokenEof TokenId = -1

const (
	TokenIdentifier TokenId = iota
	TokenDot
	TokenLambda
	TokenLeftParen
	TokenRightParen
)

const (
	TokenDotRune             rune = '.'
	TokenLambdaBackslashRune rune = '\\'
	TokenLambdaRune          rune = 'λ'
	TokenLeftParenRune       rune = '('
	TokenRightParenRune      rune = ')'
)

type Token struct {
	Tag                   TokenId
	Start, End, Line, Col int
}

func NewToken(tag TokenId, start, end, line, col int) Token {
	return Token{Tag: tag, Start: start, End: end, Line: line, Col: col}
}

func NewTokenEof() Token {
	return NewToken(TokenEof, -1, -1, -1, -1)
}
