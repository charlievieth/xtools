// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package hooks adds all the standard gopls implementations.
// This can be used in tests without needing to use the gopls main, and is
// also the place to edit for custom builds of gopls.
package hooks // import "github.com/charlievieth/xtools/gopls/hooks"

import (
	"context"

	"github.com/charlievieth/xtools/gopls/vulncheck"
	"github.com/charlievieth/xtools/lsp/diff"
	"github.com/charlievieth/xtools/lsp/source"
	"mvdan.cc/gofumpt/format"
	"mvdan.cc/xurls/v2"
)

func Options(options *source.Options) {
	options.LicensesText = licensesText
	if options.GoDiff {
		switch options.NewDiff {
		case "old":
			options.ComputeEdits = ComputeEdits
		case "new":
			options.ComputeEdits = diff.NComputeEdits
		default:
			options.ComputeEdits = BothDiffs
		}
	}
	options.URLRegexp = xurls.Relaxed()
	options.GofumptFormat = func(ctx context.Context, langVersion, modulePath string, src []byte) ([]byte, error) {
		return format.Source(src, format.Options{
			LangVersion: langVersion,
			ModulePath:  modulePath,
		})
	}
	updateAnalyzers(options)

	options.Govulncheck = vulncheck.Govulncheck
}
