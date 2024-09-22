package parser

import (
	"fmt"

	"github.com/a-h/parse"
)

var gostringExpression = parse.Func(func(pi *parse.Input) (n Node, ok bool, err error) {
	start := pi.Index()
	src, _ := pi.Peek(-1)
	// Attempt to parse the prefix first.
	if _, ok, err = openGotemplStringExprWithOptionalPadding.Parse(pi); err != nil || !ok {
		pi.Seek(start)
		return nil, false, err
	}
	fmt.Printf("gostringExpression src: %v\n", string(src))

	// Once we have a prefix, we must have an expression that returns a string, with optional err.
	var r StringExpression
	if r.Expression, err = parseGoSliceArgs(pi, "}%"); err != nil || !ok {
		fmt.Printf("r.Expression.Value: %v\n", r.Expression.Value)
		pi.Seek(start) // not an expression that returns a string, might be just text.
		return r, false, err
	}
	fmt.Printf("r.Expression.Value: %v\n", r.Expression.Value)

	// Clear any optional whitespace.
	_, _, _ = parse.OptionalWhitespace.Parse(pi)

	if _, ok, err = closeGotemplStringExprWithOptionalPadding.Parse(pi); err != nil || !ok {
		pi.Seek(start) // not an expression that returns a string, might be just text.
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
