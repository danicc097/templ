package parser

import (
	"github.com/a-h/parse"
)

var gostringExpression = parse.Func(func(pi *parse.Input) (n Node, ok bool, err error) {
	start := pi.Index()

	// do not remove leading whitespace, since string expressions in go templates can be inlined
	// and whitespace is important

	// Attempt to parse the prefix first.
	if _, ok, err = openGotemplStringExprWithOptionalPadding.Parse(pi); err != nil || !ok {
		pi.Seek(start)
		return nil, false, err
	}

	// Once we have a prefix, we must have an expression that returns a string, with optional err.
	r := StringExpression{GoTempl: true}
	if r.Expression, r.GoTemplEndMarker, err = parseGoSliceArgs(pi, "}%"); err != nil || !ok {
		pi.Seek(start) // not an expression that returns a string, might be just text.
		return r, false, err
	}
	r.Expression.GoTempl = true

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
