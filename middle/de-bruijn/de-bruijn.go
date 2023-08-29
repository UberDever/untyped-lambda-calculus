package debruijn

import (
	"lambda/ast/ast"
	"lambda/ast/tree"
	"lambda/domain"
	"lambda/syntax/source"
	"lambda/util"
)

func ToDeBruijn(source_code *source.SourceCode, tree_with_names *tree.Tree) (root domain.NodeId, nodes []domain.Node) {
	abstraction_vars := util.NewStack[string]()
	free_vars_context := make(map[string]int)
	indicies := util.NewStack[domain.NodeId]()
	node_ids := util.NewStack[domain.NodeId]()

	abs_var_id := func(variable string) domain.NodeId {
		vars := abstraction_vars.Values()
		// traverse in reverse order to encounter variable of closest lambda abstraction
		abstractions_encountered := len(vars) - 1
		for i := abstractions_encountered; i >= 0; i-- {
			id := vars[i]
			if id == variable {
				return domain.NodeId(abstractions_encountered - i)
			}
		}
		return domain.NodeNull
	}

	free_var_id := func(variable string) domain.NodeId {
		index, ok := free_vars_context[variable]
		var free_id int
		if !ok {
			free_id = len(free_vars_context)
			free_vars_context[variable] = free_id
		} else {
			free_id = index
		}
		abstractions_encountered := len(abstraction_vars.Values()) - 1
		free_id += abstractions_encountered + 1
		return domain.NodeId(free_id)
	}

	nodes = make([]domain.Node, 0)
	add_node := func(node domain.Node) domain.NodeId {
		nodes = append(nodes, node)
		return domain.NodeId(len(nodes) - 1)
	}

	onEnter := func(tree *tree.Tree, node_id domain.NodeId) {
		node := tree.Node(node_id)
		switch node.Tag {
		case domain.NodeNamedVariable:
			typed_node := ast.NewNamedVariableNode(source_code, tree, node)
			id := typed_node.Name
			index := abs_var_id(id)
			if index == domain.NodeNull {
				index = free_var_id(id)
			}
			indicies.Push(index)
		case domain.NodeApplication:
			break
		case domain.NodeAbstraction:
			typed_node := ast.NewAbstractionNode(source_code, tree, node)
			bound := tree.Node(typed_node.Bound())
			bound_node := ast.NewNamedVariableNode(source_code, tree, bound)
			abstraction_vars.Push(bound_node.Name)
		default:
			panic("Unreachable")
		}
	}

	onExit := func(tree *tree.Tree, node_id domain.NodeId) {
		node := tree.Node(node_id)
		token := node.Token
		switch node.Tag {
		case domain.NodeNamedVariable:
			index := indicies.ForcePop()
			id := add_node(domain.Node{
				Tag:   domain.NodeIndexVariable,
				Token: token,
				Lhs:   index,
				Rhs:   domain.NodeNull})
			node_ids.Push(id)
		case domain.NodeApplication:
			rhs := node_ids.ForcePop()
			lhs := node_ids.ForcePop()
			id := add_node(domain.Node{
				Tag:   domain.NodeApplication,
				Token: token,
				Lhs:   lhs,
				Rhs:   rhs})
			node_ids.Push(id)
		case domain.NodeAbstraction:
			body := node_ids.ForcePop()
			_ = node_ids.ForcePop() // variable (don't need named variable anymore)
			id := add_node(domain.Node{
				Tag:   domain.NodePureAbstraction,
				Token: token,
				Lhs:   body,
				Rhs:   domain.NodeNull})
			node_ids.Push(id)
			abstraction_vars.Pop()
		default:
			panic("Unreachable")
		}
	}

	ast.TraversePreorder(source_code, tree_with_names, onEnter, onExit)
	root = domain.NodeId(len(nodes) - 1)

	return root, nodes
}
