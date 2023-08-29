package tree

import (
	"lambda/domain"
)

type Tree struct {
	root  domain.NodeId
	nodes []domain.Node
}

func NewTree(root domain.NodeId, nodes []domain.Node) Tree {
	return Tree{root: root, nodes: nodes}
}

func (ast Tree) Node(id domain.NodeId) domain.Node {
	return ast.nodes[int(id)]
}

func (ast Tree) Root() domain.NodeId {
	return ast.root
}
