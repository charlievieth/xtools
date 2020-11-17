package testy

import (
	"testing"

	sig "github.com/charlievieth/xtools/lsp/signature"
	"github.com/charlievieth/xtools/lsp/snippets"
)

func TestSomething(t *testing.T) { //@item(TestSomething, "TestSomething(t *testing.T)", "", "func")
	var x int //@mark(testyX, "x"),diag("x", "compiler", "x declared but not used", "error"),refs("x", testyX)
	a()       //@mark(testyA, "a")
}

func _() {
	_ = snippets.X(nil) //@signature("nil", "X(_ map[sig.Alias]types.CoolAlias) map[sig.Alias]types.CoolAlias", 0)
	var _ sig.Alias
}
