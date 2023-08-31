package eval

import (
	"lambda/ast/tree"
)

type eval_context struct {
}

func NewEvalContext() eval_context {
	return eval_context{}
}

func shift_indicies(expr tree.Tree, cutoff, amount int) tree.Tree {
	return expr
}

func substitute(in tree.Tree, expr tree.Tree, level int) tree.Tree {
	return in
}

func is_redex(expr tree.Node) bool {
	return false
}

func (c *eval_context) Eval(expr tree.Tree) tree.Tree {
	// new_root := c.eval_rec(expr)
	return expr
}

func (c *eval_context) eval_rec(expr tree.Node) tree.Node {
	if is_redex(expr) {
	}
	return expr
}

// type eval_context struct {
// 	stack           util.Stack[sexpr.Sexpr]
// 	bound_variables util.Set[string]
// 	free_variables  util.Set[string]
// }
//
// func NewEvalContext() eval_context {
// 	return eval_context{
// 		stack: util.NewStack[sexpr.Sexpr](),
// 		bound_variables: util.NewSet[string](func(lhs, rhs string) bool {
// 			return lhs == rhs
// 		}),
// 		free_variables: util.NewSet[string](func(lhs, rhs string) bool {
// 			return lhs == rhs
// 		}),
// 	}
// }
//
// func ToString(expr sexpr.Sexpr, pretty bool) string {
// 	if expr.IsAtom() {
// 		return sexpr.Pretty(expr.Print())
// 	}
//
// 	lambda_symbol := '\\'
// 	if pretty {
// 		lambda_symbol = 'λ'
// 	}
// 	eval_stack := util.NewStack[any]()
// 	eval := func() {
// 		n := eval_stack.ForcePop().(int)
// 		tag := tree.NodeId(n)
// 		switch tag {
// 		case tree.NodeIndexVariable:
// 			// nothing
// 		case tree.NodeApplication:
// 			lhs := eval_stack.ForcePop()
// 			rhs := eval_stack.ForcePop()
// 			application := fmt.Sprintf(`(%s %s)`, lhs, rhs)
// 			eval_stack.Push(application)
// 		case tree.NodeAbstraction:
// 			arg := eval_stack.ForcePop()
// 			body := eval_stack.ForcePop()
// 			abstraction := fmt.Sprintf(`(%c %s %s)`, lambda_symbol, arg, body)
// 			eval_stack.Push(abstraction)
// 		default:
// 			panic("unreachable")
// 		}
// 	}
// 	onEnter := func(s sexpr.Sexpr) {
// 		if s.IsAtom() {
// 			eval_stack.Push(s.Data())
// 		} else {
// 			eval()
// 		}
// 	}
// 	sexpr.TraversePostorder(expr, onEnter)
// 	eval()
// 	return eval_stack.ForcePop().(string)
// }
//
// func (c *eval_context) GetBound() util.Set[string] {
// 	return c.bound_variables
// }
//
// func (c *eval_context) GetFree() util.Set[string] {
// 	return c.free_variables
// }
//
// func (c *eval_context) Eval(expr sexpr.Sexpr) sexpr.Sexpr {
// 	onEnter := func(s sexpr.Sexpr) {
// 		if s.IsAtom() {
// 			c.stack.Push(s)
// 		} else {
// 			fmt.Printf("%s %v\n", ToString(s, true), c.bound_variables)
// 			c.eval()
// 		}
// 	}
// 	sexpr.TraversePostorder(expr, onEnter)
// 	fmt.Printf("%s %v\n", ToString(expr, true), c.bound_variables)
// 	c.eval()
// 	return c.stack.ForcePop()
// }
//
// func (c *eval_context) eval() {
// 	n := c.stack.ForcePop().Data().(int)
// 	tag := tree.NodeId(n)
// 	switch tag {
// 	case tree.NodeIndexVariable:
// 		str := c.stack.ForcePop()
// 		identifier := sexpr.S(
// 			tree.NodeIndexVariable,
// 			str,
// 		)
// 		c.free_variables.Add(str.Data().(string))
// 		c.stack.Push(identifier)
// 	case tree.NodeApplication:
// 		lhs := c.stack.ForcePop()
// 		rhs := c.stack.ForcePop()
//
// 		rest := lhs
// 		lhs_tag := tree.NodeId(sexpr.Car(rest).Data().(int))
// 		if lhs_tag == tree.NodeAbstraction {
// 			// 	rest = ast.Cdr(rest)
// 			// 	arg := ast.Car(rest)
// 			// 	rest = ast.Cdr(rest)
// 			// 	body := ast.Car(rest)
// 			//             c.bound_variables()
// 		} else {
// 			application := sexpr.S(
// 				tree.NodeApplication,
// 				lhs, rhs,
// 			)
// 			c.stack.Push(application)
// 		}
// 	case tree.NodeAbstraction:
// 		arg := c.stack.ForcePop()
// 		body := c.stack.ForcePop()
// 		abstraction := sexpr.S(
// 			tree.NodeAbstraction,
// 			arg, body,
// 		)
// 		name := sexpr.Car(sexpr.Cdr(arg))
// 		str := name.Data().(string)
// 		c.bound_variables.Add(str)
// 		c.free_variables.Remove(str)
// 		c.stack.Push(abstraction)
// 	default:
// 		panic("unreachable")
// 	}
// }

// func (c *eval_context) bound_variables(expr ast.Sexpr) {
// 	rest := expr
// 	tag := tree.NodeId(ast.Car(rest).Data().(int))
// 	rest = ast.Cdr(rest)
// 	switch tag {
// 	case tree.NodeAbstraction:
// 		arg := ast.Car(rest)
// 		c.current_bound.Add(arg.Data().(string))
// 		rest = ast.Cdr(rest)
// 		body := ast.Car(rest)
// 		c.bound_variables(body)
// 	}
// }
