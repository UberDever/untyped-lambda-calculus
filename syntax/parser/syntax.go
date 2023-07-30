package parser

import (
	"fmt"
	"lambda/ast"
	"lambda/domain"
	"lambda/syntax/source"
	"unicode"

	"golang.org/x/exp/utf8string"
)

type tokenizer struct {
	logger *domain.Logger
}

func NewTokenizer(logger *domain.Logger) tokenizer {
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
	depth     int
	current   domain.TokenId
	atEof     bool

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
		p.logger.Add(domain.NewMessage(domain.Fatal, c.Line, c.Col, p.src.Filename(), message))
		p.atEof = true
		return
	}

	p.next()
	return
}

func (p *parser) new_node(node domain.Node) domain.NodeId {
	p.ast_nodes = append(p.ast_nodes, node)
	return domain.NodeId(len(p.ast_nodes) - 1)
}

func NewParser(logger *domain.Logger) parser {
	return parser{logger: logger, depth: 1}
}

func (p *parser) Parse(src *source.SourceCode) ast.AST {
	p.src = src

	root := p.parse_term()
	return ast.NewAST(src, root, p.ast_nodes)
}

func (p *parser) parse_term() domain.NodeId {
	id := domain.NodeInvalid

	if !p.matchTag(domain.TokenIdentifier) {
		open_paren := p.matchTag(domain.TokenLeftParen)

		if open_paren {
			p.expect(domain.TokenLeftParen)
		}
		if p.matchTag(domain.TokenLambda) {
			id = p.parse_abstraction()
		} else {
			id = p.parse_application()
		}
		if open_paren {
			p.expect(domain.TokenRightParen)
		}

	} else {
		id = p.parse_variable()
	}

	return id
}

func (p *parser) parse_variable() domain.NodeId {
	tag, token, lhs, rhs := domain.NodeInvalid, domain.TokenInvalid, domain.NodeInvalid, domain.NodeInvalid

	tag = domain.NodeVariable
	token = p.current
	lhs = domain.NodeId(p.depth)
	rhs = domain.NodeNull
	p.next()

	return p.new_node(domain.NodeConstructor[tag](token, lhs, rhs))
}

func (p *parser) parse_application() domain.NodeId {
	tag, token, lhs, rhs := domain.NodeInvalid, domain.TokenInvalid, domain.NodeInvalid, domain.NodeInvalid

	tag = domain.NodeApplication
	token = p.current
	lhs = p.parse_term()
	rhs = p.parse_term()

	return p.new_node(domain.NodeConstructor[tag](token, lhs, rhs))
}

func (p *parser) parse_abstraction() domain.NodeId {
	tag, token, lhs, rhs := domain.NodeInvalid, domain.TokenInvalid, domain.NodeInvalid, domain.NodeInvalid

	tag = domain.NodeAbstraction
	token = p.current
	p.depth++
	p.expect(domain.TokenLambda)
	lhs = p.parse_variable()
	p.expect(domain.TokenDot)
	rhs = p.parse_term()
	p.depth--

	return p.new_node(domain.NodeConstructor[tag](token, lhs, rhs))
}
