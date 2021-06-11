package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/charlievieth/xtools/jsonrpc2"
	"github.com/charlievieth/xtools/lsp"
	"github.com/charlievieth/xtools/lsp/cmd"
	"github.com/charlievieth/xtools/lsp/protocol"
	"github.com/charlievieth/xtools/lsp/source"
	"github.com/charlievieth/xtools/pkg/gopls/lsprpc"
	"github.com/charlievieth/xtools/span"
	errors "golang.org/x/xerrors"
)

// TODO: consider using cmd/Serve.Run
type Config struct {
	Remote string

	// GoplsPath allows the gopls binary used for the backend server to be
	// overridden, otherwise the first gopls found on the PATH will be used.
	GoplsPath string

	// RemoteListenTimeout configures the amount of time the auto-started gopls
	// daemon will wait with no client connections before shutting down.
	RemoteListenTimeout time.Duration

	// RemoteDebug serve debug information on the supplied address.
	RemoteDebug string

	// TODO: expose this
	//
	// RemoteLogfile configures the logfile location for the auto-started gopls
	// daemon.
	RemoteLogfile string

	// RemoteRPCTrace configures the rpc.trace option for the remote daemon.
	RemoteRPCTrace bool

	// TODO: add RemoteVerbose option
}

func (conf *Config) init() (err error) {
	if conf.Remote == "" {
		conf.Remote = lsprpc.AutoNetwork
	}
	switch {
	case conf.RemoteListenTimeout == 0:
		conf.RemoteListenTimeout = lsprpc.DefaultRemoteListenTimeout
	case conf.RemoteListenTimeout < 0:
		conf.RemoteListenTimeout = 0 // disable the timeout
	}
	if conf.GoplsPath == "" {
		conf.GoplsPath, err = exec.LookPath("gopls")
		if err != nil {
			return fmt.Errorf(`missing "gopls" executable: %w`, err)
		}
	}
	if conf.RemoteLogfile != "" {
		dir := filepath.Dir(conf.RemoteLogfile)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

func Connect(conf *Config) *Client {
	resolved := *conf
	if resolved.Remote == "" {
		resolved.Remote = lsprpc.AutoNetwork
	}
	// WARN (CEV): make sure this is correct (with our defaults)
	//
	// Use default timeout if not set
	if resolved.RemoteListenTimeout == 0 {
		// TODO (CEV): make this longer by default
		resolved.RemoteListenTimeout = lsprpc.DefaultRemoteListenTimeout
	}
	// Listen indefinitely if less than zero.
	if resolved.RemoteListenTimeout < 0 {
		resolved.RemoteListenTimeout = 0
	}
	if resolved.GoplsPath == "" {
		resolved.GoplsPath, _ = exec.LookPath("gopls")
	}
	return &Client{config: resolved}
}

// TODO (CEV): use per-workspace gopls servers
type Client struct {
	config Config
}

func (c *Client) References(ctx context.Context, location *Location, src interface{}) ([]protocol.Location, error) {
	var x protocol.Location
	_ = x
	return nil, nil
}

// TODO: consider using cmd/Serve.Run
func (c *Config) Connect(ctx context.Context) (*Client, error) {
	// TODO: consider not using a method pointer
	if err := c.init(); err != nil {
		return nil, err
	}
	lsprpc.ConnectToRemote(ctx, "network", "addr")
	return nil, nil
}

// CEV: we'll need something like this to store request scoped data
type request struct {
	fset *token.FileSet

	// filesMu sync.Mutex
	// files   map[span.URI]*cmdFile
}

// func (c *Client) References(ctx context.Context, params *protocol.ReferenceParams) ([]protocol.Location, error) {
// 	var x protocol.Location
// 	_ = x
// 	return nil, nil
// }

/*
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := exec.CommandContext(ctx, "gopls", "-remote=auto")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		Fatal(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		Fatal(err)
	}

	go func() {
		r := bufio.NewReader(stderr)
		for {
			b, e := r.ReadBytes('\n')
			if len(b) != 0 {
				fmt.Fprintf(os.Stderr, "ERROR: %s\n", b)
			}
			if e != nil {
				if e != io.EOF {
					fmt.Fprintf(os.Stderr, "ERROR: reading STDERR: %s\n", err)
				}
				break
			}
		}
	}()
}
*/

func main() {
	// {
	// 	sp := span.Parse("/usr/local/Cellar/go/1.15.5/libexec/src/go/build/build.go:46:2")
	// 	fmt.Println(sp.Start())
	// 	fmt.Println(sp.Start().Column())
	// 	fmt.Println(sp.Start().Line())
	// 	return
	// }
	const WorkingDirectory = "/Users/cvieth/go/src/github.com/charlievieth/xtools"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := Connect(&Config{})
	conn, err := client.connect(ctx, WorkingDirectory, nil)
	if err != nil {
		Fatal(err)
	}
	defer conn.Close()

	// line: 327
	// col:  10
	// /Users/cvieth/go/src/github.com/charlievieth/xtools/pkg/gopls/lsprpc/lsprpc.go
	const arg = "/Users/cvieth/go/src/github.com/charlievieth/xtools/pkg/gopls/lsprpc/lsprpc.go:327:10"
	from := span.Parse(arg)

	// WARN: this opens the file for us
	file := conn.AddFile(ctx, from.URI())
	if file.err != nil {
		Fatal(file.err)
	}

	loc, err := file.mapper.Location(from)
	if err != nil {
		Fatal(err)
	}

	// CEV: this is done by AddFile
	//
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
	locs, err := conn.Definition(ctx, &p)
	if err != nil {
		Fatal(err)
	}
	PrintJSON(locs)

	// server state
	{
		state, err := lsprpc.QueryServerState(ctx, lsprpc.AutoNetwork, "")
		if err != nil {
			Fatal(err)
		}
		fmt.Println("### Sessions:")
		PrintJSON(state)
		fmt.Println("###")
	}

	closeReq := protocol.DidCloseTextDocumentParams{
		TextDocument: protocol.TextDocumentIdentifier{
			URI: protocol.DocumentURI(file.uri),
		},
	}
	if err := conn.DidClose(ctx, &closeReq); err != nil {
		Fatal(err)
	}
}

func (c *Client) remoteOptions() []lsprpc.RemoteOption {
	var opts []lsprpc.RemoteOption
	// WARN (CEV): make sure this is correct (with our defaults)
	if c.config.RemoteListenTimeout > 0 {
		opts = append(opts, lsprpc.RemoteListenTimeout(c.config.RemoteListenTimeout))
	} else if c.config.RemoteListenTimeout < 0 {
		// WARN: use indefinite timeout for <0
		opts = append(opts, lsprpc.RemoteListenTimeout(0))
	}
	if c.config.RemoteLogfile != "" {
		opts = append(opts, lsprpc.RemoteLogfile(c.config.RemoteLogfile))
	}
	if c.config.RemoteDebug != "" {
		opts = append(opts, lsprpc.RemoteDebugAddress(c.config.RemoteDebug))
	}
	if c.config.RemoteRPCTrace {
		opts = append(opts, lsprpc.RemoteRPCTrace(c.config.RemoteRPCTrace))
	}
	return opts
}

// parseAddr parses the -listen flag in to a network, and address.
func parseAddr(listen string) (network string, address string) {
	// Allow passing just -remote=auto, as a shorthand for using automatic remote
	// resolution.
	if listen == lsprpc.AutoNetwork {
		return lsprpc.AutoNetwork, ""
	}
	if parts := strings.SplitN(listen, ";", 2); len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "tcp", listen
}

func (c *Client) connect(ctx context.Context, wd string, options func(*source.Options)) (*connection, error) {

	if c.config.RemoteLogfile != "" {
		dir := filepath.Dir(c.config.RemoteLogfile)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("creating log file directory: %w", err)
		}
	}

	network, addr := parseAddr(c.config.Remote)
	conn, err := lsprpc.ConnectToRemote(ctx, network, addr, c.remoteOptions()...)
	if err != nil {
		return nil, err
	}

	connection := newConnection()

	stream := jsonrpc2.NewHeaderStream(conn)
	connection.conn = jsonrpc2.NewConn(stream)
	connection.Server = protocol.ServerDispatcher(connection.conn)

	ctx = protocol.WithClient(ctx, connection.Client)
	connection.conn.Go(ctx, protocol.Handlers(
		protocol.ClientHandler(
			connection.Client,
			jsonrpc2.MethodNotFound,
		),
	))

	return connection, connection.initialize(ctx, wd, options)
}

func (c *connection) Close() error {
	return c.conn.Close()
}

type connection struct {
	protocol.Server
	conn   jsonrpc2.Conn
	Client *cmdClient
}

func newConnection() *connection {
	return &connection{
		Client: &cmdClient{
			fset:  token.NewFileSet(),
			files: make(map[span.URI]*cmdFile),
		},
	}
}

type cmdClient struct {
	protocol.Server

	// WARN (CEV): fix me
	// app  *Application

	fset *token.FileSet

	logMessage func(ctx context.Context, p *protocol.LogMessageParams) error

	diagnosticsMu   sync.Mutex
	diagnosticsDone chan struct{}

	filesMu sync.Mutex
	files   map[span.URI]*cmdFile
}

type cmdFile struct {
	uri         span.URI
	mapper      *protocol.ColumnMapper
	err         error
	added       bool
	diagnostics []protocol.Diagnostic
}

// fileURI converts a DocumentURI to a file:// span.URI, panicking if it's not a file.
func fileURI(uri protocol.DocumentURI) span.URI {
	sURI := uri.SpanURI()
	if !sURI.IsFile() {
		panic(fmt.Sprintf("%q is not a file URI", uri))
	}
	return sURI
}

var matcherString = map[source.SymbolMatcher]string{
	source.SymbolFuzzy:           "fuzzy",
	source.SymbolCaseSensitive:   "caseSensitive",
	source.SymbolCaseInsensitive: "caseInsensitive",
}

func (c *connection) initialize(ctx context.Context, wd string, options func(*source.Options)) error {
	params := &protocol.ParamInitialize{}
	// WARN (CEV): make this more configurable
	// params.RootURI = protocol.URIFromPath(c.Client.app.wd)
	params.RootURI = protocol.URIFromPath(wd)
	params.Capabilities.Workspace.Configuration = true

	// Make sure to respect configured options when sending initialize request.
	opts := source.DefaultOptions().Clone()
	if options != nil {
		options(opts)
	}
	params.Capabilities.TextDocument.Hover = protocol.HoverClientCapabilities{
		ContentFormat: []protocol.MarkupKind{opts.PreferredContentFormat},
	}
	params.Capabilities.TextDocument.DocumentSymbol.HierarchicalDocumentSymbolSupport = opts.HierarchicalDocumentSymbolSupport
	params.Capabilities.TextDocument.SemanticTokens = protocol.SemanticTokensClientCapabilities{}
	params.Capabilities.TextDocument.SemanticTokens.Formats = []string{"relative"}
	params.Capabilities.TextDocument.SemanticTokens.Requests.Range = true
	params.Capabilities.TextDocument.SemanticTokens.Requests.Full = true
	params.Capabilities.TextDocument.SemanticTokens.TokenTypes = lsp.SemanticTypes()
	params.Capabilities.TextDocument.SemanticTokens.TokenModifiers = lsp.SemanticModifiers()
	params.InitializationOptions = map[string]interface{}{
		"symbolMatcher": matcherString[opts.SymbolMatcher],
	}
	if _, err := c.Server.Initialize(ctx, params); err != nil {
		return err
	}
	if err := c.Server.Initialized(ctx, &protocol.InitializedParams{}); err != nil {
		return err
	}
	return nil
}

func (c *cmdClient) ShowMessage(ctx context.Context, p *protocol.ShowMessageParams) error { return nil }

func (c *cmdClient) ShowMessageRequest(ctx context.Context, p *protocol.ShowMessageRequestParams) (*protocol.MessageActionItem, error) {
	return nil, nil
}

// WARN (CEV): fix the logging logic
func (c *cmdClient) LogMessage(ctx context.Context, p *protocol.LogMessageParams) error {
	if c.logMessage != nil {
		return c.logMessage(ctx, p)
	}
	switch p.Type {
	case protocol.Error:
		log.Print("Error:", p.Message)
	case protocol.Warning:
		log.Print("Warning:", p.Message)
	case protocol.Info:
		// if c.app.verbose() {
		log.Print("Info:", p.Message)
		// }
	case protocol.Log:
		// if c.app.verbose() {
		log.Print("Log:", p.Message)
		// }
	default:
		// if c.app.verbose() {
		log.Printf("%v: %v", p.Type, p.Message)
		// }
	}
	return nil
}

func (c *cmdClient) Event(ctx context.Context, t *interface{}) error { return nil }

func (c *cmdClient) RegisterCapability(ctx context.Context, p *protocol.RegistrationParams) error {
	return nil
}

func (c *cmdClient) UnregisterCapability(ctx context.Context, p *protocol.UnregistrationParams) error {
	return nil
}

func (c *cmdClient) WorkspaceFolders(ctx context.Context) ([]protocol.WorkspaceFolder, error) {
	return nil, nil
}

func (c *cmdClient) Configuration(ctx context.Context, p *protocol.ParamConfiguration) ([]interface{}, error) {
	results := make([]interface{}, len(p.Items))
	for i, item := range p.Items {
		if item.Section != "gopls" {
			continue
		}
		env := map[string]interface{}{}
		// WARN (CEV): fix me
		//
		// for _, value := range c.app.env {
		// 	l := strings.SplitN(value, "=", 2)
		// 	if len(l) != 2 {
		// 		continue
		// 	}
		// 	env[l[0]] = l[1]
		// }

		// Docs: xtools/lsp/source/options.go
		m := map[string]interface{}{
			"env": env,
			"analyses": map[string]bool{
				"fillreturns":    true,
				"nonewvars":      true,
				"noresultvalues": true,
				"undeclaredname": true,
			},
		}
		// WARN (CEV): fix me
		// if c.app.VeryVerbose {
		m["verboseOutput"] = true
		// }
		results[i] = m
	}
	return results, nil
}

func (c *cmdClient) ApplyEdit(ctx context.Context, p *protocol.ApplyWorkspaceEditParams) (*protocol.ApplyWorkspaceEditResponse, error) {
	return &protocol.ApplyWorkspaceEditResponse{Applied: false, FailureReason: "not implemented"}, nil
}

func (c *cmdClient) PublishDiagnostics(ctx context.Context, p *protocol.PublishDiagnosticsParams) error {
	if p.URI == "gopls://diagnostics-done" {
		close(c.diagnosticsDone)
	}
	// Don't worry about diagnostics without versions.
	if p.Version == 0 {
		return nil
	}

	c.filesMu.Lock()
	defer c.filesMu.Unlock()

	file := c.getFile(ctx, fileURI(p.URI))
	file.diagnostics = p.Diagnostics
	return nil
}

func (c *cmdClient) Progress(context.Context, *protocol.ProgressParams) error {
	return nil
}

func (c *cmdClient) ShowDocument(context.Context, *protocol.ShowDocumentParams) (*protocol.ShowDocumentResult, error) {
	return nil, nil
}

func (c *cmdClient) WorkDoneProgressCreate(context.Context, *protocol.WorkDoneProgressCreateParams) error {
	return nil
}

func (c *cmdClient) getFile(ctx context.Context, uri span.URI) *cmdFile {
	file, found := c.files[uri]
	if !found || file.err != nil {
		file = &cmdFile{
			uri: uri,
		}
		c.files[uri] = file
	}
	if file.mapper == nil {
		fname := uri.Filename()
		content, err := ioutil.ReadFile(fname)
		if err != nil {
			file.err = errors.Errorf("getFile: %v: %v", uri, err)
			return file
		}
		f := c.fset.AddFile(fname, -1, len(content))
		f.SetLinesForContent(content)
		converter := span.NewContentConverter(fname, content)
		file.mapper = &protocol.ColumnMapper{
			URI:       uri,
			Converter: converter,
			Content:   content,
		}
	}
	return file
}

func (c *cmdClient) getFile_XXX(ctx context.Context, uri span.URI, modified string) *cmdFile {
	file, found := c.files[uri]
	if !found || file.err != nil {
		file = &cmdFile{
			uri: uri,
		}
		c.files[uri] = file
	}
	if file.mapper == nil {
		fname := uri.Filename()
		var content []byte
		var err error
		if modified != "" {
			content = []byte(modified)
		} else {
			content, err = ioutil.ReadFile(fname)
			if err != nil {
				file.err = errors.Errorf("getFile: %v: %v", uri, err)
				return file
			}
		}
		f := c.fset.AddFile(fname, -1, len(content))
		f.SetLinesForContent(content)
		converter := span.NewContentConverter(fname, content)
		file.mapper = &protocol.ColumnMapper{
			URI:       uri,
			Converter: converter,
			Content:   content,
		}
	}
	return file
}

// WARN: WIP
func (c *connection) AddModifiedFile(ctx context.Context, uri span.URI, content string) *cmdFile {
	c.Client.filesMu.Lock()
	defer c.Client.filesMu.Unlock()

	return nil
}

// TODO (CEV): support in-memory files
func (c *connection) AddFile(ctx context.Context, uri span.URI) *cmdFile {
	c.Client.filesMu.Lock()
	defer c.Client.filesMu.Unlock()

	file := c.Client.getFile(ctx, uri)
	// This should never happen.
	if file == nil {
		return &cmdFile{
			uri: uri,
			err: fmt.Errorf("no file found for %s", uri),
		}
	}
	if file.err != nil || file.added {
		return file
	}
	file.added = true
	p := &protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        protocol.URIFromSpanURI(uri),
			LanguageID: source.DetectLanguage("", file.uri.Filename()).String(),
			Version:    1,
			Text:       string(file.mapper.Content),
		},
	}
	if err := c.Server.DidOpen(ctx, p); err != nil {
		file.err = errors.Errorf("%v: %v", uri, err)
	}
	return file
}

