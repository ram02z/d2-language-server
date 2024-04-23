package analysis

import (
	"context"

	"github.com/ram02z/d2-language-server/log"
	"github.com/ram02z/d2-language-server/lsp"
	"oss.terrastruct.com/d2/d2ast"
	"oss.terrastruct.com/d2/d2format"
	"oss.terrastruct.com/d2/d2lib"
	"oss.terrastruct.com/d2/d2parser"
)

type State struct {
	Documents map[string]Document
	// URI: Name
	WorkspaceFolders map[string]string
	logger           *log.Logger
}

type Document struct {
	Text   string
	AST    *d2ast.Map
	Errors []d2ast.Error
}

func NewState(logger *log.Logger) State {
	return State{
		Documents:        map[string]Document{},
		WorkspaceFolders: map[string]string{},
		logger:           logger,
	}
}

func (s *State) AddWorkspaceFolders(folders []lsp.WorkspaceFolder) {
	for _, folder := range folders {
		s.WorkspaceFolders[folder.URI] = folder.Name
		s.logger.Printf("added '%s' to workspace", folder.URI)
	}
}

func (s *State) RemoveWorkspaceFolders(folders []lsp.WorkspaceFolder) {
	for _, folder := range folders {
		delete(s.WorkspaceFolders, folder.URI)
		s.logger.Printf("removed '%s' from workspace", folder.URI)
	}
}

func (s *State) OpenDocument(uri, text string) []lsp.Diagnostic {
	ctx := context.Background()
	document := parseDocument(ctx, text)
	s.Documents[uri] = document

	return getDiagnosticsFromAST(document.Errors)
}

func (s *State) UpdateDocument(uri, text string) []lsp.Diagnostic {
	ctx := context.Background()
	document := parseDocument(ctx, text)
	s.Documents[uri] = document

	return getDiagnosticsFromAST(document.Errors)
}

func (s *State) Hover(id int, uri string, position lsp.Position) lsp.HoverResponse {
	document := s.Documents[uri]
	node := getNodeUnderCursor(*document.AST, position)

	contents := ""
	if node != nil {
		contents = (*node).GetRange().String()
	}

	return lsp.HoverResponse{
		Response: lsp.NewResponse(id),
		Result: lsp.HoverResult{
			Contents: contents,
		},
	}
}

func (s *State) Definition(id int, uri string, position lsp.Position) lsp.DefinitionResponse {
	return lsp.DefinitionResponse{
		Response: lsp.NewResponse(id),
		Result: lsp.Location{
			URI: uri,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      position.Line - 1,
					Character: 0,
				},
				End: lsp.Position{
					Line:      position.Line - 1,
					Character: 0,
				},
			},
		},
	}
}

func (s *State) TextDocumentCompletion(id int, uri string, position lsp.Position) lsp.CompletionResponse {
	items := []lsp.CompletionItem{
		{
			Label:         "Test",
			Detail:        "Lorem Ipsum",
			Documentation: "Fake documentation",
		},
	}

	response := lsp.CompletionResponse{
		Response: lsp.NewResponse(id),
		Result:   items,
	}

	return response
}

func (s *State) Format(id int, uri string) lsp.FormattingResponse {
	document := s.Documents[uri]

	formattedText := d2format.Format(document.AST)
	result := ComputeTextEdits(document.Text, formattedText)

	response := lsp.FormattingResponse{
		Response: lsp.NewResponse(id),
		Result:   result,
	}

	return response
}

func getDiagnosticsFromAST(errors []d2ast.Error) []lsp.Diagnostic {
	diagnostics := []lsp.Diagnostic{}

	for _, err := range errors {
		diagnostics = append(diagnostics, lsp.Diagnostic{
			Source:  lsp.Name,
			Message: err.Message,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      err.Range.Start.Line,
					Character: err.Range.Start.Column,
				},
				End: lsp.Position{
					Line:      err.Range.End.Line,
					Character: err.Range.End.Column,
				},
			},
			Severity: lsp.Error,
		})
	}

	return diagnostics
}

func parseDocument(ctx context.Context, text string) Document {
	ast, err := d2lib.Parse(ctx, text, &d2lib.CompileOptions{
		UTF16Pos: true,
	})

	errors := []d2ast.Error{}
	if err != nil {

		errors = err.(*d2parser.ParseError).Errors
	}

	return Document{
		Text:   text,
		AST:    ast,
		Errors: errors,
	}
}

func getNodeUnderCursor(ast d2ast.Map, position lsp.Position) *d2ast.MapNode {
	for _, nodeBox := range ast.Nodes {
		node := nodeBox.Unbox()
		nodeRange := node.GetRange()
		if position.Line >= nodeRange.Start.Line && position.Line <= nodeRange.End.Line &&
			position.Character >= nodeRange.Start.Column && position.Character <= nodeRange.End.Column {
			return &node
		}
	}

	return nil
}
