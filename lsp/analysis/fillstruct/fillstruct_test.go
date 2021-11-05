// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fillstruct_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
	"github.com/charlievieth/xtools/lsp/analysis/fillstruct"
	"github.com/charlievieth/xtools/typeparams"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	tests := []string{"a"}
	if typeparams.Enabled {
		tests = append(tests, "typeparams")
	}
	analysistest.Run(t, testdata, fillstruct.Analyzer, tests...)
}
