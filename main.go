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

func main() {
	logger := log.NewLogger(lsp.Name)
	logger.Println("started lsp")

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

func handleMessage(
	logger *log.Logger,
	writer io.Writer,
	state analysis.State,
	method string,
	contents []byte,
) {
	logger.Printf("Received message with method: %s", method)
	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("error parsing %s request: %s", method, err)
			return
		}

		logger.Printf(
			"Connected to: %s %s",
			request.Params.ClientInfo.Name,
			request.Params.ClientInfo.Version,
		)

		msg := lsp.NewInitializeResponse(request.ID)
		writeResponse(writer, msg)
	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("error parsing %s request: %s", method, err)
			return
		}

		logger.Printf("Opened: %s", request.Params.TextDocument.URI)
		diagnostics := state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)
		writeResponse(writer, lsp.PublishDiagnosticsNotification{
			Notification: lsp.Notification{
				RPC:    "2.0",
				Method: "textDocument/publishDiagnostics",
			},
			Params: lsp.PublishDiagnosticsParams{
				URI:         request.Params.TextDocument.URI,
				Diagnostics: diagnostics,
			},
		})
		logger.Printf("Published %d diagnostics", len(diagnostics))
	case "textDocument/didChange":
		var request lsp.DidChangeTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("error parsing %s request: %s", method, err)
			return
		}

		logger.Printf("Changed: %s", request.Params.TextDocument.URI)
		for _, change := range request.Params.ContentChanges {
			diagnostics := state.UpdateDocument(request.Params.TextDocument.URI, change.Text)
			writeResponse(writer, lsp.PublishDiagnosticsNotification{
				Notification: lsp.Notification{
					RPC:    "2.0",
					Method: "textDocument/publishDiagnostics",
				},
				Params: lsp.PublishDiagnosticsParams{
					URI:         request.Params.TextDocument.URI,
					Diagnostics: diagnostics,
				},
			})
			logger.Printf("Published %d diagnostics", len(diagnostics))
		}
	case "textDocument/hover":
		var request lsp.HoverRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("error parsing %s request: %s", method, err)
			return
		}

		msg := state.Hover(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		writeResponse(writer, msg)
	case "textDocument/definition":
		var request lsp.DefinitionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("error parsing %s request: %s", method, err)
			return
		}

		msg := state.Definition(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		writeResponse(writer, msg)
	case "textDocument/completion":
		var request lsp.CompletionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("error parsing %s request: %s", method, err)
			return
		}

		msg := state.TextDocumentCompletion(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		writeResponse(writer, msg)
	}
}

func writeResponse(writer io.Writer, msg any) {
	reply := rpc.EncodeMessage(msg)
	writer.Write([]byte(reply))
}
