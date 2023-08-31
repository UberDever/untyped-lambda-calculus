package tree

import (
	"lambda/syntax/source"
	"math"
)

type NodeId int

const NodeInvalid NodeId = math.MinInt
const NodeNull NodeId = -1

const (
	NodeNamedVariable NodeId = iota
	NodeApplication
	NodeAbstraction
	NodeIndexVariable
	NodePureAbstraction
	NodeMax
)

type Node struct {
	Tag      NodeId
	Token    source.TokenId
	Lhs, Rhs NodeId
}

func NewNodeInvalid() Node {
	return Node{Tag: NodeInvalid, Token: source.TokenInvalid, Lhs: NodeInvalid, Rhs: NodeInvalid}
}

type Tree struct {
	root  NodeId
	nodes []Node
}

func NewTree(root NodeId, nodes []Node) Tree {
	return Tree{root: root, nodes: nodes}
}

func (ast Tree) Node(id NodeId) Node {
	return ast.nodes[int(id)]
}

func (ast Tree) Root() NodeId {
	return ast.root
}
