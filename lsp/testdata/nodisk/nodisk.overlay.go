package nodisk

import (
	"github.com/charlievieth/xtools/lsp/foo"
)

func _() {
	foo.Foo() //@complete("F", Foo, IntFoo, StructFoo)
}
