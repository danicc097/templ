package parser

import (
	"strings"
	"testing"

	"github.com/a-h/parse"
	"github.com/google/go-cmp/cmp"
)

func TestGoGoCodeParser(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    GoCode
		errContains string
	}{
		{
			name:  "basic expression",
			input: `{{ p := "this" }}`,
			expected: GoCode{
				GoTempl: true,
				Expression: Expression{
					Value: `p := "this"`,
					Range: Range{
						From: Position{
							Index: 3,
							Line:  0,
							Col:   3,
						},
						To: Position{
							Index: 14,
							Line:  0,
							Col:   14,
						},
					},
				},
			},
		},
		{
			name:  "basic expression, no space",
			input: `{{p:="this"}}`,
			expected: GoCode{
				GoTempl: true,
				Expression: Expression{
					Value: `p:="this"`,
					Range: Range{
						From: Position{
							Index: 2,
							Line:  0,
							Col:   2,
						},
						To: Position{
							Index: 11,
							Line:  0,
							Col:   11,
						},
					},
				},
			},
		},
		{
			name: "multiline function decl",
			input: `{{
				p := func() {
					dosomething()
				}
			}}`,
			expected: GoCode{
				GoTempl: true,
				Expression: Expression{
					Value: `p := func() {
					dosomething()
				}`,
					Range: Range{
						From: Position{
							Index: 7,
							Line:  1,
							Col:   4,
						},
						To: Position{
							Index: 45,
							Line:  3,
							Col:   5,
						},
					},
				},
				Multiline: true,
			},
		},
		{
			name: "comments in expression 1",
			input: `{{ /* Comment at the start of expression. */
	one := "one"
	two := "two"
	// Comment in middle of expression.
	four := "four"
	// Comment at end of expression.
}}`,
			expected: GoCode{
				GoTempl: true,
				Expression: Expression{
					Value: `/* Comment at the start of expression. */
	one := "one"
	two := "two"
	// Comment in middle of expression.
	four := "four"
	// Comment at end of expression.`,
					Range: Range{
						From: Position{Index: 3, Line: 0, Col: 3},
						To:   Position{Index: 159, Line: 5, Col: 33},
					},
				},
				TrailingSpace: SpaceNone,
				Multiline:     true,
			},
		},
		{
			name:  "line comments in expression 1",
			input: `{{ // Comment only }}`,
			expected: GoCode{
				GoTempl: true,
				Expression: Expression{
					Value: "// Comment only",
					Range: Range{From: Position{Index: 3, Line: 0, Col: 3}, To: Position{Index: 18, Line: 0, Col: 18}},
				},
				TrailingSpace: SpaceNone,
				Multiline:     false,
			},
			errContains: "Use /**/ syntax for gotempl comments: line 0, col 3",
		},
		{
			name: "line comments in expression 2",
			input: `{{// Comment only
				// Comment only
				}}`,
			expected: GoCode{
				GoTempl: true,
				Expression: Expression{
					Value: "// Comment only",
					Range: Range{From: Position{Index: 3, Line: 0, Col: 3}, To: Position{Index: 18, Line: 0, Col: 18}},
				},
				TrailingSpace: SpaceNone,
				Multiline:     false,
			},
			errContains: "Use /**/ syntax for gotempl comments: line 0, col 2",
		},
		{
			name:  "block comments in expression 1",
			input: `{{ /* Comment only */ }}`,
			expected: GoCode{
				GoTempl: true,
				Expression: Expression{
					Value: "/* Comment only */",
					Range: Range{From: Position{Index: 3, Line: 0, Col: 3}, To: Position{Index: 21, Line: 0, Col: 21}},
				},
				TrailingSpace: SpaceNone,
				Multiline:     false,
			},
		},
		{
			name: "block comments in expression 1",
			input: `{{ /* Comment only
				*/ }}`,
			expected: GoCode{
				GoTempl: true,
				Expression: Expression{
					Value: `/* Comment only
				*/`,
					Range: Range{From: Position{Index: 3, Line: 0, Col: 3}, To: Position{Index: 25, Line: 1, Col: 6}},
				},
				TrailingSpace: SpaceNone,
				Multiline:     true,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			input := parse.NewInput(tt.input)
			an, ok, err := gogoCode.Parse(input)
			if (err != nil) != (tt.errContains != "") {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.errContains != "" {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Fatalf("expected error to contain %q, got %q", tt.errContains, err.Error())
				}
				return
			}
			if !ok {
				t.Fatalf("unexpected failure for input %q", tt.input)
			}
			if (an == nil) != (tt.expected.Expression.Value == "") {
				t.Fatalf("no node was returned, but an expression value was expected")
			}
			if an == nil {
				return
			}
			actual := an.(GoCode)
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
