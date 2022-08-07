//go:build gopls_test
// +build gopls_test

// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package debug

import (
	"testing"

	"github.com/charlievieth/xtools/gopls/hooks"
	"github.com/charlievieth/xtools/lsp/bug"
	. "github.com/charlievieth/xtools/lsp/regtest"
)

func TestMain(m *testing.M) {
	Main(m, hooks.Options)
}

func TestBugNotification(t *testing.T) {
	// Verify that a properly configured session gets notified of a bug on the
	// server.
	WithOptions(
		Modes(Default), // must be in-process to receive the bug report below
		Settings{"showBugReports": true},
	).Run(t, "", func(t *testing.T, env *Env) {
		const desc = "got a bug"
		bug.Report(desc, nil)
		env.Await(ShownMessage(desc))
	})
}