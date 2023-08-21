package middle

import (
	AST "lambda/ast"
	"lambda/domain"
	"lambda/util"
)

type deBruijnContext struct {
	abstraction_vars  util.Stack[string]
	free_vars_context map[string]int
	indicies          util.Stack[domain.NodeId]

	node_ids util.Stack[domain.NodeId]
}

func ToDeBruijn(namedAST AST.AST) AST.AST {
	ctx := deBruijnContext{
		abstraction_vars:  util.NewStack[string](),
		free_vars_context: make(map[string]int),
		indicies:          util.NewStack[domain.NodeId](),
		node_ids:          util.NewStack[domain.NodeId](),
	}

	abs_var_id := func(variable string) domain.NodeId {
		vars := ctx.abstraction_vars.Values()
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
		index, ok := ctx.free_vars_context[variable]
		var free_id int
		if !ok {
			free_id = len(ctx.free_vars_context)
			ctx.free_vars_context[variable] = free_id
		} else {
			free_id = index
		}
		abstractions_encountered := len(ctx.abstraction_vars.Values()) - 1
		free_id += abstractions_encountered + 1
		return domain.NodeId(free_id)
	}

	nodes := make([]domain.Node, 0)
	new_node := func(node domain.Node) domain.NodeId {
		nodes = append(nodes, node)
		return domain.NodeId(len(nodes) - 1)
	}

	onEnter := func(ast *AST.AST, node_id domain.NodeId) {
		node := ast.Node(node_id)
		switch node.Tag {
		case domain.NodeNamedVariable:
			id := ast.NamedVariableNode(node).Name
			index := abs_var_id(id)
			if index == domain.NodeNull {
				index = free_var_id(id)
			}
			ctx.indicies.Push(index)
		case domain.NodeApplication:
			break
		case domain.NodeAbstraction:
			bound := ast.Node(ast.AbstractionNode(node).Bound)
			ctx.abstraction_vars.Push(ast.NamedVariableNode(bound).Name)
		default:
			panic("Unreachable")
		}
	}

	onExit := func(ast *AST.AST, node_id domain.NodeId) {
		node := ast.Node(node_id)
		token := node.Token
		switch node.Tag {
		case domain.NodeNamedVariable:
			index := ctx.indicies.ForcePop()
			id := new_node(domain.NodeConstructor[domain.NodeIndexVariable](token, index, domain.NodeNull))
			ctx.node_ids.Push(id)
		case domain.NodeApplication:
			rhs := ctx.node_ids.ForcePop()
			lhs := ctx.node_ids.ForcePop()
			id := new_node(domain.NodeConstructor[domain.NodeApplication](token, lhs, rhs))
			ctx.node_ids.Push(id)
		case domain.NodeAbstraction:
			body := ctx.node_ids.ForcePop()
			_ = ctx.node_ids.ForcePop() // variable
			id := new_node(domain.NodeConstructor[domain.NodePureAbstraction](token, body, domain.NodeNull))
			ctx.node_ids.Push(id)
			ctx.abstraction_vars.Pop()
		default:
			panic("Unreachable")
		}
	}

	namedAST.TraversePreorder(onEnter, onExit)
	root := domain.NodeId(len(nodes) - 1)

	// return namedAST
	return AST.NewAST(namedAST.SourceCode(), root, nodes)
}
