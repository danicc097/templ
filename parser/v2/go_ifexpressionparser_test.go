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
			name: "if: simple expression",
			input: `{{ if p.Test }}
<span>
  %{ "span content" }%
</span>
{{ end }}`,
			expected: GoIfExpression{
				Expression: Expression{
					Value: `p.Test`,
					Range: Range{
						From: Position{
							Index: 4,
							Line:  0,
							Col:   4,
						},
						To: Position{
							Index: 10,
							Line:  0,
							Col:   10,
						},
					},
				},
				Then: []Node{
					Element{
						Name: "span",
						NameRange: Range{
							From: Position{Index: 16, Line: 1, Col: 2},
							To:   Position{Index: 20, Line: 1, Col: 6},
						},

						Children: []Node{
							Whitespace{Value: "\n  "},
							StringExpression{
								Expression: Expression{
									Value: `"span content"`,
									Range: Range{
										From: Position{
											Index: 28,
											Line:  2,
											Col:   6,
										},
										To: Position{
											Index: 42,
											Line:  2,
											Col:   20,
										},
									},
								},
								TrailingSpace: SpaceVertical,
							},
						},
						IndentChildren: true,
						TrailingSpace:  SpaceVertical,
					},
				},
			},
		},
		{
			name: "if: else",
			input: `{{ if p.A }}
	{{ "A" }}
{{ else }}
	{{ "B" }}
{{ end }}`,
			expected: GoIfExpression{
				Expression: Expression{
					Value: `p.A`,
					Range: Range{
						From: Position{
							Index: 4,
							Line:  0,
							Col:   4,
						},
						To: Position{
							Index: 7,
							Line:  0,
							Col:   7,
						},
					},
				},
				Then: []Node{
					Whitespace{Value: "\t"},
					StringExpression{
						Expression: Expression{
							Value: `"A"`,
							Range: Range{
								From: Position{
									Index: 13,
									Line:  1,
									Col:   4,
								},
								To: Position{
									Index: 16,
									Line:  1,
									Col:   7,
								},
							},
						},
						TrailingSpace: SpaceVertical,
					},
				},
				Else: []Node{
					StringExpression{
						Expression: Expression{
							Value: `"B"`,
							Range: Range{
								From: Position{
									Index: 35,
									Line:  3,
									Col:   4,
								},
								To: Position{
									Index: 38,
									Line:  3,
									Col:   7,
								},
							},
						},
						TrailingSpace: SpaceVertical,
					},
				},
			},
		},
		{
			name: "if: expressions can have a space after the opening brace",
			input: `{{ if p.Test }}
  text
{{ end }}`,
			expected: GoIfExpression{
				Expression: Expression{
					Value: `p.Test`,
					Range: Range{
						From: Position{
							Index: 4,
							Line:  0,
							Col:   4,
						},
						To: Position{
							Index: 10,
							Line:  0,
							Col:   10,
						},
					},
				},
				Then: []Node{
					Whitespace{Value: "  "},
					Text{
						Value: "text",
						Range: Range{
							From: Position{Index: 18, Line: 1, Col: 4},
							To:   Position{Index: 22, Line: 1, Col: 8},
						},
						TrailingSpace: SpaceVertical,
					},
				},
			},
		},
		{
			name: "if: simple expression, without spaces",
			input: `{{ if p.Test }}
<span>
  {{ "span content" }}
</span>
{{ end }}`,
			expected: GoIfExpression{
				Expression: Expression{
					Value: `p.Test`,
					Range: Range{
						From: Position{
							Index: 4,
							Line:  0,
							Col:   4,
						},
						To: Position{
							Index: 10,
							Line:  0,
							Col:   10,
						},
					},
				},
				Then: []Node{
					Element{
						Name: "span",
						NameRange: Range{
							From: Position{Index: 16, Line: 1, Col: 2},
							To:   Position{Index: 20, Line: 1, Col: 6},
						},

						Children: []Node{
							Whitespace{Value: "\n  "},
							StringExpression{
								Expression: Expression{
									Value: `"span content"`,
									Range: Range{
										From: Position{
											Index: 28,
											Line:  2,
											Col:   6,
										},
										To: Position{
											Index: 42,
											Line:  2,
											Col:   20,
										},
									},
								},
								TrailingSpace: SpaceVertical,
							},
						},
						IndentChildren: true,
						TrailingSpace:  SpaceVertical,
					},
				},
			},
		},
		{
			name: "if: else, without spaces",
			input: `{{ if p.A}}
	{{ "A" }}
{{ else }}
	{{ "B" }}
{{ end }}`,
			expected: GoIfExpression{
				Expression: Expression{
					Value: `p.A`,
					Range: Range{
						From: Position{
							Index: 4,
							Line:  0,
							Col:   4,
						},
						To: Position{
							Index: 7,
							Line:  0,
							Col:   7,
						},
					},
				},
				Then: []Node{
					Whitespace{Value: "\t"},
					StringExpression{
						Expression: Expression{
							Value: `"A"`,
							Range: Range{
								From: Position{
									Index: 12,
									Line:  1,
									Col:   4,
								},
								To: Position{
									Index: 15,
									Line:  1,
									Col:   7,
								},
							},
						},
						TrailingSpace: SpaceVertical,
					},
				},
				Else: []Node{
					StringExpression{
						Expression: Expression{
							Value: `"B"`,
							Range: Range{
								From: Position{
									Index: 33,
									Line:  3,
									Col:   4,
								},
								To: Position{
									Index: 36,
									Line:  3,
									Col:   7,
								},
							},
						},
						TrailingSpace: SpaceVertical,
					},
				},
			},
		},
		{
			name: "if: nested",
			input: `{{ if p.A }}
					{{ if p.B }}
						<div>{{ "B" }}</div>
					{{ end }}
				{{ end }}`,
			expected: GoIfExpression{
				Expression: Expression{
					Value: `p.A`,
					Range: Range{
						From: Position{
							Index: 4,
							Line:  0,
							Col:   4,
						},
						To: Position{
							Index: 7,
							Line:  0,
							Col:   7,
						},
					},
				},
				Then: []Node{
					Whitespace{Value: "\t\t\t\t\t"},
					GoIfExpression{
						Expression: Expression{
							Value: `p.B`,
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
							Whitespace{Value: "\t\t\t\t\t\t"},
							Element{
								Name: "div",
								NameRange: Range{
									From: Position{Index: 30, Line: 2, Col: 12},
									To:   Position{Index: 33, Line: 2, Col: 15},
								},
								Children: []Node{
									StringExpression{
										Expression: Expression{
											Value: `"B"`,
											Range: Range{
												From: Position{
													Index: 35,
													Line:  2,
													Col:   17,
												},
												To: Position{
													Index: 38,
													Line:  2,
													Col:   20,
												},
											},
										},
										TrailingSpace: SpaceNone,
									},
								},
								TrailingSpace: SpaceNone,
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
			if diff := cmp.Diff(tt.expected, result, cmp.AllowUnexported(GoIfExpression{}, StringExpression{}, Element{}, Whitespace{}, Expression{}, Position{})); diff != "" {
				t.Errorf("unexpected result, diff (-want +got):\n%s", diff)
			}
		})
	}
}
