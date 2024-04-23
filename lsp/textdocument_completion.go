package lsp

type CompletionRequest struct {
	Request
	Params CompletionParams `json:"params"`
}

type CompletionParams struct {
	TextDocumentPositionParams
	Context CompletionContext `json:"context"`
}

type CompletionContext struct {
	TriggerKind      CompletionTriggerKind `json:"triggerKind"`
	TriggerCharacter string                `json:"triggerCharacter"`
}

type CompletionTriggerKind int

const (
	Invoked                         CompletionTriggerKind = 1
	TriggerCharacter                CompletionTriggerKind = 2
	TriggerForIncompleteCompletions CompletionTriggerKind = 3
)

type CompletionResponse struct {
	Response
	Result []CompletionItem `json:"result"`
}

type CompletionItem struct {
	Label         string `json:"label"`
	Detail        string `json:"detail"`
	Documentation string `json:"documentation"`
}
