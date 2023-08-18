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

func (ast AST) Node(id domain.NodeId) domain.Node {
	return ast.nodes[int(id)]
}

func (ast AST) Print() string {
	str := strings.Builder{}
	onEnter := func(ast *AST, id domain.NodeId) {
		node := ast.Node(id)
		if node.Tag != domain.NodeVariable {
			str.WriteByte('(')
		} else {
			str.WriteByte(' ')
		}
		str.WriteString(NodeString[node.Tag](ast, id))
	}
	onExit := func(ast *AST, id domain.NodeId) {
		node := ast.Node(id)
		if node.Tag != domain.NodeVariable {
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

var NodeChildren = [...]func(*AST, domain.NodeId) (domain.NodeId, domain.NodeId){
	domain.NodeVariable:    VariableNode_Children,
	domain.NodeApplication: ApplicationNode_Children,
	domain.NodeAbstraction: AbstractionNode_Children,
}

var NodeString = [...]func(*AST, domain.NodeId) string{
	domain.NodeVariable:    VariableNode_String,
	domain.NodeApplication: ApplicationNode_String,
	domain.NodeAbstraction: AbstractionNode_String,
}

type VariableNode struct {
	Index int
}

func (ast *AST) VariableNode(node domain.Node) VariableNode {
	return VariableNode{Index: int(node.Lhs)}
}

func VariableNode_Children(ast *AST, id domain.NodeId) (domain.NodeId, domain.NodeId) {
	return domain.NodeNull, domain.NodeNull
}

func VariableNode_String(ast *AST, id domain.NodeId) string {
	n := ast.VariableNode(ast.nodes[id])
	return fmt.Sprintf("%d", n.Index)
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
	Body domain.NodeId
}

func (ast *AST) AbstractionNode(node domain.Node) AbstractionNode {
	return AbstractionNode{
		Body: node.Rhs,
	}
}

func AbstractionNode_Children(ast *AST, id domain.NodeId) (domain.NodeId, domain.NodeId) {
	n := ast.AbstractionNode(ast.nodes[id])
	return n.Body, domain.NodeNull
}

func AbstractionNode_String(ast *AST, id domain.NodeId) string {
	return "λ"
}
