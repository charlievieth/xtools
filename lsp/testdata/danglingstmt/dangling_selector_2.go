package danglingstmt

import "github.com/charlievieth/xtools/lsp/foo"

func _() {
	foo. //@rank(" //", Foo)
	var _ = []string{foo.} //@rank("}", Foo)
}
