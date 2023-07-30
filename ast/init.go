package ast

import (
	"fmt"
	"lambda/domain"
)

func init() {
	for i := 0; i < int(domain.NodeMax); i++ {
		{
			h := domain.NodeConstructor[i]
			if h == nil {
				panic(fmt.Sprintf("Node tag to constructor table is not full, failed at %d", i))
			}
		}
		{
			h := NodeChildren[i]
			if h == nil {
				panic(fmt.Sprintf("Node tag to children table is not full, failed at %d", i))
			}
		}
		{
			h := NodeString[i]
			if h == nil {
				panic(fmt.Sprintf("Node tag to string table is not full, failed at %d", i))
			}
		}
	}
}
