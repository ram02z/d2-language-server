package lsp

import protocol "github.com/tliron/glsp/protocol_3_16"

var Handler protocol.Handler

func init() {
	// Lifecycle messages
	Handler.Initialize = Initialize
	Handler.Initialized = Initialized
	Handler.Shutdown = Shutdown
	Handler.LogTrace = LogTrace
	Handler.SetTrace = SetTrace
}
