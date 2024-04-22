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
	CompletionProvider         map[string]any `json:"completionProvider"`
	TextDocumentSync           int            `json:"textDocumentSync"`
	HoverProvider              bool           `json:"hoverProvider"`
	DefinitionProvider         bool           `json:"definitionProvider"`
	DocumentFormattingProvider bool           `json:"documentFormattingProvider"`
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
				TextDocumentSync:           1,
				HoverProvider:              true,
				DefinitionProvider:         true,
				CompletionProvider:         map[string]any{},
				DocumentFormattingProvider: true,
			},
			ServerInfo: ServerInfo{
				Name:    Name,
				Version: Version,
			},
		},
	}
}
