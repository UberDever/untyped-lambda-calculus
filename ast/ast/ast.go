package ast

import (
	"fmt"
	"lambda/ast/tree"
	"lambda/syntax/source"
	"strings"
)

type NodeIterable interface {
	Children() (tree.NodeId, tree.NodeId)
}

func NewNodeIterable(node tree.Node) NodeIterable {
	switch node.Tag {
	case tree.NodeNamedVariable:
		return named_variable_node{n: node}
	case tree.NodeApplication:
		return application_node{n: node}
	case tree.NodeAbstraction:
		return abstraction_node{n: node}
	case tree.NodeIndexVariable:
		return index_variable_node{n: node}
	case tree.NodePureAbstraction:
		return pure_abstraction_node{n: node}
	}
	panic("Unreachable")
}

func NewNodeStringer(src source.SourceCode, t tree.Tree, node tree.Node) fmt.Stringer {
	switch node.Tag {
	case tree.NodeNamedVariable:
		return ToNamedVariableNode(src, t, node)
	case tree.NodeApplication:
		return ToApplicationNode(t, node)
	case tree.NodeAbstraction:
		return ToAbstractionNode(t, node)
	case tree.NodeIndexVariable:
		return ToIndexVariableNode(t, node)
	case tree.NodePureAbstraction:
		return ToPureAbstractionNode(t, node)
	}
	panic("Unreachable")
}

type named_variable_node struct {
	n    tree.Node
	Name string
}

type application_node struct {
	n tree.Node
}

type abstraction_node struct {
	n tree.Node
}

type index_variable_node struct {
	n tree.Node
}

type pure_abstraction_node struct {
	n tree.Node
}

func ToNamedVariableNode(src source.SourceCode, tree tree.Tree, node tree.Node) named_variable_node {
	return named_variable_node{
		n:    node,
		Name: src.Lexeme(node.Token),
	}
}

func (n named_variable_node) String() string {
	return n.Name
}

func (n named_variable_node) Children() (tree.NodeId, tree.NodeId) {
	return tree.NodeNull, tree.NodeNull
}

func ToApplicationNode(tree tree.Tree, node tree.Node) application_node {
	return application_node{
		n: node,
	}
}

func (n application_node) String() string {
	return ""
}

func (n application_node) Children() (tree.NodeId, tree.NodeId) {
	return n.Lhs(), n.Rhs()
}

func (n application_node) Lhs() tree.NodeId {
	return n.n.Lhs
}

func (n application_node) Rhs() tree.NodeId {
	return n.n.Rhs
}

func ToAbstractionNode(tree tree.Tree, node tree.Node) abstraction_node {
	return abstraction_node{
		n: node,
	}
}

func (n abstraction_node) String() string {
	return "λ"
}

func (n abstraction_node) Children() (tree.NodeId, tree.NodeId) {
	return n.Bound(), n.Body()
}

func (n abstraction_node) Bound() tree.NodeId {
	return n.n.Lhs
}

func (n abstraction_node) Body() tree.NodeId {
	return n.n.Rhs
}

func ToIndexVariableNode(tree tree.Tree, node tree.Node) index_variable_node {
	return index_variable_node{
		n: node,
	}
}

func (n index_variable_node) String() string {
	return fmt.Sprintf("%d", n.n.Lhs)
}

func (n index_variable_node) Children() (tree.NodeId, tree.NodeId) {
	return tree.NodeNull, tree.NodeNull
}

func (n index_variable_node) Index() int {
	return int(n.n.Lhs)
}

func (n index_variable_node) NameIndex() int {
	return int(n.n.Rhs)
}

func ToPureAbstractionNode(tree tree.Tree, node tree.Node) pure_abstraction_node {
	return pure_abstraction_node{
		n: node,
	}
}

func (n pure_abstraction_node) String() string {
	return "λ"
}

func (n pure_abstraction_node) Children() (tree.NodeId, tree.NodeId) {
	return n.Body(), tree.NodeNull
}

func (n pure_abstraction_node) Body() tree.NodeId {
	return n.n.Lhs
}

type NodeAction = func(tree.Tree, tree.NodeId)

func TraversePreorder(tree tree.Tree, onEnter, onExit NodeAction) {
	traversePreorder(tree, onEnter, onExit, tree.RootId())
}

func traversePreorder(t tree.Tree, onEnter, onExit NodeAction, id tree.NodeId) {
	if id == tree.NodeNull {
		return
	}
	onEnter(t, id)
	defer onExit(t, id)
	node := t.Node(id)
	iterable_node := NewNodeIterable(node)
	lhs, rhs := iterable_node.Children()
	traversePreorder(t, onEnter, onExit, lhs)
	traversePreorder(t, onEnter, onExit, rhs)
}

func Print(src source.SourceCode, in_tree tree.Tree) string {
	str := strings.Builder{}
	onEnter := func(t tree.Tree, id tree.NodeId) {
		node := t.Node(id)
		stringer_node := NewNodeStringer(src, t, node)
		if node.Tag != tree.NodeIndexVariable &&
			node.Tag != tree.NodeNamedVariable {
			str.WriteByte('(')
		} else {
			str.WriteByte(' ')
		}
		str.WriteString(stringer_node.String())
	}
	onExit := func(t tree.Tree, id tree.NodeId) {
		node := t.Node(id)
		if node.Tag != tree.NodeIndexVariable &&
			node.Tag != tree.NodeNamedVariable {
			str.WriteByte(')')
		}
	}
	TraversePreorder(in_tree, onEnter, onExit)

	return str.String()
}
