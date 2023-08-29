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
