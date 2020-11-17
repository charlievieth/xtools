package other

import "github.com/charlievieth/xtools/lsp/rename/crosspkg"

func Other() {
	crosspkg.Bar
	crosspkg.Foo() //@rename("Foo", "Flamingo")
}
