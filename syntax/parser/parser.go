package parser

import (
	"fmt"
	"lambda/ast/tree"
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
	tokens := make([]source.Token, 0, 16)
	pos := 0
	line, col := 1, 0

	add_token := func(tag source.TokenId, length int) {
		start, end := pos, pos+length
		tokens = append(tokens, source.NewToken(tag, start, end, line, col))
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
		return r != source.TokenDotRune &&
			r != source.TokenLambdaRune &&
			r != source.TokenLambdaBackslashRune &&
			r != source.TokenLeftParenRune &&
			r != source.TokenRightParenRune &&
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
		case source.TokenDotRune:
			add_token(source.TokenDot, 1)
		case source.TokenLambdaRune:
			add_token(source.TokenLambda, 1)
		case source.TokenLambdaBackslashRune:
			add_token(source.TokenLambda, 1)
		case source.TokenLeftParenRune:
			add_token(source.TokenLeftParen, 1)
		case source.TokenRightParenRune:
			add_token(source.TokenRightParen, 1)
		default:
			length := identifier_length()
			if length > 0 {
				add_token(source.TokenIdentifier, length)
			}
		}
	}

	tokens = append(tokens, source.NewTokenEof())
	return source.NewSourceCode(filename, text, tokens)
}

type parser struct {
	src *source.SourceCode

	ast_nodes []tree.Node
	current   source.TokenId
	atEof     bool

	logger *util.Logger
}

func (p *parser) next() {
	for {
		p.current++
		c := p.src.Token(p.current)
		p.atEof = c.Tag == source.TokenEof
		return
	}
}

func (p *parser) matchTag(tag source.TokenId) bool {
	if p.atEof {
		return false
	}
	c := p.src.Token(p.current)
	return c.Tag == tag
}

// TODO: Refactor this crap to support tag and tag + lexeme reporting
func (p *parser) expect(tag source.TokenId, lexeme string) (ok bool) {
	report_error := func(expected, got string, line, col int) {
		message := fmt.Sprintf("\nExpected\n %s but got\n %s", expected, got)
		p.logger.Add(util.NewMessage(util.Fatal, line, col, p.src.Filename(), message))
	}

	expected := p.src.TraceToken(tag, lexeme, int(source.TokenEof), int(source.TokenEof))

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

func (p *parser) new_node(node tree.Node) tree.NodeId {
	p.ast_nodes = append(p.ast_nodes, node)
	return tree.NodeId(len(p.ast_nodes) - 1)
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

func (p *parser) parse_term() tree.NodeId {
	id := tree.NodeInvalid

	if !p.matchTag(source.TokenIdentifier) {
		open_paren := p.matchTag(source.TokenLeftParen)

		if open_paren {
			p.expect(source.TokenLeftParen, "")
		}
		if p.matchTag(source.TokenLambda) {
			id = p.parse_abstraction()
		} else {
			id = p.parse_application()
		}
		if open_paren {
			p.expect(source.TokenRightParen, "")
		}

	} else {
		id = p.parse_variable()
	}

	return id
}

func (p *parser) parse_variable() tree.NodeId {
	tag, token, lhs, rhs := tree.NodeInvalid, source.TokenInvalid, tree.NodeInvalid, tree.NodeInvalid

	tag = tree.NodeNamedVariable
	token = p.current
	identifier := p.src.Lexeme(token)
	if identifier == "let" {
		return p.parse_let_binding()
	}

	p.next()

	return p.new_node(tree.Node{
		Tag:   tag,
		Token: token,
		Lhs:   lhs,
		Rhs:   rhs})
}

func (p *parser) parse_application() tree.NodeId {
	tag, token, lhs, rhs := tree.NodeInvalid, source.TokenInvalid, tree.NodeInvalid, tree.NodeInvalid

	tag = tree.NodeApplication
	token = p.current
	lhs = p.parse_term()
	rhs = p.parse_term()

	return p.new_node(tree.Node{
		Tag:   tag,
		Token: token,
		Lhs:   lhs,
		Rhs:   rhs})
}

func (p *parser) parse_abstraction() tree.NodeId {
	tag, token, lhs, rhs := tree.NodeInvalid, source.TokenInvalid, tree.NodeInvalid, tree.NodeInvalid

	tag = tree.NodeAbstraction
	token = p.current
	p.expect(source.TokenLambda, "")
	lhs = p.parse_variable()
	p.expect(source.TokenDot, "")
	rhs = p.parse_term()

	return p.new_node(tree.Node{
		Tag:   tag,
		Token: token,
		Lhs:   lhs,
		Rhs:   rhs})
}

func (p *parser) parse_let_binding() tree.NodeId {

	token := p.current
	p.expect(source.TokenIdentifier, "let")
	bound := p.parse_variable()
	p.expect(source.TokenIdentifier, "=")
	value := p.parse_term()
	p.expect(source.TokenIdentifier, "in")
	expr := p.parse_term()

	absraction := p.new_node(tree.Node{
		Tag:   tree.NodeAbstraction,
		Token: token,
		Lhs:   bound,
		Rhs:   expr})

	return p.new_node(tree.Node{
		Tag:   tree.NodeApplication,
		Token: token,
		Lhs:   absraction,
		Rhs:   value})
}
