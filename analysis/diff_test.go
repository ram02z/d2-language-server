package analysis_test

import (
	"testing"

	"github.com/ram02z/d2-language-server/analysis"
	"github.com/ram02z/d2-language-server/lsp"
)

func TestComputeTextEdits(t *testing.T) {
	tests := []struct {
		name     string
		before   string
		after    string
		expected []lsp.TextEdit
	}{
		{
			name:     "no change",
			before:   "hello world",
			after:    "hello world",
			expected: nil,
		},
		{
			name:   "multiline edit 1",
			before: "line1\n\n",
			after:  "line1\n",
			expected: []lsp.TextEdit{
				{
					Range: lsp.Range{
						Start: lsp.Position{Line: 1, Character: 0},
						End:   lsp.Position{Line: 2, Character: 0},
					},
					NewText: "",
				},
			},
		},
		{
			name:   "multiline edit 2",
			before: "line1\n\nline2",
			after:  "line1\nline2",
			expected: []lsp.TextEdit{
				{
					Range: lsp.Range{
						Start: lsp.Position{Line: 1, Character: 0},
						End:   lsp.Position{Line: 2, Character: 0},
					},
					NewText: "",
				},
			},
		},
		{
			name:   "multiline edit 3",
			before: "\nline1\n\nline2",
			after:  "line1\nline2",
			expected: []lsp.TextEdit{
				{
					Range: lsp.Range{
						Start: lsp.Position{Line: 0, Character: 0},
						End:   lsp.Position{Line: 1, Character: 0},
					},
					NewText: "",
				},
				{
					Range: lsp.Range{
						Start: lsp.Position{Line: 1, Character: 5},
						End:   lsp.Position{Line: 2, Character: 0},
					},
					NewText: "",
				},
			},
		},
		{
			name:   "line edit 1",
			before: "helloXworldY",
			after:  "helloYworldX",
			expected: []lsp.TextEdit{
				{
					Range: lsp.Range{
						Start: lsp.Position{Line: 0, Character: 5},
						End:   lsp.Position{Line: 0, Character: 6},
					},
					NewText: "Y",
				},
				{
					Range: lsp.Range{
						Start: lsp.Position{Line: 0, Character: 11},
						End:   lsp.Position{Line: 0, Character: 12},
					},
					NewText: "X",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := analysis.ComputeTextEdits(test.before, test.after)
			if !equalTextEdits(actual, test.expected) {
				t.Errorf("ComputeTextEdits(%q, %q) = %v, expected %v", test.before, test.after, actual, test.expected)
			}
		})
	}
}

func equalTextEdits(a, b []lsp.TextEdit) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
