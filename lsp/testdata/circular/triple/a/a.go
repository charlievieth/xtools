package a

import (
	_ "github.com/charlievieth/xtools/lsp/circular/triple/b" //@diag("_ \"github.com/charlievieth/xtools/lsp/circular/triple/b\"", "compiler", "import cycle not allowed", "error")
)
