package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/ram02z/d2-language-server/analysis"
	"github.com/ram02z/d2-language-server/log"
	"github.com/ram02z/d2-language-server/lsp"
	"github.com/ram02z/d2-language-server/rpc"
	"github.com/ram02z/d2-language-server/utils"
)

type HandlerFunc func(*log.Logger, io.Writer, analysis.State, []byte)

var handlers = map[lsp.Method]HandlerFunc{
	lsp.Initialize:                handleInitialize,
	lsp.DidOpenTextDocument:       handleDidOpenTextDocument,
	lsp.DidChangeTextDocument:     handleDidChangeTextDocument,
	lsp.DidCloseTextDocument:      handleDidCloseTextDocument,
	lsp.Hover:                     handleHover,
	lsp.Definition:                handleDefinition,
	lsp.Completion:                handleCompletion,
	lsp.Formatting:                handleFormatting,
	lsp.DidChangeWorkspaceFolders: handleDidChangeWorkspaceFolders,
	lsp.DidChangeWatchedFiles:     handleDidChangeWatchedFiles,
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
		logger.Printf("unsupported method: %s", method)
		return
	}

	logger.Printf("received message with method: %s", method)
	handler(logger, writer, state, contents)
}

func handleInitialize(logger *log.Logger, writer io.Writer, state analysis.State, contents []byte) {
	var request lsp.InitializeRequest
	if err := json.Unmarshal(contents, &request); err != nil {
		logger.Printf("error parsing %s request: %s", lsp.Initialize, err)
		return
	}

	logger.Printf(
		"connected to: %s %s",
		request.Params.ClientInfo.Name,
		request.Params.ClientInfo.Version,
	)

	if folders := request.Params.WorkspaceFolders; folders != nil {
		state.AddWorkspaceFolders(folders)
	}

	msg := lsp.NewInitializeResponse(request.ID)
	if err := writeResponse(writer, msg); err != nil {
		logger.Printf("could not write response: %s", err)
	}

	registerCapabilityRequest := lsp.NewRequestWithParams(
		lsp.ClientRegisterCapability,
		lsp.RegistrationParams{
			Registrations: []lsp.Registration{
				lsp.NewRegistration(
					lsp.DidChangeWatchedFiles,
					lsp.DidChangeWatchedFilesRegistrationOptions{
						Watchers: []lsp.FileSystemWatcher{
							{
								GlobPattern: "**/*.d2",
							},
						},
					},
				),
			},
		},
	)

	if err := writeResponse(writer, registerCapabilityRequest); err != nil {
		logger.Printf("could not register capability: %s", err)
	}
}

func handleDidOpenTextDocument(logger *log.Logger, writer io.Writer, state analysis.State, contents []byte) {
	var request lsp.DidOpenTextDocumentNotification
	if err := json.Unmarshal(contents, &request); err != nil {
		logger.Printf("error parsing %s request: %s", lsp.DidOpenTextDocument, err)
		return
	}

	logger.Printf("opened document: %s", request.Params.TextDocument.URI)
	diagnostics := state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)
	writeResponse(writer, lsp.PublishDiagnosticsNotification{
		Notification: lsp.NewNotification(lsp.PublishDiagnostics),
		Params: lsp.PublishDiagnosticsParams{
			URI:         request.Params.TextDocument.URI,
			Diagnostics: diagnostics,
		},
	})
	logger.Printf("published %d diagnostics", len(diagnostics))
}

