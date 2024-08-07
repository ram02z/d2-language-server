package lsp

type Method string

const (
	Initialize                Method = "initialize"
	DidOpenTextDocument       Method = "textDocument/didOpen"
	DidChangeTextDocument     Method = "textDocument/didChange"
	DidCloseTextDocument      Method = "textDocument/didClose"
	PublishDiagnostics        Method = "textDocument/publishDiagnostics"
	Hover                     Method = "textDocument/hover"
	Definition                Method = "textDocument/definition"
	Completion                Method = "textDocument/completion"
	Formatting                Method = "textDocument/formatting"
	DidChangeWorkspaceFolders Method = "workspace/didChangeWorkspaceFolders"
)
