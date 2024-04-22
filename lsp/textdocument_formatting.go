package lsp

type FormattingRequest struct {
	Request
	Params FormattingParams `json:"params"`
}

type FormattingParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type FormattingResponse struct {
	Response
	Result []TextEdit `json:"result"`
}