func (c *connection) semanticTokens(ctx context.Context, file span.URI) (*protocol.SemanticTokens, error) {
	p := &protocol.SemanticTokensParams{
		TextDocument: protocol.TextDocumentIdentifier{
			URI: protocol.URIFromSpanURI(file),
		},
	}
	resp, err := c.Server.SemanticTokensFull(ctx, p)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *connection) diagnoseFiles(ctx context.Context, files []span.URI) error {
	var untypedFiles []interface{}
	for _, file := range files {
		untypedFiles = append(untypedFiles, string(file))
	}
	c.Client.diagnosticsMu.Lock()
	defer c.Client.diagnosticsMu.Unlock()

	c.Client.diagnosticsDone = make(chan struct{})
	_, err := c.Server.NonstandardRequest(ctx, "gopls/diagnoseFiles", map[string]interface{}{"files": untypedFiles})
	<-c.Client.diagnosticsDone
	return err
}

func (c *connection) terminate(ctx context.Context) {
	// WARN (CEV): fix me
	//
	// if strings.HasPrefix(c.Client.app.Remote, "internal@") {
	// 	// internal connections need to be left alive for the next test
	// 	return
	// }

	//TODO: do we need to handle errors on these calls?
	c.Shutdown(ctx)
	//TODO: right now calling exit terminates the process, we should rethink that
	//server.Exit(ctx)
}

// Implement io.Closer.
func (c *cmdClient) Close() error {
	return nil
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

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
