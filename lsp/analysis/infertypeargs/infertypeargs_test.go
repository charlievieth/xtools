// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package infertypeargs_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
	"github.com/charlievieth/xtools/lsp/analysis/infertypeargs"
	"github.com/charlievieth/xtools/testenv"
	"github.com/charlievieth/xtools/typeparams"
)

func Test(t *testing.T) {
	testenv.NeedsGo1Point(t, 13)
	if !typeparams.Enabled {
		t.Skip("type params are not enabled")
	}
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, infertypeargs.Analyzer, "a")
}
