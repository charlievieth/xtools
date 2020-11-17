package errors

import (
	"github.com/charlievieth/xtools/lsp/types"
)

func _() {
	bob.Bob() //@complete(".")
	types.b //@complete(" //", Bob_interface)
}
