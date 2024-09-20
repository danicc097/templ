package parser

import (
	"testing"

	"github.com/a-h/parse"
	"github.com/google/go-cmp/cmp"
)

func TestForExpressionParser_Go(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			name: "for: simple",
			input: `{{ for _, item := range p.Items }}
					{ item }
				{{ end }}`,
			expected: GoForExpression{
				Expression: Expression{
					Value: `_, item := range p.Items`,
					Range: Range{
						From: Position{
							Index: 7,
							Line:  0,
							Col:   7,
						},
						To: Position{
							Index: 31,
							Line:  0,
							Col:   31,
						},
					},
				},
				Children: []Node{
					Whitespace{Value: "\n\t\t\t\t\t"},
					StringExpression{
						Expression: Expression{Value: "item", Range: Range{
							From: Position{Index: 42, Line: 1, Col: 7},
							To:   Position{Index: 46, Line: 1, Col: 11},
						}},
						TrailingSpace: SpaceVertical,
					},
				},
			},
		},
		{
			name:  "for: no newlines",
			input: `{{ for _, item := range p.Items }}{item}{{ end }}`,
			expected: GoForExpression{
				Expression: Expression{
					Value: "_, item := range p.Items",
					Range: Range{
						From: Position{
							Index: 7,
							Line:  0,
							Col:   7,
						},
						To: Position{
							Index: 31,
							Line:  0,
							Col:   31,
						},
					},
				}, Children: []Node{
					StringExpression{Expression: Expression{Value: "item", Range: Range{
						From: Position{Index: 35, Line: 0, Col: 35},
						To:   Position{Index: 39, Line: 0, Col: 39},
					}}},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			input := parse.NewInput(tt.input)
			actual, ok, err := goForExpression.Parse(input)
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

func TestIncompleteFor_Go(t *testing.T) {
	t.Run("no body and no end is ignored", func(t *testing.T) {
		input := parse.NewInput(`{{ for _, item := range p.Items }}`)
		_, ok, err := goForExpression.Parse(input)
		if err == nil {
			t.Fatalf("expected error but got none")
		}
		if ok {
			t.Fatal("expected a non match, but got a match")
		}
	})
	t.Run("no end is ignored", func(t *testing.T) {
		input := parse.NewInput(`{{ for _, item := range p.Items }}{item}`)
		_, ok, err := goForExpression.Parse(input)
		if err == nil {
			t.Fatalf("expected error but got none")
		}
		if ok {
			t.Fatal("expected a non match, but got a match")
		}
	})
	t.Run("capitalised For is ignored", func(t *testing.T) {
		input := parse.NewInput(`{{ For _, item := range p.Items }}{item}{{ end }}`)
		_, ok, err := goForExpression.Parse(input)
		if err != nil {
			t.Fatalf("expected no error but got %v", err)
		}
		if ok {
			t.Fatal("expected a non match, but got a match")
		}
	})
	t.Run("for without body is ignored", func(t *testing.T) {
		input := parse.NewInput(`{{ For _, item := range p.Items }}{{ end }}`)
		_, ok, err := goForExpression.Parse(input)
		if err != nil {
			t.Fatalf("expected no error but got %v", err)
		}
		if ok {
			t.Fatal("expected a non match, but got a match")
		}
	})
	t.Run("go for expression inside {{}} gives error", func(t *testing.T) {
		input := parse.NewInput(`{{ for _, item := range p.Items { }}{item}{{ } }}`)
		_, ok, err := goForExpression.Parse(input)
		if err == nil {
			t.Fatalf("expected error but got none")
		}
		if ok {
			t.Fatal("expected a non match, but got a match")
		}
	})
	t.Run("go syntax is ignored (2)", func(t *testing.T) {
		input := parse.NewInput(`for _, item := range p.Items{ {item} }`)
		_, ok, err := goForExpression.Parse(input)
		if err != nil {
			t.Fatalf("expected no error but got %v", err)
		}
		if ok {
			t.Fatal("expected a non match, but got a match")
		}
	})
}
