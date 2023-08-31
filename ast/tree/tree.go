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

func (t Tree) Count() int {
	return len(t.nodes)
}

func (t Tree) Node(id NodeId) Node {
	return t.nodes[int(id)]
}

func (t Tree) Root() Node {
	return t.Node(t.root)
}

func (t Tree) RootId() NodeId {
	return t.root
}

func (t Tree) Clone() Tree {
	return Tree{
		root:  t.RootId(),
		nodes: append([]Node{}, t.nodes...),
	}
}

type MutableTree struct {
	Tree
}

func NewMutableTree(tree Tree) MutableTree {
	t := tree.Clone()
	return MutableTree{t}
}

func (t *MutableTree) SetRoot(root NodeId) {
	t.root = root
}

func (t *MutableTree) SetNode(id NodeId, node Node) {
	t.nodes[int(id)] = node
}

func (t *MutableTree) SetNodes(nodes []Node) {
	t.nodes = nodes
}

func (t *MutableTree) Nodes() []Node {
	return t.nodes
}
