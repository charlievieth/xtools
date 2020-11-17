// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fillstruct_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
	"github.com/charlievieth/xtools/lsp/analysis/fillstruct"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, fillstruct.Analyzer, "a")
}