package parser

import (
	"testing"

	"github.com/a-h/parse"
	"github.com/google/go-cmp/cmp"
)

func TestGoStringExpressionParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected StringExpression
	}{
		{
			name:  "basic expression 1",
			input: `%{ fmt.Sprintf("%s", "this") }%`,
			expected: StringExpression{
				GoTempl: true,
				Expression: Expression{
					GoTempl: true,
					Value:   `fmt.Sprintf("%s", "this")`,
					Range: Range{
						From: Position{
							Index: 3,
							Line:  0,
							Col:   3,
						},
						To: Position{
							Index: 28,
							Line:  0,
							Col:   28,
						},
					},
				},
			},
		},
		{
			name:  "basic expression 2",
			input: `%{ "this" }%`,
			expected: StringExpression{
				GoTempl: true,
				Expression: Expression{
					GoTempl: true,
					Value:   `"this"`,
					Range: Range{
						From: Position{
							Index: 3,
							Line:  0,
							Col:   3,
						},
						To: Position{
							Index: 9,
							Line:  0,
							Col:   9,
						},
					},
				},
			},
		},
		{
			name:  "no spaces",
			input: `%{"this"}%`,
			expected: StringExpression{
				GoTempl: true,
				Expression: Expression{
					GoTempl: true,
					Value:   `"this"`,
					Range: Range{
						From: Position{
							Index: 2,
							Line:  0,
							Col:   2,
						},
						To: Position{
							Index: 8,
							Line:  0,
							Col:   8,
						},
					},
				},
			},
		},
		{
			name: "multiple lines",
			input: `%{ test{}.Call(a,
		b,
	  c) }%`,
			expected: StringExpression{
				GoTempl: true,
				Expression: Expression{
					GoTempl: true,
					Value:   "test{}.Call(a,\n\t\tb,\n\t  c)",
					Range: Range{
						From: Position{
							Index: 3,
							Line:  0,
							Col:   3,
						},
						To: Position{
							Index: 28,
							Line:  2,
							Col:   5,
						},
					},
				},
			},
		},
		{
			name:  "basic expression with end marker",
			input: `%{ "this" -}%`,
			expected: StringExpression{
				GoTempl:          true,
				GoTemplEndMarker: true,
				Expression: Expression{
					GoTempl: true,
					Value:   `"this"`,
					Range: Range{
						From: Position{
							Index: 3,
							Line:  0,
							Col:   3,
						},
						To: Position{
							Index: 9,
							Line:  0,
							Col:   9,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			input := parse.NewInput(tt.input)
			an, ok, err := gostringExpression.Parse(input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !ok {
				t.Fatalf("unexpected failure for input %q", tt.input)
			}
			actual := an.(StringExpression)
			if diff := cmp.Diff(tt.expected, actual); diff != "" {
				t.Error(diff)
			}

			// Check the index.
			cut := tt.input[actual.Expression.Range.From.Index:actual.Expression.Range.To.Index]
			if tt.expected.Expression.Value != cut {
				t.Errorf("range, expected %q, got %q", tt.expected.Expression.Value, cut)
			}
		})
	}
}
