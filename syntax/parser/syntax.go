package parser

import (
	"fmt"
	"lambda/ast/tree"
	"lambda/domain"
	"lambda/syntax/source"
	"lambda/util"
	"unicode"

	"golang.org/x/exp/utf8string"
)

type tokenizer struct {
	logger *util.Logger
}

func NewTokenizer(logger *util.Logger) tokenizer {
	return tokenizer{logger: logger}
}

func (tok tokenizer) Tokenize(filename string, text utf8string.String) source.SourceCode {
	tokens := make([]domain.Token, 0, 16)
	pos := 0
	line, col := 1, 0

	add_token := func(tag domain.TokenId, length int) {
		start, end := pos, pos+length
		tokens = append(tokens, domain.NewToken(tag, start, end, line, col))
		pos = end
		col = end
	}

	skip_spaces := func() {
		for pos < text.RuneCount() {
			c := text.At(pos)
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
			r != domain.TokenLambdaBackslashRune &&
			r != domain.TokenLeftParenRune &&
			r != domain.TokenRightParenRune &&
			!unicode.IsSpace(r)
	}
	identifier_length := func() int {
		start, end := pos, pos
		for end < text.RuneCount() {
			c := text.At(end)
			if !identifier_rune(c) {
				return end - start
			}
			end++
		}
		return end - start
	}

	for {
		skip_spaces()
		if pos >= text.RuneCount() {
			break
		}

		switch text.At(pos) {
		case domain.TokenDotRune:
			add_token(domain.TokenDot, 1)
		case domain.TokenLambdaRune:
			add_token(domain.TokenLambda, 1)
		case domain.TokenLambdaBackslashRune:
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

	tokens = append(tokens, domain.NewTokenEof())
	return source.NewSourceCode(filename, text, tokens)
}

type parser struct {
	src *source.SourceCode

	ast_nodes []domain.Node
	current   domain.TokenId
	atEof     bool

	logger *util.Logger
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

// TODO: Refactor this crap to support tag and tag + lexeme reporting
func (p *parser) expect(tag domain.TokenId, lexeme string) (ok bool) {
	report_error := func(expected, got string, line, col int) {
		message := fmt.Sprintf("\nExpected\n %s but got\n %s", expected, got)
		p.logger.Add(util.NewMessage(util.Fatal, line, col, p.src.Filename(), message))
	}

	expected := p.src.TraceToken(tag, lexeme, int(domain.TokenEof), int(domain.TokenEof))

	if p.atEof {
		got := "EOF"
		report_error(expected, got, -1, -1)
		return
	}

	ok = p.matchTag(tag)
	if !ok {
		c := p.src.Token(p.current)
		got := p.src.TraceToken(c.Tag, p.src.Lexeme(p.current), c.Line, c.Col)
		report_error(expected, got, c.Line, c.Col)
		return
	}

	p.next()
	return
}

func (p *parser) new_node(node domain.Node) domain.NodeId {
	p.ast_nodes = append(p.ast_nodes, node)
	return domain.NodeId(len(p.ast_nodes) - 1)
}

func NewParser(logger *util.Logger) parser {
	return parser{
		logger: logger,
	}
}

func (p *parser) Parse(src *source.SourceCode) tree.Tree {
	p.src = src

	root := p.parse_term()
	if !p.atEof {
		message := "Unexpected EOF"
		// NOTE: line and column values here are handy to reporting, but I have removed them
		// from parser implementation, don't remember why
		p.logger.Add(util.NewMessage(util.Fatal, -1, -1, p.src.Filename(), message))
	}
	return tree.NewTree(root, p.ast_nodes)
}

func (p *parser) parse_term() domain.NodeId {
	id := domain.NodeInvalid

	if !p.matchTag(domain.TokenIdentifier) {
		open_paren := p.matchTag(domain.TokenLeftParen)

		if open_paren {
			p.expect(domain.TokenLeftParen, "")
		}
		if p.matchTag(domain.TokenLambda) {
			id = p.parse_abstraction()
		} else {
			id = p.parse_application()
		}
		if open_paren {
			p.expect(domain.TokenRightParen, "")
		}

	} else {
		id = p.parse_variable()
	}

	return id
}

func (p *parser) parse_variable() domain.NodeId {
	tag, token, lhs, rhs := domain.NodeInvalid, domain.TokenInvalid, domain.NodeInvalid, domain.NodeInvalid

	tag = domain.NodeNamedVariable
	token = p.current
	identifier := p.src.Lexeme(token)
	if identifier == "let" {
		return p.parse_let_binding()
	}

	p.next()

	return p.new_node(domain.Node{
		Tag:   tag,
		Token: token,
		Lhs:   lhs,
		Rhs:   rhs})
}

func (p *parser) parse_application() domain.NodeId {
	tag, token, lhs, rhs := domain.NodeInvalid, domain.TokenInvalid, domain.NodeInvalid, domain.NodeInvalid

	tag = domain.NodeApplication
	token = p.current
	lhs = p.parse_term()
	rhs = p.parse_term()

	return p.new_node(domain.Node{
		Tag:   tag,
		Token: token,
		Lhs:   lhs,
		Rhs:   rhs})
}

func (p *parser) parse_abstraction() domain.NodeId {
	tag, token, lhs, rhs := domain.NodeInvalid, domain.TokenInvalid, domain.NodeInvalid, domain.NodeInvalid

	tag = domain.NodeAbstraction
	token = p.current
	p.expect(domain.TokenLambda, "")
	lhs = p.parse_variable()
	p.expect(domain.TokenDot, "")
	rhs = p.parse_term()

	return p.new_node(domain.Node{
		Tag:   tag,
		Token: token,
		Lhs:   lhs,
		Rhs:   rhs})
}

func (p *parser) parse_let_binding() domain.NodeId {

	token := p.current
	p.expect(domain.TokenIdentifier, "let")
	bound := p.parse_variable()
	p.expect(domain.TokenIdentifier, "=")
	value := p.parse_term()
	p.expect(domain.TokenIdentifier, "in")
	expr := p.parse_term()

	absraction := p.new_node(domain.Node{
		Tag:   domain.NodeAbstraction,
		Token: token,
		Lhs:   bound,
		Rhs:   expr})

	return p.new_node(domain.Node{
		Tag:   domain.NodeApplication,
		Token: token,
		Lhs:   absraction,
		Rhs:   value})
}
