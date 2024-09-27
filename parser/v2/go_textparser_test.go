package parser

import (
	"testing"

	"github.com/a-h/parse"
	"github.com/google/go-cmp/cmp"
)

func TestGoTextParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Text
	}{
		{
			name:  "Text ends at a gotempl string expression start",
			input: `abcdef%{ "test" }`,
			expected: Text{
				GoTempl: true,
				Value:   "abcdef",
				Range: Range{
					From: Position{Index: 0, Line: 0, Col: 0},
					To:   Position{Index: 6, Line: 0, Col: 6},
				},
			},
		},
		{
			name:  "Text ends at a gotempl expression start",
			input: `abcdef{{ const a := "test" }}`,
			expected: Text{
				GoTempl: true,
				Value:   "abcdef",
				Range: Range{
					From: Position{Index: 0, Line: 0, Col: 0},
					To:   Position{Index: 6, Line: 0, Col: 6},
				},
			},
		},
		{
			name:  "Text may contain brackets and braces",
			input: `abcdef ({}) {}( ) { } ghijk%{ "test" }%`,
			expected: Text{
				GoTempl: true,
				Value:   "abcdef ({}) {}( ) { } ghijk",
				Range: Range{
					From: Position{Index: 0, Line: 0, Col: 0},
					To:   Position{Index: 27, Line: 0, Col: 27},
				},
			},
		},
		{
			name:  "Multiline text is collected line by line",
			input: "Line 1\n  Line 2",
			expected: Text{
				GoTempl: true,
				Value:   "Line 1",
				Range: Range{
					From: Position{Index: 0, Line: 0, Col: 0},
					To:   Position{Index: 6, Line: 0, Col: 6},
				},
				TrailingSpace:    "\n",
				TrailingSpaceLit: "\n  ",
			},
		},
		{
			name:  "Multiline text is collected line by line keeping leading space",
			input: "\t  Line 1\n    Line 2",
			expected: Text{
				GoTempl: true,
				Value:   "\t  Line 1",
				Range: Range{
					From: Position{Index: 0, Line: 0, Col: 0},
					To:   Position{Index: 9, Line: 0, Col: 9},
				},
				TrailingSpace:    "\n",
				TrailingSpaceLit: "\n    ",
			},
		},
		{
			name:  "Multiline text is collected line by line (Windows)",
			input: "Line 1\r\nLine 2",
			expected: Text{
				GoTempl: true,
				Value:   "Line 1",
				Range: Range{
					From: Position{Index: 0, Line: 0, Col: 0},
					To:   Position{Index: 6, Line: 0, Col: 6},
				},
				TrailingSpace:    "\n",
				TrailingSpaceLit: "\r\n",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			input := parse.NewInput(tt.input)
			actual, ok, err := gotextParser.Parse(input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !ok {
				t.Fatalf("unexpected failure for input %q", tt.input)
			}
			if diff := cmp.Diff(tt.expected, actual); diff != "" {
				t.Error(diff)
			}
		})
	}
}
