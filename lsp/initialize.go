package lsp

type InitializeRequest struct {
	Request
	Params InitializeRequestParams `json:"params"`
}

type InitializeRequestParams struct {
	ClientInfo       *ClientInfo       `json:"clientInfo"`
	WorkspaceFolders []WorkspaceFolder `json:"workspaceFolders"`
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
	TextDocumentSync           int               `json:"textDocumentSync"`
	CompletionProvider         CompletionOptions `json:"completionProvider"`
	HoverProvider              bool              `json:"hoverProvider"`
	DefinitionProvider         bool              `json:"definitionProvider"`
	DocumentFormattingProvider bool              `json:"documentFormattingProvider"`
	Workspace                  Workspace         `json:"workspace"`
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
				TextDocumentSync: 1,
				CompletionProvider: CompletionOptions{
					TriggerCharacters: []string{
						"@",
						// ".",
						// ":",
					},
					ResolveProvider: true,
				},
				HoverProvider:              true,
				DefinitionProvider:         true,
				DocumentFormattingProvider: true,
				Workspace: Workspace{
					WorkspaceFolders: WorkspaceFoldersServerCapabilities{
						Supported:           true,
						ChangeNotifications: true,
					},
				},
			},
			ServerInfo: ServerInfo{
				Name:    Name,
				Version: Version,
			},
		},
	}
}