func handleDidChangeTextDocument(logger *log.Logger, writer io.Writer, state analysis.State, contents []byte) {
	var request lsp.DidChangeTextDocumentNotification
	if err := json.Unmarshal(contents, &request); err != nil {
		logger.Printf("error parsing %s request: %s", lsp.DidChangeTextDocument, err)
		return
	}

	logger.Printf("changed document: %s", request.Params.TextDocument.URI)
	diagnostics := make([]lsp.Diagnostic, 0)
	contentChangesLen := len(request.Params.ContentChanges)
	// HACK: only considering the final change
	if contentChangesLen > 0 {
		lastChangeEvent := request.Params.ContentChanges[contentChangesLen-1]
		diagnostics = append(diagnostics, state.UpdateDocument(request.Params.TextDocument.URI, lastChangeEvent.Text)...)
	}
	writeResponse(writer, lsp.PublishDiagnosticsNotification{
		Notification: lsp.NewNotification(lsp.PublishDiagnostics),
		Params: lsp.PublishDiagnosticsParams{
			URI:         request.Params.TextDocument.URI,
			Diagnostics: diagnostics,
		},
	})
	logger.Printf("published %d diagnostics", len(diagnostics))
}

func handleDidCloseTextDocument(logger *log.Logger, writer io.Writer, state analysis.State, contents []byte) {
	var request lsp.DidCloseTextDocumentNotification
	if err := json.Unmarshal(contents, &request); err != nil {
		logger.Printf("error parsing %s request: %s", lsp.DidChangeTextDocument, err)
		return
	}

	logger.Printf("closed document: %s", request.Params.TextDocument.URI)
	state.RemoveDocument(request.Params.TextDocument.URI)
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

	var msg lsp.CompletionResponse
	if request.Params.Context.TriggerKind == lsp.TriggerCharacter {
		if request.Params.Context.TriggerCharacter == "@" {
			msg = state.ImportCompletion(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		}
	} else {
		msg = state.TextDocumentCompletion(request.ID, request.Params.TextDocument.URI, request.Params.Position)
	}
	writeResponse(writer, msg)
}

func handleFormatting(logger *log.Logger, writer io.Writer, state analysis.State, contents []byte) {
	var request lsp.FormattingRequest
	if err := json.Unmarshal(contents, &request); err != nil {
		logger.Printf("error parsing %s request: %s", lsp.Formatting, err)
		return
	}

	msg := state.Format(request.ID, request.Params.TextDocument.URI)
	writeResponse(writer, msg)
	logger.Printf("formatted: %s", request.Params.TextDocument.URI)
}

func handleDidChangeWorkspaceFolders(logger *log.Logger, writer io.Writer, state analysis.State, contents []byte) {
	var request lsp.DidChangeWorkspaceFoldersNotifications
	if err := json.Unmarshal(contents, &request); err != nil {
		logger.Printf("error parsing %s request: %s", lsp.DidChangeWorkspaceFolders, err)
		return
	}

	logger.Printf(
		"changed workspace: +%d folders, -%d folders",
		len(request.Params.Event.Added),
		len(request.Params.Event.Removed),
	)

	state.RemoveWorkspaceFolders(request.Params.Event.Removed)
	state.AddWorkspaceFolders(request.Params.Event.Added)
}

func handleDidChangeWatchedFiles(logger *log.Logger, writer io.Writer, state analysis.State, contents []byte) {
	var request lsp.DidChangeWatchedFilesNotification
	if err := json.Unmarshal(contents, &request); err != nil {
		logger.Printf("error parsing %s request: %s", lsp.DidChangeWatchedFiles, err)
		return
	}

	for _, change := range request.Params.Changes {
		path, err := utils.GetPathFromURI(change.URI)
		if err != nil {
			logger.Printf("could not get path from uri: %s", err)
			continue
		}
		state.UpdateFile(path, change.Type)
	}
}

func writeResponse(writer io.Writer, msg any) error {
	reply, err := rpc.EncodeMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to encode messaged: %w", err)
	}
	_, err = writer.Write([]byte(reply))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

func main() {
	logger := log.NewLogger(lsp.Name)
	logger.Println("started lsp")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	state := analysis.NewState(logger)
	writer := os.Stdout

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("decoding error: %s", err)
			continue
		}
		handleMessage(logger, writer, state, method, contents)
	}
}
