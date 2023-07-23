package domain

import (
	"fmt"
	"math"
	"regexp"
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
	TokenDotString        rune = '.'
	TokenLambdaString     rune = '\\'
	TokenLeftParenString  rune = '('
	TokenRightParenString rune = ')'
)

var (
	TokenIdentifierRegex  = regexp.MustCompile(`[a-zA-Z+\-*/=<>?!_.][a-zA-Z0-9+\-*/=<>?!_.]*`)
	TokenDotRegexp        = regexp.MustCompile(fmt.Sprintf(`\%c`, TokenDotString))
	TokenLambdaRegexp     = regexp.MustCompile(fmt.Sprintf(`\%c`, TokenLambdaString))
	TokenLeftParenRegexp  = regexp.MustCompile(fmt.Sprintf(`\%c`, TokenLeftParenString))
	TokenRightParenRegexp = regexp.MustCompile(fmt.Sprintf(`\%c`, TokenRightParenString))
)
