package main

import (
	"context"
	"runtime"
	"sync"
	"testing"

	"github.com/charlievieth/xtools/lsp/protocol"
	"github.com/charlievieth/xtools/span"
)

func TestInitialization_Race(t *testing.T) {
	const WorkingDirectory = "/Users/cvieth/go/src/github.com/charlievieth/xtools"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := Connect(&Config{})
	conn, err := client.connect(ctx, WorkingDirectory, nil)
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				if err := conn.Initialized(ctx, &protocol.InitializedParams{}); err != nil {
					t.Error(err)
					return
				}
			}
		}()
	}
	wg.Wait()

	return
}

func BenchmarkDefinition(b *testing.B) {
	const WorkingDirectory = "/Users/cvieth/go/src/github.com/charlievieth/xtools"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := Connect(&Config{})
	conn, err := client.connect(ctx, WorkingDirectory, nil)
	if err != nil {
		b.Fatal(err)
	}
	defer conn.Close()

	// line: 327
	// col:  10
	// /Users/cvieth/go/src/github.com/charlievieth/xtools/pkg/gopls/lsprpc/lsprpc.go
	const arg = "/Users/cvieth/go/src/github.com/charlievieth/xtools/pkg/gopls/lsprpc/lsprpc.go:327:10"
	from := span.Parse(arg)
	file := conn.AddFile(ctx, from.URI())
	if file.err != nil {
		b.Fatal(file.err)
	}
	loc, err := file.mapper.Location(from)
	if err != nil {
		b.Fatal(err)
	}

	// text, err := ioutil.ReadFile(file.uri.Filename())
	// if err != nil {
	// 	Fatal(err)
	// }
	// doc := protocol.TextDocumentItem{
	// 	URI:        protocol.DocumentURI(file.uri),
	// 	LanguageID: "go",
	// 	Text:       string(text),
	// 	Version:    1,
	// }
	// err = conn.DidOpen(ctx, &protocol.DidOpenTextDocumentParams{
	// 	TextDocument: doc,
	// })
	// if err != nil {
	// 	Fatal(err)
	// }

	tdpp := protocol.TextDocumentPositionParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: loc.URI},
		Position:     loc.Range.Start,
	}
	p := protocol.DefinitionParams{
		TextDocumentPositionParams: tdpp,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := conn.Definition(ctx, &p)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkInitialization(b *testing.B) {
	const WorkingDirectory = "/Users/cvieth/go/src/github.com/charlievieth/xtools"

	connect := func() error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		client := Connect(&Config{})
		conn, err := client.connect(ctx, WorkingDirectory, nil)
		if conn != nil {
			conn.Close()
		}
		return err
	}

	// warmup
	if err := connect(); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := connect(); err != nil {
			b.Fatal(err)
		}
	}

	return
}

/*
func BenchmarkInitialization_OLD(b *testing.B) {
	const WorkingDirectory = "/Users/cvieth/go/src/github.com/charlievieth/xtools"
	ctx := context.Background()

	connect := func() error {
		ropts := []lsprpc.RemoteOption{
			// lsprpc.RemoteLogfile("/Users/cvieth/go/src/github.com/charlievieth/xtools/pkg/gopls/tmp/daemon.log"),
			// lsprpc.RemoteDebugAddress(":6060"),
			// lsprpc.RemoteRPCTrace(true),
			lsprpc.RemoteListenTimeout(time.Second * 5),
		}
		conn, err := lsprpc.ConnectToRemote(ctx, lsprpc.AutoNetwork, "", ropts...)
		if err != nil {
			return err
		}
		defer conn.Close()

		connection := newConnection()

		stream := jsonrpc2.NewHeaderStream(conn)
		defer stream.Close()

		cc := jsonrpc2.NewConn(stream)
		connection.Server = protocol.ServerDispatcher(cc)

		ctx = protocol.WithClient(ctx, connection.Client)
		cc.Go(ctx,
			protocol.Handlers(
				protocol.ClientHandler(connection.Client,
					jsonrpc2.MethodNotFound)))
		return cc.Close()
	}

	// warmup
	if err := connect(); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := connect(); err != nil {
			b.Fatal(err)
		}
	}

	return
}
*/
