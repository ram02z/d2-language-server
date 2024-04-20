package analysis

import (
	"context"
	"fmt"
	"github.com/ram02z/d2-language-server/lsp"
	"oss.terrastruct.com/d2/d2lib"
	"oss.terrastruct.com/d2/d2parser"
)

type State struct {
	Documents map[string]string
}

func NewState() State {
	return State{
		Documents: map[string]string{},
	}
}

func LineRange(line, start, end int) lsp.Range {
	return lsp.Range{
		Start: lsp.Position{
			Line:      line,
			Character: start,
		},
		End: lsp.Position{
			Line:      line,
			Character: end,
		},
	}
}

func getDiagnosticsForFile(text string) []lsp.Diagnostic {
	diagnostics := []lsp.Diagnostic{}
	ctx := context.Background()
	_, err := d2lib.Parse(ctx, text, &d2lib.CompileOptions{
		UTF16Pos: true,
	})

	if err != nil {
		parserErrors := err.(*d2parser.ParseError).Errors
		for _, err := range parserErrors {
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
	}

	return diagnostics
}

func (s *State) OpenDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = text

	return getDiagnosticsForFile(text)
}

func (s *State) UpdateDocument(uri, text string) []lsp.Diagnostic {
	s.Documents[uri] = text

	return getDiagnosticsForFile(text)
}

func (s *State) Hover(id int, uri string, position lsp.Position) lsp.HoverResponse {
	document := s.Documents[uri]

	return lsp.HoverResponse{
		Response: lsp.NewResponse(id),
		Result: lsp.HoverResult{
			Contents: fmt.Sprintf("File: %s, Characters: %d", uri, len(document)),
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
