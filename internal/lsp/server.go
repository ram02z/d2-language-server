package lsp

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	glspserver "github.com/tliron/glsp/server"
)

type Server struct {
	server *glspserver.Server
}

type ServerOpts struct {
	Name    string
	Version string
}

func NewServer(opts ServerOpts) *Server {
	handler := protocol.Handler{}
	glspServer := glspserver.NewServer(&handler, "d2", false)
	// var clientCapabilites protocol.ClientCapabilities

	handler.Initialize = func(_ *glsp.Context, params *protocol.InitializeParams) (interface{}, error) {
		// clientCapabilites = params.Capabilities

		if params.Trace != nil {
			protocol.SetTraceValue(*params.Trace)
		}

		serverCapabilities := handler.CreateServerCapabilities()
		serverCapabilities.TextDocumentSync = protocol.TextDocumentSyncKindIncremental
		serverCapabilities.CompletionProvider = &protocol.CompletionOptions{}

		return &protocol.InitializeResult{
			Capabilities: serverCapabilities,
			ServerInfo: &protocol.InitializeResultServerInfo{
				Name:    opts.Name,
				Version: &opts.Version,
			},
		}, nil
	}

	handler.Initialized = func(context *glsp.Context, params *protocol.InitializedParams) error {
		return nil
	}

	handler.Shutdown = func(_ *glsp.Context) error {
		protocol.SetTraceValue(protocol.TraceValueOff)
		return nil
	}

	handler.LogTrace = func(context *glsp.Context, params *protocol.LogTraceParams) error {
		return nil
	}

	handler.SetTrace = func(_ *glsp.Context, params *protocol.SetTraceParams) error {
		protocol.SetTraceValue(params.Value)
		return nil
	}

	handler.TextDocumentCompletion = func(context *glsp.Context, params *protocol.CompletionParams) (interface{}, error) {
		var completionItems []protocol.CompletionItem
		print(params.TextDocument.URI)

		return completionItems, nil
	}

	return &Server{
		server: glspServer,
	}
}

func (s *Server) Run() error {
	return s.server.RunStdio()
}
