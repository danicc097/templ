package parser

import (
	"testing"

	"github.com/a-h/parse"
	"github.com/google/go-cmp/cmp"
)

func TestGoIfExpression(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected GoIfExpression
	}{
		{
			name: "if simple expression",
			input: `{{ if p.Test }}
  %{ "content" }%
{{ end }}`,
			expected: GoIfExpression{
				Expression: Expression{
					GoTempl: true,
					Value:   `p.Test`,
					Range: Range{
						From: Position{
							Index: 6,
							Line:  0,
							Col:   6,
						},
						To: Position{
							Index: 12,
							Line:  0,
							Col:   12,
						},
					},
				},
				Then: []Node{
					Whitespace{Value: "\n  ", GoTempl: true},
					StringExpression{
						Expression: Expression{
							Value: `"content"`,
							Range: Range{
								From: Position{Index: 21, Line: 1, Col: 5},
								To:   Position{Index: 30, Line: 1, Col: 14},
							},
							GoTempl: true,
						},
						TrailingSpace: "\n",
						Gotempl:       true,
					},
				},
			},
		},
		{
			name: "if else",
			input: `{{ if p.A }}
	{{ "A" }}
{{ else }}
	{{ "B" }}
{{ end }}`,
			expected: GoIfExpression{
				Expression: Expression{
					GoTempl: true,
					Value:   `p.A`,
					Range: Range{
						From: Position{
							Index: 6,
							Line:  0,
							Col:   6,
						},
						To: Position{
							Index: 9,
							Line:  0,
							Col:   9,
						},
					},
				},
				Then: []Node{
					Whitespace{Value: "\n\t", GoTempl: true},
					GoCode{
						Expression: Expression{
							Value: `"A"`,
							Range: Range{
								From: Position{Index: 17, Line: 1, Col: 4},
								To:   Position{Index: 20, Line: 1, Col: 7},
							},
						},
						TrailingSpace: "\n",
					},
				},
				Else: []Node{
					Whitespace{Value: "\n\t", GoTempl: true},
					GoCode{
						Expression: Expression{
							Value: `"B"`,
							Range: Range{
								From: Position{Index: 39, Line: 3, Col: 4},
								To:   Position{Index: 42, Line: 3, Col: 7},
							},
						},
						TrailingSpace: "\n",
					},
				},
			},
		},
		{
			name: "if expressions can have a space after the opening brace",
			input: `{{ if p.Test }}
  text
{{ end }}`,
			expected: GoIfExpression{
				Expression: Expression{
					GoTempl: true,
					Value:   `p.Test`,
					Range: Range{
						From: Position{
							Index: 6,
							Line:  0,
							Col:   6,
						},
						To: Position{
							Index: 12,
							Line:  0,
							Col:   12,
						},
					},
				},
				Then: []Node{
					Whitespace{Value: "\n  ", GoTempl: true},
					Text{
						GoTempl: true,
						Value:   "text",
						Range: Range{
							From: Position{Index: 18, Line: 1, Col: 2},
							To:   Position{Index: 22, Line: 1, Col: 6},
						},
						TrailingSpace: SpaceVertical,
					},
				},
			},
		},
		{
			name: "if else, without spaces",
			input: `{{ if p.A}}
	{{ "A" }}
{{ else }}
	{{ "B" }}
{{ end }}`,
			expected: GoIfExpression{
				Expression: Expression{
					GoTempl: true,
					Value:   `p.A`,
					Range: Range{
						From: Position{
							Index: 6,
							Line:  0,
							Col:   6,
						},
						To: Position{
							Index: 9,
							Line:  0,
							Col:   9,
						},
					},
				},
				Then: []Node{
					Whitespace{Value: "\n\t", GoTempl: true},
					GoCode{
						Expression: Expression{
							Value: `"A"`,
							Range: Range{
								From: Position{Index: 16, Line: 1, Col: 4},
								To:   Position{Index: 19, Line: 1, Col: 7},
							},
						},
						TrailingSpace: "\n",
					},
				},
				Else: []Node{
					Whitespace{Value: "\n\t", GoTempl: true},
					GoCode{
						Expression: Expression{
							Value: `"B"`,
							Range: Range{
								From: Position{Index: 38, Line: 3, Col: 4},
								To:   Position{Index: 41, Line: 3, Col: 7},
							},
						},
						TrailingSpace: "\n",
					},
				},
			},
		},
		{
			name: "if nested",
			input: `{{ if p.A }}
					{{ if p.B }}
						{{ "C" }}
					{{ end }}
				{{ end }}`,
			expected: GoIfExpression{
				Expression: Expression{
					GoTempl: true,
					Value:   `p.A`,
					Range: Range{
						From: Position{
							Index: 6,
							Line:  0,
							Col:   6,
						},
						To: Position{
							Index: 9,
							Line:  0,
							Col:   9,
						},
					},
				},
				Then: []Node{
					Whitespace{Value: "\t\t\t\t\t", GoTempl: true},
					GoIfExpression{
						Expression: Expression{
							GoTempl: true,
							Value:   `p.B`,
							Range: Range{
								From: Position{
									Index: 20,
									Line:  1,
									Col:   10,
								},
								To: Position{
									Index: 23,
									Line:  1,
									Col:   13,
								},
							},
						},
						Then: []Node{
							Whitespace{Value: "\t\t\t\t\t\t", GoTempl: true},
							GoCode{
								Expression: Expression{
									Value: `"C"`,
									Range: Range{
										From: Position{Index: 34, Line: 2, Col: 12},
										To:   Position{Index: 37, Line: 2, Col: 15},
									},
								},
								TrailingSpace: "\n\t\t\t\t\t",
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok, err := goIfExpression.Parse(parse.NewInput(tt.input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !ok {
				t.Fatal("expected parser to succeed, but it didn't")
			}

			// Ignore ranges in comparison, they are tested elsewhere.
			if diff := cmp.Diff(tt.expected, result, cmp.AllowUnexported(GoIfExpression{}, StringExpression{}, Element{}, Whitespace{}, Expression{
				GoTempl: true,
			}, Position{})); diff != "" {
				t.Errorf("unexpected result, diff (-want +got):\n%s", diff)
			}
		})
	}
}
