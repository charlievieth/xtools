// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd_test

import (
	"os"
	"testing"

	cmdtest "github.com/charlievieth/xtools/lsp/cmd/test"
	"github.com/charlievieth/xtools/lsp/tests"
	"github.com/charlievieth/xtools/testenv"
)

func TestMain(m *testing.M) {
	testenv.ExitIfSmallMachine()
	os.Exit(m.Run())
}

func TestCommandLine(t *testing.T) {
	cmdtest.TestCommandLine(t, "../testdata", tests.DefaultOptions)
}
