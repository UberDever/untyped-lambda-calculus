package debruijn

import (
	"lambda/ast/ast"
	"lambda/ast/tree"
	"lambda/syntax/source"
	"lambda/util"
)

type DeBruijnResult struct {
	Tree          tree.Tree
	VariableNames map[int]string
}

func ToDeBruijn(source_code source.SourceCode, tree_with_names tree.Tree) DeBruijnResult {
	abstraction_vars := util.NewStack[string]()
	free_vars_context := make(map[string]int)
	indicies := util.NewStack[tree.NodeId]()
	node_ids := util.NewStack[tree.NodeId]()
	variable_names := make(map[int]string)

	abs_var_id := func(variable string) tree.NodeId {
		vars := abstraction_vars.Values()
		// traverse in reverse order to encounter variable of closest lambda abstraction
		abstractions_encountered := len(vars) - 1
		for i := abstractions_encountered; i >= 0; i-- {
			id := vars[i]
			if id == variable {
				return tree.NodeId(abstractions_encountered - i)
			}
		}
		return tree.NodeNull
	}

	free_var_id := func(variable string) tree.NodeId {
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
		return tree.NodeId(free_id)
	}

	nodes := make([]tree.Node, 0)
	add_node := func(node tree.Node) tree.NodeId {
		nodes = append(nodes, node)
		return tree.NodeId(len(nodes) - 1)
	}

	onEnter := func(t tree.Tree, node_id tree.NodeId) {
		node := t.Node(node_id)
		switch node.Tag {
		case tree.NodeNamedVariable:
			typed_node := ast.ToNamedVariableNode(source_code, t, node)
			id := typed_node.Name
			index := abs_var_id(id)
			if index == tree.NodeNull {
				index = free_var_id(id)
			}
			variable_names[int(index)] = id
			indicies.Push(index)
		case tree.NodeApplication:
			break
		case tree.NodeAbstraction:
			typed_node := ast.ToAbstractionNode(t, node)
			bound := t.Node(typed_node.Bound())
			bound_node := ast.ToNamedVariableNode(source_code, t, bound)
			abstraction_vars.Push(bound_node.Name)
		default:
			panic("Unreachable")
		}
	}

	onExit := func(t tree.Tree, node_id tree.NodeId) {
		node := t.Node(node_id)
		token := node.Token
		switch node.Tag {
		case tree.NodeNamedVariable:
			index := indicies.ForcePop()
			id := add_node(tree.Node{
				Tag:   tree.NodeIndexVariable,
				Token: token,
				Lhs:   index,
				Rhs:   tree.NodeNull})
			node_ids.Push(id)
		case tree.NodeApplication:
			rhs := node_ids.ForcePop()
			lhs := node_ids.ForcePop()
			id := add_node(tree.Node{
				Tag:   tree.NodeApplication,
				Token: token,
				Lhs:   lhs,
				Rhs:   rhs})
			node_ids.Push(id)
		case tree.NodeAbstraction:
			body := node_ids.ForcePop()
			_ = node_ids.ForcePop() // variable (don't need named variable anymore)
			id := add_node(tree.Node{
				Tag:   tree.NodePureAbstraction,
				Token: token,
				Lhs:   body,
				Rhs:   tree.NodeNull})
			node_ids.Push(id)
			abstraction_vars.Pop()
		default:
			panic("Unreachable")
		}
	}

	ast.TraversePreorder(tree_with_names, tree_with_names.RootId(), onEnter, onExit)
	root := tree.NodeId(len(nodes) - 1)

	return DeBruijnResult{
		Tree:          tree.NewTree(root, nodes),
		VariableNames: variable_names,
	}
}
