package parser

import (
	"github.com/a-h/parse"
)

var gostringExpression = parse.Func(func(pi *parse.Input) (n Node, ok bool, err error) {
	// Check the prefix first.
	if _, ok, err = parse.Or(parse.String("(( "), parse.String("((")).Parse(pi); err != nil || !ok {
		return
	}

	// Once we have a prefix, we must have an expression that returns a string, with optional err.
	var r StringExpression
	if r.Expression, err = parseGoSliceArgs(pi, "))"); err != nil {
		return r, false, err
	}
	// fmt.Printf("r.Expression.Value: %v\n", r.Expression.Value)

	// Clear any optional whitespace.
	_, _, _ = parse.OptionalWhitespace.Parse(pi)

	// }
	if _, ok, err = dblCloseParensWithOptionalPadding.Parse(pi); err != nil || !ok {
		err = parse.Error("gostring expression: missing close braces", pi.Position())
		return
	}

	// Parse trailing whitespace.
	ws, _, err := parse.Whitespace.Parse(pi)
	if err != nil {
		return r, false, err
	}
	r.TrailingSpace, err = NewTrailingSpace(ws)
	if err != nil {
		return r, false, err
	}

	return r, true, nil
})
