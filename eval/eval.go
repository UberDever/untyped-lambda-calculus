package eval

import (
	"lambda/ast/ast"
	"lambda/ast/tree"
)

func replicate_subtree(t *tree.MutableTree, root tree.NodeId) (new_root tree.NodeId) {
	nodes := t.Nodes()

	var rec func(tree.NodeId) tree.NodeId
	rec = func(r tree.NodeId) tree.NodeId {
		node := t.Node(r)
		iterable := ast.NewNodeIterable(node)
		lhs, rhs := iterable.Children()
		if lhs != tree.NodeNull {
			node.Lhs = rec(lhs)
		}
		if rhs != tree.NodeNull {
			node.Rhs = rec(rhs)
		}
		nodes = append(nodes, node)
		return tree.NodeId(len(nodes) - 1)
	}

	new_root = rec(root)
	t.SetNodes(nodes)
	return
}

func shift_indicies(t *tree.MutableTree, in tree.NodeId, cutoff, amount int) {
	node := t.Node(in)
	switch node.Tag {
	case tree.NodeIndexVariable:
		v := ast.ToIndexVariableNode(t.Tree, node)
		index := v.Index()
		if index >= cutoff {
			index += amount
		}
		new_node := tree.Node{
			Tag:   node.Tag,
			Token: node.Token,
			Lhs:   tree.NodeId(index),
			Rhs:   node.Rhs,
		}
		t.SetNode(in, new_node)
	case tree.NodePureAbstraction:
		v := ast.ToPureAbstractionNode(t.Tree, node)
		shift_indicies(t, v.Body(), cutoff+1, amount)
	case tree.NodeApplication:
		v := ast.ToApplicationNode(t.Tree, node)
		shift_indicies(t, v.Lhs(), cutoff, amount)
		shift_indicies(t, v.Rhs(), cutoff, amount)
	default:
		panic("unreachable")
	}
}

func substitute(t *tree.MutableTree, in tree.NodeId, expr tree.NodeId, level int) {
	node := t.Node(in)
	switch node.Tag {
	case tree.NodeIndexVariable:
		v := ast.ToIndexVariableNode(t.Tree, node)
		index := v.Index()
		if index == level {
			expr_cloned := replicate_subtree(t, expr)
			node := t.Node(expr_cloned)
			t.SetNode(in, node)
		}
	case tree.NodePureAbstraction:
		v := ast.ToPureAbstractionNode(t.Tree, node)
		body := v.Body()
		shift_indicies(t, expr, 0, 1)
		substitute(t, body, expr, level+1)
	case tree.NodeApplication:
		v := ast.ToApplicationNode(t.Tree, node)
		substitute(t, v.Lhs(), expr, level)
		substitute(t, v.Rhs(), expr, level)
	default:
		panic("unreachable")
	}
}

func is_redex(t tree.Tree, expr tree.Node) bool {
	if expr.Tag == tree.NodeApplication {
		v := ast.ToApplicationNode(t, expr)
		lhs := t.Node(v.Lhs())
		if lhs.Tag == tree.NodePureAbstraction {
			return true
		}
	}
	return false
}

func Eval(in_tree tree.Tree) tree.Tree {
	t := tree.NewMutableTree(in_tree)

	for is_redex(t.Tree, t.Root()) {
		application := ast.ToApplicationNode(t.Tree, t.Root())
		lhs := t.Node(application.Lhs())
		abstraction := ast.ToPureAbstractionNode(t.Tree, lhs)
		body := abstraction.Body()
		rhs := application.Rhs()

		shift_indicies(&t, rhs, 0, 1)
		substitute(&t, body, rhs, 0)
		shift_indicies(&t, body, 0, -1)
		t.SetRoot(body)
	}
	return t.Clone()
}
