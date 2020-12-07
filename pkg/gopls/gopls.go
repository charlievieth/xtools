package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/charlievieth/xtools/lsp/cmd"
	"github.com/charlievieth/xtools/lsp/lsprpc"
	"github.com/charlievieth/xtools/lsp/protocol"
	"github.com/charlievieth/xtools/span"
)

const DefaultRemoteListenTimeout = time.Minute * 10

// TODO: consider using cmd/Serve.Run
type Config struct {
	Remote string

	// RemoteListenTimeout configures the amount of time the auto-started gopls
	// daemon will wait with no client connections before shutting down.
	RemoteListenTimeout time.Duration

	// TODO: expose this
	//
	// RemoteLogfile configures the logfile location for the auto-started gopls
	// daemon.
	// RemoteLogfile string

	// GoplsPath allows the gopls binary used for the backend server to be
	// overridden, otherwise the first gopls found on the PATH will be used.
	GoplsPath string
}

func (c *Config) init() (err error) {
	if c.Remote == "" {
		c.Remote = lsprpc.AutoNetwork
	}
	switch {
	case c.RemoteListenTimeout == 0:
		c.RemoteListenTimeout = DefaultRemoteListenTimeout
	case c.RemoteListenTimeout < 0:
		c.RemoteListenTimeout = 0 // disable the timeout
	}
	if c.GoplsPath == "" {
		c.GoplsPath, err = exec.LookPath("gopls")
		if err != nil {
			return fmt.Errorf(`missing "gopls" executable: %w`, err)
		}
	}
	return nil
}

// parseAddr parses the -listen flag in to a network, and address.
// func parseAddr(listen string) (network string, address string) {
// 	// Allow passing just -remote=auto, as a shorthand for using automatic remote
// 	// resolution.
// 	if listen == lsprpc.AutoNetwork {
// 		return lsprpc.AutoNetwork, ""
// 	}
// 	if parts := strings.SplitN(listen, ";", 2); len(parts) == 2 {
// 		return parts[0], parts[1]
// 	}
// 	return "tcp", listen
// }

// TODO: consider using cmd/Serve.Run
func (c *Config) Connect(ctx context.Context) (*Client, error) {
	// TODO: consider not using a method pointer
	if err := c.init(); err != nil {
		return nil, err
	}
	lsprpc.ConnectToRemote(ctx, "network", "addr")
	return nil, nil
}

type Server struct {
	Address string

	// ListenTimeout configures the amount of time the auto-started gopls
	// daemon will wait with no client connections before shutting down.
	ListenTimeout time.Duration

	// TODO: remove the socket when we create it by default.
	//
	// Remove the underlying socket when the server stops.
	UnlinkOnClose bool

	// TODO: expose this
	//
	// RemoteLogfile configures the logfile location for the auto-started gopls
	// daemon.
	// RemoteLogfile string

	// GoplsPath allows the gopls binary used for the backend server to be
	// overridden, otherwise the first gopls found on the PATH will be used.
	GoplsPath string

	protocol.Server            // WARN
	srv             *cmd.Serve // TODO: use or remove
}

func (s *Server) Run(ctx context.Context) error {

	return nil
}

// CEV: we'll need something like this to store request scoped data
type request struct {
	fset *token.FileSet

	// filesMu sync.Mutex
	// files   map[span.URI]*cmdFile
}

type Client struct {
	server protocol.Server // WARN
}

func (c *Client) References(ctx context.Context, location *Location, src interface{}) ([]protocol.Location, error) {
	var x protocol.Location
	_ = x
	return nil, nil
}

// func (c *Client) References(ctx context.Context, params *protocol.ReferenceParams) ([]protocol.Location, error) {
// 	var x protocol.Location
// 	_ = x
// 	return nil, nil
// }

func main() {
	var _ = span.URI("")
	const filename = "/Users/cvieth/go/src/github.com/charlievieth/xtools/pkg/gopls/gopls.go:#1023"
	// u := span.URIFromPath(filename)
	// fmt.Println(u.Filename())

	loc := Location{
		Path: "/Users/cvieth/go/src/github.com/charlievieth/xtools/pkg/gopls/gopls.go",
		Point: Point{
			Offset: 1023,
			Line:   2,
			Column: 4,
		},
	}
	p := loc.Span()
	fmt.Printf("%s\n", p)
	fmt.Println(loc)

	// loc := Location{
	// 	Path:  "/Users/cvieth/go/src/github.com/charlievieth/xtools/pkg/gopls/gopls.go",
	// 	Point: Point{Offset: 1023},
	// }
	// p1 := span.Parse(filename)
	// p2 := loc.Span()

	// PrintJSON(&p1)
	// PrintJSON(&p2)
	// fmt.Println(span.Compare(p1, p2))
	return

	// fmt.Printf("%+v\n", p)
	// fmt.Println(p.URI().Filename(), p.URI().IsFile())
	// fmt.Println("HasOffset:", p.HasOffset())
	// fmt.Println("HasPosition:", p.HasPosition())
}

func PrintJSON(v interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		Fatal(err)
	}
}

func Fatal(err interface{}) {
	if err == nil {
		return
	}
	var s string
	if _, file, line, ok := runtime.Caller(1); ok && file != "" {
		s = fmt.Sprintf("Error (%s:%d)", filepath.Base(file), line)
	} else {
		s = "Error"
	}
	switch err.(type) {
	case error, string, fmt.Stringer:
		fmt.Fprintf(os.Stderr, "%s: %s\n", s, err)
	default:
		fmt.Fprintf(os.Stderr, "%s: %#v\n", s, err)
	}
	os.Exit(1)
}
