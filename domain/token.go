package domain

import (
	"math"
)

type Token int

const TokenInvalid = math.MinInt
const TokenEof = -1

const (
	TokenIdentifier = iota
	TokenDot
	TokenLambda
	TokenLeftParen
	TokenRightParen
)

const (
	TokenDotRune        rune = '.'
	TokenLambdaRune     rune = '\\'
	TokenLeftParenRune  rune = '('
	TokenRightParenRune rune = ')'
)
