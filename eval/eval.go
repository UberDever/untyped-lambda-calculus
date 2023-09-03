package eval

import (
	"lambda/ast/ast"
	"lambda/ast/tree"
)

func replicate_subtree(t *tree.MutableTree, root tree.NodeId) (new_root tree.NodeId) {
	nodes := t.Nodes()

	var aux func(tree.NodeId) tree.NodeId
	aux = func(r tree.NodeId) tree.NodeId {
		node := t.Node(r)
		iterable := ast.NewNodeIterable(node)
		lhs, rhs := iterable.Children()
		if lhs != tree.NodeNull {
			node.Lhs = aux(lhs)
		}
		if rhs != tree.NodeNull {
			node.Rhs = aux(rhs)
		}
		nodes = append(nodes, node)
		return tree.NodeId(len(nodes) - 1)
	}

	new_root = aux(root)
	t.SetNodes(nodes)
	return
}

// gc baby (stop the world, mark and sweep)
func collect_garbage(t *tree.MutableTree, root tree.NodeId) {
	nodes := make([]tree.Node, 0, len(t.Nodes())/4+1)

	var aux func(tree.NodeId) tree.NodeId
	aux = func(r tree.NodeId) tree.NodeId {
		cur := t.Node(r)
		nodes = append(nodes, cur)
		last := len(nodes) - 1

		lhs, rhs := ast.NewNodeIterable(cur).Children()
		if lhs != tree.NodeNull {
			nodes[last].Lhs = aux(lhs)
		}
		if rhs != tree.NodeNull {
			nodes[last].Rhs = aux(rhs)
		}

		return tree.NodeId(last)
	}
	new_root := aux(root)
	t.Tree = tree.NewTree(new_root, nodes)
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
			expr_cloned := t.Node(replicate_subtree(t, expr))
			t.SetNode(in, expr_cloned)
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

func find_redex_whnf(t tree.Tree, expr tree.NodeId) tree.NodeId {
	expr_node := t.Node(expr)
	if expr_node.Tag == tree.NodeApplication {
		v := ast.ToApplicationNode(t, expr_node)
		cur := v.Lhs()
		cur_node := t.Node(cur)
		for cur_node.Tag == tree.NodeApplication {
			cur_app := ast.ToApplicationNode(t, cur_node)

			new_cur := cur_app.Lhs()
			expr = cur
			cur = new_cur
			cur_node = t.Node(cur)
		}
		if cur_node.Tag == tree.NodePureAbstraction {
			return expr
		}
	}
	return tree.NodeNull
}

func Eval(log_eval func(t tree.Tree), in_tree tree.Tree, root tree.NodeId) tree.Tree {

	t := tree.NewMutableTree(in_tree.Clone())
	for reductions_count := 1; true; reductions_count++ {
		app_id := find_redex_whnf(t.Tree, root)
		if app_id == tree.NodeNull {
			break
		}

		log_eval(t.Tree)

		app := ast.ToApplicationNode(t.Tree, t.Node(app_id))
		lambda := ast.ToPureAbstractionNode(t.Tree, t.Node(app.Lhs()))
		app_rhs := app.Rhs()
		lambda_body := lambda.Body()

		shift_indicies(&t, app_rhs, 0, 1)
		substitute(&t, lambda_body, app_rhs, 0)
		shift_indicies(&t, lambda_body, 0, -1)

		t.SetNode(app_id, t.Node(lambda_body))
		if reductions_count%10 == 0 {
			// collect_garbage(&t, root)
		}
	}
	// collect_garbage(&t, root)
	return t.Tree
}
