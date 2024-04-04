package main

import (
	"bufio"
	"encoding/json"
	"io"
	"os"

	"github.com/ram02z/d2-language-server/analysis"
	"github.com/ram02z/d2-language-server/log"
	"github.com/ram02z/d2-language-server/lsp"
	"github.com/ram02z/d2-language-server/rpc"
)

type HandlerFunc func(*log.Logger, io.Writer, analysis.State, []byte)

var handlers = map[lsp.Method]HandlerFunc{
	lsp.Initialize:            handleInitialize,
	lsp.DidOpenTextDocument:   handleDidOpenTextDocument,
	lsp.DidChangeTextDocument: handleDidChangeTextDocument,
	lsp.DidSaveTextDocument:   handleDidSaveTextDocument,
	lsp.Hover:                 handleHover,
	lsp.Definition:            handleDefinition,
	lsp.Completion:            handleCompletion,
}

func handleMessage(
	logger *log.Logger,
	writer io.Writer,
	state analysis.State,
	method string,
	contents []byte,
) {
	handler, ok := handlers[lsp.Method(method)]
	if !ok {
		logger.Printf("Unsupported method: %s", method)
		return
	}

	logger.Printf("Received message with method: %s", method)
	handler(logger, writer, state, contents)
}

func handleInitialize(logger *log.Logger, writer io.Writer, state analysis.State, contents []byte) {
	var request lsp.InitializeRequest
	if err := json.Unmarshal(contents, &request); err != nil {
		logger.Printf("error parsing %s request: %s", lsp.Initialize, err)
		return
	}

	logger.Printf(
		"Connected to: %s %s",
		request.Params.ClientInfo.Name,
		request.Params.ClientInfo.Version,
	)

	msg := lsp.NewInitializeResponse(request.ID)
	writeResponse(writer, msg)
}

func handleDidOpenTextDocument(logger *log.Logger, writer io.Writer, state analysis.State, contents []byte) {
	var request lsp.DidOpenTextDocumentNotification
	if err := json.Unmarshal(contents, &request); err != nil {
		logger.Printf("error parsing %s request: %s", lsp.DidOpenTextDocument, err)
		return
	}

	logger.Printf("Opened: %s", request.Params.TextDocument.URI)
	diagnostics := state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)
	writeResponse(writer, lsp.PublishDiagnosticsNotification{
		Notification: lsp.NewNotification(lsp.PublishDiagnostics),
		Params: lsp.PublishDiagnosticsParams{
			URI:         request.Params.TextDocument.URI,
			Diagnostics: diagnostics,
		},
	})
	logger.Printf("Published %d diagnostics", len(diagnostics))
}

func handleDidChangeTextDocument(logger *log.Logger, writer io.Writer, state analysis.State, contents []byte) {
	var request lsp.DidChangeTextDocumentNotification
	if err := json.Unmarshal(contents, &request); err != nil {
		logger.Printf("error parsing %s request: %s", lsp.DidChangeTextDocument, err)
		return
	}

	logger.Printf("Changed: %s", request.Params.TextDocument.URI)
	for _, change := range request.Params.ContentChanges {
		state.UpdateDocument(request.Params.TextDocument.URI, change.Text)
	}
	writeResponse(writer, lsp.PublishDiagnosticsNotification{
		Notification: lsp.NewNotification(lsp.PublishDiagnostics),
		Params: lsp.PublishDiagnosticsParams{
			URI:         request.Params.TextDocument.URI,
			Diagnostics: []lsp.Diagnostic{},
		},
	})

}

func handleDidSaveTextDocument(logger *log.Logger, writer io.Writer, state analysis.State, contents []byte) {
	var request lsp.DidSaveTextDocumentNotification
	if err := json.Unmarshal(contents, &request); err != nil {
		logger.Printf("error parsing %s request: %s", lsp.DidSaveTextDocument, err)
		return
	}

	logger.Printf("Saved: %s", request.Params.TextDocument.URI)
	diagnostics := state.ParseDocument(request.Params.TextDocument.URI)
	writeResponse(writer, lsp.PublishDiagnosticsNotification{
		Notification: lsp.NewNotification(lsp.PublishDiagnostics),
		Params: lsp.PublishDiagnosticsParams{
			URI:         request.Params.TextDocument.URI,
			Diagnostics: diagnostics,
		},
	})
	logger.Printf("Published %d diagnostics", len(diagnostics))

}

func handleHover(logger *log.Logger, writer io.Writer, state analysis.State, contents []byte) {
	var request lsp.HoverRequest
	if err := json.Unmarshal(contents, &request); err != nil {
		logger.Printf("error parsing %s request: %s", lsp.Hover, err)
		return
	}

	msg := state.Hover(request.ID, request.Params.TextDocument.URI, request.Params.Position)
	writeResponse(writer, msg)
}

func handleDefinition(logger *log.Logger, writer io.Writer, state analysis.State, contents []byte) {
	var request lsp.DefinitionRequest
	if err := json.Unmarshal(contents, &request); err != nil {
		logger.Printf("error parsing %s request: %s", lsp.Definition, err)
		return
	}

	msg := state.Definition(request.ID, request.Params.TextDocument.URI, request.Params.Position)
	writeResponse(writer, msg)
}

func handleCompletion(logger *log.Logger, writer io.Writer, state analysis.State, contents []byte) {
	var request lsp.CompletionRequest
	if err := json.Unmarshal(contents, &request); err != nil {
		logger.Printf("error parsing %s request: %s", lsp.Completion, err)
		return
	}

	msg := state.TextDocumentCompletion(request.ID, request.Params.TextDocument.URI, request.Params.Position)
	writeResponse(writer, msg)
}

func writeResponse(writer io.Writer, msg any) {
	reply := rpc.EncodeMessage(msg)
	writer.Write([]byte(reply))
}

func main() {
	logger := log.NewLogger(lsp.Name)
	logger.Println("Started lsp")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	state := analysis.NewState()
	writer := os.Stdout

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Decoding error: %s", err)
			continue
		}
		handleMessage(logger, writer, state, method, contents)
	}
}

