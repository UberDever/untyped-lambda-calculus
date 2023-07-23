package domain

import "math"

type NodeId int

const NodeInvalid = math.MinInt

const (
	NodeIdentifier = iota
	NodeApplication
	NodeAbstraction
)
