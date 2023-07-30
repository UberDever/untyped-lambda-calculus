package domain

import (
	"math"
)

type NodeId int

const NodeInvalid NodeId = math.MinInt
const NodeNull NodeId = -1

const (
	NodeVariable NodeId = iota
	NodeApplication
	NodeAbstraction
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

func NewVariableNode(token TokenId, lhs, rhs NodeId) Node {
	return Node{
		Tag:   NodeVariable,
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

var NodeConstructor = [...]func(TokenId, NodeId, NodeId) Node{
	NodeVariable:    NewVariableNode,
	NodeApplication: NewApplicationNode,
	NodeAbstraction: NewAbstractionNode,
}
