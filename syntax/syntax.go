package syntax

import (
	"fmt"
	"lambda/ast"
	"lambda/domain"
	"unicode"

	"golang.org/x/exp/utf8string"
)

type token struct {
	Tag                   domain.TokenId
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

func (s source_code) Location(id domain.TokenId) (line, col int) {
	t := s.Token(id)
	line = t.Line
	col = t.Col
	return
}

func (s source_code) Lexeme(id domain.TokenId) string {
	t := s.Token(id)
	return s.text.Slice(int(t.Start), int(t.End))
}

func (s source_code) Filename() string {
	return s.filename
}

func (s source_code) Token(id domain.TokenId) token {
	return s.tokens[id]
}

func (s source_code) TraceToken(tag domain.TokenId, lexeme string, line int, col int) string {
	str := fmt.Sprintf("\ttag = %d\n", tag)
	if lexeme != "" {
		str += fmt.Sprintf("\tlexeme = %#v\n", lexeme)
	}
	if line != -1 && col != -1 {
		str += "\tloc = " + fmt.Sprintf("%d", line) + ":" + fmt.Sprintf("%d", col) + "\n"
	}
	return str
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

	add_token := func(t domain.TokenId, length int) {
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

type parser struct {
	ast ast.Sexpr
	src *source_code

	current domain.TokenId
	atEof   bool

	logger *domain.Logger
}

func (p *parser) next() {
	for {
		p.current++
		c := p.src.Token(p.current)
		p.atEof = c.Tag == domain.TokenEof
		return
	}
}

func (p *parser) matchTag(tag domain.TokenId) bool {
	if p.atEof {
		return false
	}

	c := p.src.Token(p.current)
	return c.Tag == tag
}

func (p *parser) expect(tag domain.TokenId) (ok bool) {
	if p.atEof {
		return
	}

	ok = p.matchTag(tag)

	if !ok {
		c := p.src.Token(p.current)
		expected := p.src.TraceToken(tag, "", int(domain.TokenEof), int(domain.TokenEof))
		got := p.src.TraceToken(c.Tag, p.src.Lexeme(p.current), c.Line, c.Col)
		message := fmt.Sprintf("\nExpected\n %s but got\n %s", expected, got)
		p.logger.Add(domain.NewMessage(domain.Fatal, c.Line, c.Col, p.src.filename, message))
		p.atEof = true
		return
	}

	p.next()
	return
}

func NewParser(logger *domain.Logger) parser {
	return parser{logger: logger}
}

func (p *parser) Parse(src *source_code) ast.Sexpr {
	p.src = src
	return p.parse_term()
}

func (p *parser) parse_term() ast.Sexpr {
	var node ast.Sexpr
	if !p.matchTag(domain.TokenIdentifier) {
		open_paren := p.matchTag(domain.TokenLeftParen)

		if open_paren {
			p.expect(domain.TokenLeftParen)
		}
		if p.matchTag(domain.TokenLambda) {
			node = p.parse_abstraction()
		} else {
			node = p.parse_application()
		}
		if open_paren {
			p.expect(domain.TokenRightParen)
		}

	} else {
		node = p.parse_identifier()
	}
	return node
}

func (p *parser) parse_identifier() ast.Sexpr {
	identifier := p.src.Lexeme(p.current)
	p.next()
	return ast.S(
		domain.NodeIdentifier,
		identifier,
	)
}

func (p *parser) parse_application() ast.Sexpr {
	return ast.S(
		domain.NodeApplication,
		p.parse_term(),
		p.parse_term(),
	)
}

func (p *parser) parse_abstraction() ast.Sexpr {
	p.expect(domain.TokenLambda)
	identifier := p.parse_identifier()
	p.expect(domain.TokenDot)
	term := p.parse_term()
	return ast.S(
		domain.NodeAbstraction,
		identifier,
		term,
	)
}
