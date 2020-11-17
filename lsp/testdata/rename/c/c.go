package c

import "github.com/charlievieth/xtools/lsp/rename/b"

func _() {
	b.Hello() //@rename("Hello", "Goodbye")
}
