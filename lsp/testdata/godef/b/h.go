package b

import . "github.com/charlievieth/xtools/lsp/godef/a"

func _() {
	// variable of type a.A
	var _ A //@mark(AVariable, "_"),hover("_", AVariable)

	AStuff() //@hover("AStuff", AStuff)
}
