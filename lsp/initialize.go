package lsp

type InitializeRequest struct {
	Request
	Params InitializeRequestParams `json:"params"`
}

type InitializeRequestParams struct {
	ClientInfo *ClientInfo `json:"clientInfo"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResponse struct {
	Response
	Result InitializeResult `json:"result"`
}

type InitializeResult struct {
	ServerInfo   ServerInfo         `json:"serverInfo"`
	Capabilities ServerCapabilities `json:"capabilities"`
}

type ServerCapabilities struct {
	CompletionProvider         CompletionOptions `json:"completionProvider"`
	TextDocumentSync           int               `json:"textDocumentSync"`
	HoverProvider              bool              `json:"hoverProvider"`
	DefinitionProvider         bool              `json:"definitionProvider"`
	DocumentFormattingProvider bool              `json:"documentFormattingProvider"`
}

type CompletionOptions struct {
	TriggerCharacters []string `json:"triggerCharacters"`
	ResolveProvider   bool     `json:"resolveProvider"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func NewInitializeResponse(id int) InitializeResponse {
	return InitializeResponse{
		Response: NewResponse(id),
		Result: InitializeResult{
			Capabilities: ServerCapabilities{
				// Documents are synced by always sending the full content of the document
				TextDocumentSync:   1,
				HoverProvider:      true,
				DefinitionProvider: true,
				CompletionProvider: CompletionOptions{
					TriggerCharacters: []string{".", ":", "@"},
					ResolveProvider: true,
				},
				DocumentFormattingProvider: true,
			},
			ServerInfo: ServerInfo{
				Name:    Name,
				Version: Version,
			},
		},
	}
}
