package middle

import (
	"lambda/ast"
	"lambda/domain"
	"lambda/util"
)

type deBruijnContext struct {
	abstraction_vars  util.Stack[string]
	free_vars_context map[string]int
}

func ToDeBruijn(namedAST ast.AST) (deBruijn ast.AST) {
	ctx := deBruijnContext{
		free_vars_context: make(map[string]int),
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

	onEnter := func(ast *ast.AST, node_id domain.NodeId) {
		node := ast.Node(node_id)
		switch node.Tag {
		case domain.NodeNamedVariable:
			id := ast.NamedVariableNode(node).Name
			index := abs_var_id(id)
			if index == domain.NodeNull {
				index = free_var_id(id)
			}

		case domain.NodeApplication:
		case domain.NodeAbstraction:
			n := ast.AbstractionNode(node)
			id := ast.SourceCode().Lexeme(domain.TokenId(n.Bound))
			ctx.abstraction_vars.Push(id)
		default:
			panic("Unreachable")
		}
	}

	onExit := func(ast *ast.AST, node_id domain.NodeId) {
		node := ast.Node(node_id)
		switch node.Tag {
		case domain.NodeNamedVariable:
		case domain.NodeApplication:
		case domain.NodeAbstraction:
			ctx.abstraction_vars.Pop()
		default:
			panic("Unreachable")
		}
	}

	namedAST.TraversePreorder(onEnter, onExit)

	return namedAST
}
