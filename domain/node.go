package domain

import (
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
	Token    TokenId
	Lhs, Rhs NodeId
}

func NewNodeInvalid() Node {
	return Node{Tag: NodeInvalid, Token: TokenInvalid, Lhs: NodeInvalid, Rhs: NodeInvalid}
}

func NewNamedVariableNode(token TokenId, lhs, rhs NodeId) Node {
	return Node{
		Tag:   NodeNamedVariable,
		Token: token,
		Lhs:   lhs,
		Rhs:   rhs,
	}
}

func NewApplicationNode(token TokenId, lhs, rhs NodeId) Node {
	return Node{
		Tag:   NodeApplication,
		Token: token,
		Lhs:   lhs,
		Rhs:   rhs,
	}
}

func NewAbstractionNode(token TokenId, lhs, rhs NodeId) Node {
	return Node{
		Tag:   NodeAbstraction,
		Token: token,
		Lhs:   lhs,
		Rhs:   rhs,
	}
}

func NewIndexVariableNode(token TokenId, lhs, rhs NodeId) Node {
	return Node{
		Tag:   NodeIndexVariable,
		Token: token,
		Lhs:   lhs,
		Rhs:   rhs,
	}
}

func NewPureAbstractionNode(token TokenId, lhs, rhs NodeId) Node {
	return Node{
		Tag:   NodePureAbstraction,
		Token: token,
		Lhs:   lhs,
		Rhs:   rhs,
	}
}

var NodeConstructor = [...]func(TokenId, NodeId, NodeId) Node{
	NodeNamedVariable:   NewNamedVariableNode,
	NodeApplication:     NewApplicationNode,
	NodeAbstraction:     NewAbstractionNode,
	NodeIndexVariable:   NewIndexVariableNode,
	NodePureAbstraction: NewPureAbstractionNode,
}
