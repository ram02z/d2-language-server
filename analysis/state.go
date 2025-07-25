package analysis

import (
	"context"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"

	"github.com/ram02z/d2-language-server/log"
	"github.com/ram02z/d2-language-server/lsp"
	"oss.terrastruct.com/d2/d2ast"
	"oss.terrastruct.com/d2/d2format"
	"oss.terrastruct.com/d2/d2lib"
	"oss.terrastruct.com/d2/d2lsp"
	"oss.terrastruct.com/d2/d2parser"
)

type State struct {
	Documents        map[lsp.DocumentURI]Document
	WorkspaceFolders map[lsp.URI]Workspace
	logger           *log.Logger
}

type Workspace struct {
	Name  string
	Files []string
}

type Document struct {
	Text   string
	AST    *d2ast.Map
	Errors []d2ast.Error
}

func NewState(logger *log.Logger) State {
	return State{
		Documents:        map[lsp.DocumentURI]Document{},
		WorkspaceFolders: map[lsp.URI]Workspace{},
		logger:           logger,
	}
}

func (s *State) AddWorkspaceFolders(folders []lsp.WorkspaceFolder) {
	for _, folder := range folders {
		folderPaths := findFilesByExt(folder.URI.Filename(), ".d2")
		s.WorkspaceFolders[folder.URI] = Workspace{
			Name:  folder.Name,
			Files: folderPaths,
		}
		s.logger.Printf("added '%s' to workspace", folder.URI)
	}
}

func (s *State) RemoveWorkspaceFolders(folders []lsp.WorkspaceFolder) {
	for _, folder := range folders {
		delete(s.WorkspaceFolders, folder.URI)
		s.logger.Printf("removed '%s' from workspace", folder.URI)
	}
}

func (s *State) OpenDocument(uri lsp.DocumentURI, text string) []lsp.Diagnostic {
	ctx := context.Background()
	document := parseDocument(ctx, text)
	s.Documents[uri] = document

	return getDiagnosticsFromAST(document.Errors)
}

func (s *State) UpdateDocument(uri lsp.DocumentURI, text string) []lsp.Diagnostic {
	ctx := context.Background()
	document := parseDocument(ctx, text)
	s.Documents[uri] = document

	return getDiagnosticsFromAST(document.Errors)
}

func (s *State) RemoveDocument(uri lsp.DocumentURI) {
	delete(s.Documents, uri)
}

func (s *State) UpdateFile(path string, event lsp.FileChangeType) {
	// TODO: This is inefficient and doesn't scale well with many workspace folders.
	// A better approach would be to find the parent workspace folder directly from the file's path
	// instead of iterating through all of them. This could be done by iterating up the file's
	// path components or by using a more efficient data structure for lookups (like a trie).
	for uri, workspace := range s.WorkspaceFolders {
		if !strings.HasPrefix(path, uri.Filename()) {
			continue
		}

		switch event {
		case lsp.Created:
			workspace.Files = append(workspace.Files, path)
			s.WorkspaceFolders[uri] = workspace
			s.logger.Printf("added %s to %s", path, workspace.Name)
		case lsp.Deleted:
			for i, file := range workspace.Files {
				if file == path {
					workspace.Files = slices.Delete(workspace.Files, i, i+1)
					s.WorkspaceFolders[uri] = workspace
					s.logger.Printf("removed %s from %s", path, workspace.Name)
					break
				}
			}
		}
	}
}

func (s *State) Hover(id any, uri lsp.DocumentURI, position lsp.Position) lsp.HoverResponse {
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

func (s *State) Definition(id any, uri lsp.DocumentURI, position lsp.Position) lsp.DefinitionResponse {
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

func (s *State) ImportCompletion(id any, uri lsp.DocumentURI, position lsp.Position) lsp.CompletionResponse {
	var files []string
	path := uri.Filename()
	root := filepath.Dir(path)
	if len(s.WorkspaceFolders) == 0 {
		files = findFilesByExt(root, ".d2")
	} else {
		for _, workspace := range s.WorkspaceFolders {
			files = append(files, workspace.Files...)
		}
	}

	items := []lsp.CompletionItem{}
	for _, file := range files {
		if file == path {
			continue
		}

		relPath, err := filepath.Rel(root, file)
		if err != nil {
			s.logger.Printf("d2 does not support absolute imports: %v", err)
			continue
		}

		// D2 formatter gets rid of the suffix
		fileName := strings.TrimSuffix(relPath, filepath.Ext(relPath))

		items = append(items, lsp.CompletionItem{
			Label: fileName,
			Kind:  lsp.CompletionItemKindFile,
		})
	}

	response := lsp.CompletionResponse{
		Response: lsp.NewResponse(id),
		Result:   items,
	}

	return response
}

func (s *State) TextDocumentCompletion(id any, uri lsp.DocumentURI, position lsp.Position) lsp.CompletionResponse {
	d2Items, err := d2lsp.GetCompletionItems(s.Documents[uri].Text, position.Line, position.Character)
	if err != nil {
		s.logger.Printf("Error while getting completion items: %v", err)
		return lsp.CompletionResponse{
			Response: lsp.NewResponse(id),
			Result:   []lsp.CompletionItem{},
		}
	}

	lspItems := make([]lsp.CompletionItem, len(d2Items))
	for i, d2Item := range d2Items {
		lspItems[i] = mapToLspCompletionItem(d2Item)
	}

	return lsp.CompletionResponse{
		Response: lsp.NewResponse(id),
		Result:   lspItems,
	}
}

func (s *State) Format(id any, uri lsp.DocumentURI) lsp.FormattingResponse {
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

func findFilesByExt(root, ext string) []string {
	files := []string{}
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(d.Name()) == ext {
			files = append(files, path)
		}
		return nil
	})

	return files
}

func mapToLspCompletionItem(d2Item d2lsp.CompletionItem) lsp.CompletionItem {
	return lsp.CompletionItem{
		Label:      d2Item.Label,
		Detail:     d2Item.Detail,
		Kind:       mapToLspCompletionItemKind(d2Item.Kind),
		InsertText: d2Item.InsertText,
	}
}

func mapToLspCompletionItemKind(d2Kind d2lsp.CompletionKind) lsp.CompletionItemKind {
	if d2Kind == d2lsp.StyleCompletion {
		return lsp.CompletionItemKindProperty
	}
	return lsp.CompletionItemKindKeyword
}
