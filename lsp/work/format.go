// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package work

import (
	"context"

	"golang.org/x/mod/modfile"
	"github.com/charlievieth/xtools/event"
	"github.com/charlievieth/xtools/lsp/protocol"
	"github.com/charlievieth/xtools/lsp/source"
)

func Format(ctx context.Context, snapshot source.Snapshot, fh source.FileHandle) ([]protocol.TextEdit, error) {
	ctx, done := event.Start(ctx, "work.Format")
	defer done()

	pw, err := snapshot.ParseWork(ctx, fh)
	if err != nil {
		return nil, err
	}
	formatted := modfile.Format(pw.File.Syntax)
	// Calculate the edits to be made due to the change.
	diff, err := snapshot.View().Options().ComputeEdits(fh.URI(), string(pw.Mapper.Content), string(formatted))
	if err != nil {
		return nil, err
	}
	return source.ToProtocolEdits(pw.Mapper, diff)
}
