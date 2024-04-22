package analysis

import (
	"strings"

	"github.com/akedrou/textdiff"
	"github.com/ram02z/d2-language-server/lsp"
)

func offsetRange(before string, startOffset, endOffset int) lsp.Range {
	startPos := offsetToPosition(before, startOffset)
	endPos := offsetToPosition(before, endOffset)
	return lsp.Range{
		Start: startPos,
		End:   endPos,
	}
}

func offsetToPosition(text string, offset int) lsp.Position {
	lines := strings.Split(text, "\n")
	var lineIdx, charIdx int
	for i, line := range lines {
		if offset >= len(line)+1 {
			offset -= len(line) + 1
			continue
		}
		lineIdx = i
		charIdx = offset
		break
	}
	return lsp.Position{Line: lineIdx, Character: charIdx}
}

func ComputeTextEdits(before, after string) []lsp.TextEdit {
	edits := textdiff.Strings(before, after)
	if edits == nil {
		return nil
	}

	var result []lsp.TextEdit
	for _, edit := range edits {
		rng := offsetRange(before, edit.Start, edit.End)
		if rng.Start == rng.End && edit.New == "" {
			continue
		}
		result = append(result, lsp.TextEdit{
			Range:   rng,
			NewText: edit.New,
		})
	}

	return result
}
