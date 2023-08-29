package ast

import (
	"fmt"
	"lambda/ast/tree"
	"lambda/domain"
	"lambda/syntax/source"
	"strings"
)

type INode interface {
	String() string
	Children() (domain.NodeId, domain.NodeId)
}

func NewINode(src *source.SourceCode, tree *tree.Tree, node domain.Node) INode {
	switch node.Tag {
	case domain.NodeNamedVariable:
		return NewNamedVariableNode(src, tree, node)
	case domain.NodeApplication:
		return NewApplicationNode(src, tree, node)
	case domain.NodeAbstraction:
		return NewAbstractionNode(src, tree, node)
	case domain.NodeIndexVariable:
		return NewIndexVariableNode(src, tree, node)
	case domain.NodePureAbstraction:
		return NewPureAbstractionNode(src, tree, node)
	}
	panic("Unreachable")
}

type named_variable_node struct {
	n    domain.Node
	Name string
}

type application_node struct {
	n domain.Node
}

type abstraction_node struct {
	n domain.Node
}

type index_variable_node struct {
	n domain.Node
}

type pure_abstraction_node struct {
	n domain.Node
}

func NewNamedVariableNode(src *source.SourceCode, tree *tree.Tree, node domain.Node) named_variable_node {
	return named_variable_node{
		n:    node,
		Name: src.Lexeme(node.Token),
	}
}

func (n named_variable_node) String() string {
	return n.Name
}

func (n named_variable_node) Children() (domain.NodeId, domain.NodeId) {
	return domain.NodeNull, domain.NodeNull
}

func NewApplicationNode(src *source.SourceCode, tree *tree.Tree, node domain.Node) application_node {
	return application_node{
		n: node,
	}
}

func (n application_node) String() string {
	return ""
}

func (n application_node) Children() (domain.NodeId, domain.NodeId) {
	return n.Lhs(), n.Rhs()
}

func (n application_node) Lhs() domain.NodeId {
	return n.n.Lhs
}

func (n application_node) Rhs() domain.NodeId {
	return n.n.Rhs
}

func NewAbstractionNode(src *source.SourceCode, tree *tree.Tree, node domain.Node) abstraction_node {
	return abstraction_node{
		n: node,
	}
}

func (n abstraction_node) String() string {
	return "λ"
}

func (n abstraction_node) Children() (domain.NodeId, domain.NodeId) {
	return n.Bound(), n.Body()
}

func (n abstraction_node) Bound() domain.NodeId {
	return n.n.Lhs
}

func (n abstraction_node) Body() domain.NodeId {
	return n.n.Rhs
}

func NewIndexVariableNode(src *source.SourceCode, tree *tree.Tree, node domain.Node) index_variable_node {
	return index_variable_node{
		n: node,
	}
}

func (n index_variable_node) String() string {
	return fmt.Sprintf("%d", n.n.Lhs)
}

func (n index_variable_node) Children() (domain.NodeId, domain.NodeId) {
	return domain.NodeNull, domain.NodeNull
}

func (n index_variable_node) Index() int {
	return int(n.n.Lhs)
}

func (n index_variable_node) NameIndex() int {
	return int(n.n.Rhs)
}

func NewPureAbstractionNode(src *source.SourceCode, tree *tree.Tree, node domain.Node) pure_abstraction_node {
	return pure_abstraction_node{
		n: node,
	}
}

func (n pure_abstraction_node) String() string {
	return "λ"
}

func (n pure_abstraction_node) Children() (domain.NodeId, domain.NodeId) {
	return n.Body(), domain.NodeNull
}

func (n pure_abstraction_node) Body() domain.NodeId {
	return n.n.Lhs
}

type NodeAction = func(*tree.Tree, domain.NodeId)

func TraversePreorder(source_code *source.SourceCode, tree *tree.Tree, onEnter, onExit NodeAction) {
	traversePreorder(source_code, tree, onEnter, onExit, tree.Root())
}

func traversePreorder(source_code *source.SourceCode, tree *tree.Tree, onEnter, onExit NodeAction, id domain.NodeId) {
	if id == domain.NodeNull {
		return
	}
	onEnter(tree, id)
	defer onExit(tree, id)
	node := tree.Node(id)
	typed_node := NewINode(source_code, tree, node)
	lhs, rhs := typed_node.Children()
	traversePreorder(source_code, tree, onEnter, onExit, lhs)
	traversePreorder(source_code, tree, onEnter, onExit, rhs)
}

func Print(src *source.SourceCode, in_tree *tree.Tree) string {
	str := strings.Builder{}
	onEnter := func(tree *tree.Tree, id domain.NodeId) {
		node := tree.Node(id)
		inode := NewINode(src, tree, node)
		if node.Tag != domain.NodeIndexVariable &&
			node.Tag != domain.NodeNamedVariable {
			str.WriteByte('(')
		} else {
			str.WriteByte(' ')
		}
		str.WriteString(inode.String())
	}
	onExit := func(tree *tree.Tree, id domain.NodeId) {
		node := tree.Node(id)
		if node.Tag != domain.NodeIndexVariable &&
			node.Tag != domain.NodeNamedVariable {
			str.WriteByte(')')
		}
	}
	TraversePreorder(src, in_tree, onEnter, onExit)

	return str.String()
}
