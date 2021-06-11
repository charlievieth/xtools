// +build gopls_test

// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hooks_test

import (
	"testing"

	"github.com/charlievieth/xtools/gopls/hooks"
	"github.com/charlievieth/xtools/lsp/diff/difftest"
)

func TestDiff(t *testing.T) {
	difftest.DiffTest(t, hooks.ComputeEdits)
}
