package parser

import "github.com/a-h/parse"

var gowhitespaceExpression = parse.Func(func(pi *parse.Input) (n Node, ok bool, err error) {
	var r Whitespace
	r.GoTempl = true
	if r.Value, ok, err = parse.OptionalWhitespace.Parse(pi); err != nil || !ok {
		return
	}
	return r, len(r.Value) > 0, nil
})
