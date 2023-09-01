package eval

import (
	"lambda/ast/ast"
	"lambda/ast/sexpr"
	"lambda/ast/tree"
	"lambda/syntax/source"
	"lambda/util"
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

func find_redex(t tree.Tree, expr tree.NodeId) tree.NodeId {
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

func Eval(logger *util.Logger, source_code source.SourceCode, in_tree tree.Tree) tree.Tree {
	log_computation := func(t tree.MutableTree, id tree.NodeId) {
		old_root := t.RootId()
		t.SetRoot(id)
		tree := ast.Print(source_code, t.Tree)
		t.SetRoot(old_root)

		pretty := sexpr.Spaced(tree)
		logger.Add(util.NewMessage(util.Debug, 0, 0, "e", pretty))
	}

	t := tree.NewMutableTree(in_tree.Clone())

	for {
		app_id := find_redex(t.Tree, t.RootId())
		if app_id == tree.NodeNull {
			break
		}

		log_computation(t, app_id)

		app := ast.ToApplicationNode(t.Tree, t.Node(app_id))
		lambda := ast.ToPureAbstractionNode(t.Tree, t.Node(app.Lhs()))
		app_rhs := app.Rhs()
		lambda_body := lambda.Body()

		shift_indicies(&t, app_rhs, 0, 1)
		substitute(&t, lambda_body, app_rhs, 0)
		shift_indicies(&t, lambda_body, 0, -1)

		t.SetNode(app_id, t.Node(lambda_body))
		t.SetNode(app_rhs, tree.Node{
			Tag:   tree.NodeIndexVariable,
			Token: source.TokenInvalid,
			Lhs:   420,
			Rhs:   -1,
		})
	}
	return t.Tree
}
