//go:build gopls_test
// +build gopls_test

// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bench

import (
	"context"
	"testing"

	"github.com/charlievieth/xtools/lsp/fake"
	. "github.com/charlievieth/xtools/lsp/regtest"
)

// BenchmarkIWL benchmarks the initial workspace load time for a new editing
// session.
func BenchmarkIWL(b *testing.B) {
	dir := benchmarkDir()
	b.ResetTimer()

	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		_, editor, awaiter, err := connectEditor(dir, fake.EditorConfig{})
		if err != nil {
			b.Fatal(err)
		}
		if err := awaiter.Await(ctx, InitialWorkspaceLoad); err != nil {
			b.Fatal(err)
		}
		b.StopTimer()
		if err := editor.Close(ctx); err != nil {
			b.Fatal(err)
		}
		b.StartTimer()
	}
}
