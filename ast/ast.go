package ast

import (
	"fmt"
	"lambda/domain"
	"lambda/syntax/source"
	"strings"
)

type AST struct {
	src   *source.SourceCode
	root  domain.NodeId
	nodes []domain.Node
}

func NewAST(src *source.SourceCode, root domain.NodeId, nodes []domain.Node) AST {
	return AST{src: src, root: root, nodes: nodes}
}

func (ast AST) SourceCode() *source.SourceCode {
	return ast.src
}

func (ast AST) Node(id domain.NodeId) domain.Node {
	return ast.nodes[int(id)]
}

func (ast AST) Print() string {
	str := strings.Builder{}
	onEnter := func(ast *AST, id domain.NodeId) {
		node := ast.Node(id)
		if node.Tag != domain.NodeIndexVariable &&
			node.Tag != domain.NodeNamedVariable {
			str.WriteByte('(')
		} else {
			str.WriteByte(' ')
		}
		str.WriteString(NodeString[node.Tag](ast, id))
	}
	onExit := func(ast *AST, id domain.NodeId) {
		node := ast.Node(id)
		if node.Tag != domain.NodeIndexVariable &&
			node.Tag != domain.NodeNamedVariable {
			str.WriteByte(')')
		}
	}
	ast.TraversePreorder(onEnter, onExit)

	return str.String()
}

type NodeAction = func(*AST, domain.NodeId)

func (ast AST) TraversePreorder(onEnter, onExit NodeAction) {
	ast.traversePreorder(onEnter, onExit, ast.root)
}

func (ast AST) traversePreorder(onEnter, onExit NodeAction, id domain.NodeId) {
	if id == domain.NodeNull {
		return
	}
	onEnter(&ast, id)
	defer onExit(&ast, id)
	n := ast.Node(id)
	lhs, rhs := NodeChildren[n.Tag](&ast, id)
	ast.traversePreorder(onEnter, onExit, lhs)
	ast.traversePreorder(onEnter, onExit, rhs)
}

func (ast AST) TraversePostorder(onEnter NodeAction) {
	ast.traversePostorder(onEnter, ast.root)
}

func (ast AST) traversePostorder(onEnter NodeAction, id domain.NodeId) {
	if id == domain.NodeNull {
		return
	}
	n := ast.Node(id)
	lhs, rhs := NodeChildren[n.Tag](&ast, id)
	ast.traversePostorder(onEnter, lhs)
	ast.traversePostorder(onEnter, rhs)
	onEnter(&ast, id)
}

var NodeChildren = [...]func(*AST, domain.NodeId) (domain.NodeId, domain.NodeId){
	domain.NodeNamedVariable:   NamedVariableNode_Children,
	domain.NodeApplication:     ApplicationNode_Children,
	domain.NodeAbstraction:     AbstractionNode_Children,
	domain.NodeIndexVariable:   IndexVariableNode_Children,
	domain.NodePureAbstraction: PureAbstractionNode_Children,
}

var NodeString = [...]func(*AST, domain.NodeId) string{
	domain.NodeNamedVariable:   NamedVariableNode_String,
	domain.NodeApplication:     ApplicationNode_String,
	domain.NodeAbstraction:     AbstractionNode_String,
	domain.NodeIndexVariable:   IndexVariableNode_String,
	domain.NodePureAbstraction: PureAbstractionNode_String,
}

type NamedVariableNode struct {
	Name string
}

func (ast *AST) NamedVariableNode(node domain.Node) NamedVariableNode {
	return NamedVariableNode{Name: ast.src.Lexeme(domain.TokenId(node.Lhs))}
}

func NamedVariableNode_Children(ast *AST, id domain.NodeId) (domain.NodeId, domain.NodeId) {
	return domain.NodeNull, domain.NodeNull
}

func NamedVariableNode_String(ast *AST, id domain.NodeId) string {
	n := ast.NamedVariableNode(ast.nodes[id])
	return n.Name
}

type ApplicationNode struct {
	Lhs, Rhs domain.NodeId
}

func (ast *AST) ApplicationNode(node domain.Node) ApplicationNode {
	return ApplicationNode{
		Lhs: node.Lhs, Rhs: node.Rhs,
	}
}

func ApplicationNode_Children(ast *AST, id domain.NodeId) (domain.NodeId, domain.NodeId) {
	n := ast.ApplicationNode(ast.nodes[id])
	return n.Lhs, n.Rhs
}

func ApplicationNode_String(ast *AST, id domain.NodeId) string {
	return ""
}

type AbstractionNode struct {
	Bound, Body domain.NodeId
}

func (ast *AST) AbstractionNode(node domain.Node) AbstractionNode {
	return AbstractionNode{
		Bound: node.Lhs,
		Body:  node.Rhs,
	}
}

func AbstractionNode_Children(ast *AST, id domain.NodeId) (domain.NodeId, domain.NodeId) {
	n := ast.AbstractionNode(ast.nodes[id])
	return n.Bound, n.Body
}

func AbstractionNode_String(ast *AST, id domain.NodeId) string {
	return "λ"
}

type IndexVariableNode struct {
	Index int
}

func (ast *AST) IndexVariableNode(node domain.Node) IndexVariableNode {
	return IndexVariableNode{Index: int(node.Lhs)}
}

func IndexVariableNode_Children(ast *AST, id domain.NodeId) (domain.NodeId, domain.NodeId) {
	return domain.NodeNull, domain.NodeNull
}

func IndexVariableNode_String(ast *AST, id domain.NodeId) string {
	n := ast.IndexVariableNode(ast.nodes[id])
	return fmt.Sprintf("%d", n.Index)
}

type PureAbstractionNode struct {
	Body domain.NodeId
}

func (ast *AST) PureAbstractionNode(node domain.Node) PureAbstractionNode {
	return PureAbstractionNode{
		Body: node.Lhs,
	}
}

func PureAbstractionNode_Children(ast *AST, id domain.NodeId) (domain.NodeId, domain.NodeId) {
	n := ast.PureAbstractionNode(ast.nodes[id])
	return n.Body, domain.NodeNull
}

func PureAbstractionNode_String(ast *AST, id domain.NodeId) string {
	return "λ"
}
