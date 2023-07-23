package ast

import (
	"lambda/util"
	"testing"
)

func TestEquality(t *testing.T) {
	cmp := func(lhs, rhs any) bool { return lhs == rhs }
	{
		lhs := Cons(Sexpr{"Cons"}, Sexpr{"Cell"})
		rhs := Cons(Sexpr{"Cons"}, Sexpr{"Cell"})
		if !Equals(lhs, rhs, cmp) {
			t.Errorf("Equality test failed:\n%s\n%s", lhs.Print(),
				rhs.Print())
		}
	}

	{
		lhs := S(1, 2, 3)
		rhs := S(1, 2, 3)
		if !Equals(lhs, rhs, cmp) {
			t.Errorf("Equality test failed:\n%s\n%s", lhs.Print(),
				rhs.Print())
		}
	}

	{
		lhs := S(1, S(2, "a"), 3)
		rhs := S(1, S(2, "a"), 3)
		if !Equals(lhs, rhs, cmp) {
			t.Errorf("Equality test failed:\n%s\n%s", lhs.Print(),
				rhs.Print())
		}
	}

	{
		lhs := S(1, nil)
		rhs := S(nil, 1)
		if Equals(lhs, rhs, cmp) {
			t.Errorf("Equality test failed:\n%s\n%s", lhs.Print(),
				rhs.Print())
		}
	}
}

func TestPrintDotted(t *testing.T) {
	{
		var l Sexpr = Sexpr{nil}
		expected := `nil`
		if l.PrintDotted() != expected {
			t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(l.PrintDotted(), expected))
		}
	}
	{
		l := Cons(Sexpr{1}, Sexpr{2})
		expected := `(1 . 2)`
		if l.PrintDotted() != expected {
			t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(l.PrintDotted(), expected))
		}
	}
	{
		l := S(1, 2, 3, 4, 5)
		expected := `(1 . (2 . (3 . (4 . (5 . nil)))))`
		if l.PrintDotted() != expected {
			t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(l.PrintDotted(), expected))
		}
	}
}

func TestPrint(t *testing.T) {
	{
		l := S()
		expected := `()`
		if l.Print() != expected {
			t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(l.Print(), expected))
		}
	}
	{
		l := S(1, 2, 3, 4, 5)
		expected := `(1 2 3 4 5 )`
		if l.Print() != expected {
			t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(l.Print(), expected))
		}
	}
	{
		l := S("A", S("B", "C"), "D", S("E", S("F", "G")))
		expected := `(A (B C )D (E (F G )))`
		if l.Print() != expected {
			t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(l.Print(), expected))
		}
	}
}

func TestMinify(t *testing.T) {
	{
		l := S()
		expected := `()`
		got := Minified(l.Print())
		if got != expected {
			t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(got, expected))
		}
	}
	{
		l := S(1, 2, 3, 4, 5)
		expected := `(1 2 3 4 5)`
		got := Minified(l.Print())
		if got != expected {
			t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(got, expected))
		}
	}
	{
		l := S("A", S("B", "C"), "D", S("E", S("F", "G")))
		expected := `(A(B C)D(E(F G)))`
		got := Minified(l.Print())
		if got != expected {
			t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(got, expected))
		}
	}
	{
		l := S("Long", "words", S("right", "here"))
		expected := `(Long words(right here))`
		got := Minified(l.Print())
		if got != expected {
			t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(got, expected))
		}
	}
}

func TestPrettify(t *testing.T) {
	{
		l := S()
		expected := `()`
		got := Pretty(l.Print())
		if got != expected {
			t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(got, expected))
		}
	}
	{
		l := Sexpr{"Hello"}
		expected := `Hello`
		got := Pretty(l.Print())
		if got != expected {
			t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(got, expected))
		}
	}
	{
		l := S(1, 2, 3, 4, 5)
		expected := `(1 2 3 4 5)`
		got := Pretty(l.Print())
		if got != expected {
			t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(got, expected))
		}
	}
	{
		l := S("A", S("B", "C"), "D", S("E", S("F", "G")))
		expected := `(A (B C) D (E (F G)))`
		got := Pretty(l.Print())
		if got != expected {
			t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(got, expected))
		}
	}
	{
		l := S("A", S(S("Sub", S("expr")), "C"), "D", S("E", S("F", "G")))
		expected := `(A ((Sub (expr)) C) D (E (F G)))`
		got := Pretty(l.Print())
		if got != expected {
			t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(got, expected))
		}
	}
	{
		l := S("Long", "words", S("right", "here"))
		expected := `(Long words (right here))`
		got := Pretty(l.Print())
		if got != expected {
			t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(got, expected))
		}
	}
	{
		l := S("this", "list", "contains", "very", "big", "amount", "of", "words", "as", "bad", "example", "of", "very", "long", "list")
		expected := `(this list contains very big amount of words as bad example of very long list)`
		got := Pretty(l.Print())
		if got != expected {
			t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(got, expected))
		}
	}
	{
		l := S("this", "list", "contains", "very", "big", "amount", "of", "words", S("this", "should", "go", "on", "next", "line"))
		expected := "(this list contains very big amount of words\n    (this should go on next line))"
		got := Pretty(l.Print())
		if got != expected {
			t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(got, expected))
		}
	}
}
