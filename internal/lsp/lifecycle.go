package lsp

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var clientCapabilities *protocol.ClientCapabilities


func Initialize(_ *glsp.Context, params *protocol.InitializeParams) (interface{}, error) {
	serverVersion := Version
	clientCapabilities = &params.Capabilities

	if params.Trace != nil {
		protocol.SetTraceValue(*params.Trace)
	}

	serverCapabilities := Handler.CreateServerCapabilities()
	serverCapabilities.TextDocumentSync = protocol.TextDocumentSyncKindIncremental
	serverCapabilities.CompletionProvider = &protocol.CompletionOptions{}

	return &protocol.InitializeResult{
		Capabilities: serverCapabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    Name,
			Version: &serverVersion,
		},
	}, nil
}

func Initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func Shutdown(_ *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func LogTrace(context *glsp.Context, params *protocol.LogTraceParams) error {
	return nil
}

func SetTrace(_ *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}
