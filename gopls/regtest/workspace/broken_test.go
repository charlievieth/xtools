//go:build gopls_test
// +build gopls_test

// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package workspace

import (
	"strings"
	"testing"

	"github.com/charlievieth/xtools/lsp"
	. "github.com/charlievieth/xtools/lsp/regtest"
	"github.com/charlievieth/xtools/testenv"
)

// This file holds various tests for UX with respect to broken workspaces.
//
// TODO: consolidate other tests here.

// Test for golang/go#53933
func TestBrokenWorkspace_DuplicateModules(t *testing.T) {
	testenv.NeedsGo1Point(t, 18)

	// This proxy module content is replaced by the workspace, but is still
	// required for module resolution to function in the Go command.
	const proxy = `
-- example.com/foo@v0.0.1/go.mod --
module example.com/foo

go 1.12
-- example.com/foo@v1.2.3/foo.go --
package foo
`

	const src = `
-- go.work --
go 1.18

use (
	./package1
	./package1/vendor/example.com/foo
	./package2
	./package2/vendor/example.com/foo
)

-- package1/go.mod --
module mod.test

go 1.18

require example.com/foo v0.0.1
-- package1/main.go --
package main

import "example.com/foo"

func main() {
	_ = foo.CompleteMe
}
-- package1/vendor/example.com/foo/go.mod --
module example.com/foo

go 1.18
-- package1/vendor/example.com/foo/foo.go --
package foo

const CompleteMe = 111
-- package2/go.mod --
module mod2.test

go 1.18

require example.com/foo v0.0.1
-- package2/main.go --
package main

import "example.com/foo"

func main() {
	_ = foo.CompleteMe
}
-- package2/vendor/example.com/foo/go.mod --
module example.com/foo

go 1.18
-- package2/vendor/example.com/foo/foo.go --
package foo

const CompleteMe = 222
`

	WithOptions(
		ProxyFiles(proxy),
	).Run(t, src, func(t *testing.T, env *Env) {
		env.OpenFile("package1/main.go")
		env.Await(
			OutstandingWork(lsp.WorkspaceLoadFailure, `found module "example.com/foo" multiple times in the workspace`),
		)

		// Remove the redundant vendored copy of example.com.
		env.WriteWorkspaceFile("go.work", `go 1.18
		use (
			./package1
			./package2
			./package2/vendor/example.com/foo
		)
		`)
		env.Await(NoOutstandingWork())

		// Check that definitions in package1 go to the copy vendored in package2.
		location, _ := env.GoToDefinition("package1/main.go", env.RegexpSearch("package1/main.go", "CompleteMe"))
		const wantLocation = "package2/vendor/example.com/foo/foo.go"
		if !strings.HasSuffix(location, wantLocation) {
			t.Errorf("got definition of CompleteMe at %q, want %q", location, wantLocation)
		}
	})
}

// Test for golang/go#43186: correcting the module path should fix errors
// without restarting gopls.
func TestBrokenWorkspace_WrongModulePath(t *testing.T) {
	const files = `
-- go.mod --
module mod.testx

go 1.18
-- p/internal/foo/foo.go --
package foo

const C = 1
-- p/internal/bar/bar.go --
package bar

import "mod.test/p/internal/foo"

const D = foo.C + 1
-- p/internal/bar/bar_test.go --
package bar_test

import (
	"mod.test/p/internal/foo"
	. "mod.test/p/internal/bar"
)

const E = D + foo.C
-- p/internal/baz/baz_test.go --
package baz_test

import (
	named "mod.test/p/internal/bar"
)

const F = named.D - 3
`

	Run(t, files, func(t *testing.T, env *Env) {
		env.OpenFile("p/internal/bar/bar.go")
		env.Await(
			OnceMet(
				env.DoneWithOpen(),
				env.DiagnosticAtRegexp("p/internal/bar/bar.go", "\"mod.test/p/internal/foo\""),
			),
		)
		env.OpenFile("go.mod")
		env.RegexpReplace("go.mod", "mod.testx", "mod.test")
		env.SaveBuffer("go.mod") // saving triggers a reload
		env.Await(NoOutstandingDiagnostics())
	})
}
