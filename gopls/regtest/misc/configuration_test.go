//go:build gopls_test
// +build gopls_test

// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package misc

import (
	"testing"

	. "github.com/charlievieth/xtools/lsp/regtest"

	"github.com/charlievieth/xtools/lsp/fake"
	"github.com/charlievieth/xtools/testenv"
)

// Test that enabling and disabling produces the expected results of showing
// and hiding staticcheck analysis results.
func TestChangeConfiguration(t *testing.T) {
	// Staticcheck only supports Go versions > 1.14.
	testenv.NeedsGo1Point(t, 15)

	const files = `
-- go.mod --
module mod.com

go 1.12
-- a/a.go --
package a

import "errors"

// FooErr should be called ErrFoo (ST1012)
var FooErr = errors.New("foo")
`
	Run(t, files, func(t *testing.T, env *Env) {
		env.OpenFile("a/a.go")
		env.Await(
			env.DoneWithOpen(),
			NoDiagnostics("a/a.go"),
		)
		cfg := &fake.EditorConfig{}
		*cfg = env.Editor.Config
		cfg.EnableStaticcheck = true
		env.ChangeConfiguration(t, cfg)
		env.Await(
			DiagnosticAt("a/a.go", 5, 4),
		)
	})
}
