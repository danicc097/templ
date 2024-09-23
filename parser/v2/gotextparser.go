package parser

import (
	"github.com/a-h/parse"
)

var goTemplOrNewLine = parse.Any(parse.String("{{"), openGotemplStringExpr, parse.String("\r\n"), parse.Rune('\n'))

var gotextParser = parse.Func(func(pi *parse.Input) (n Node, ok bool, err error) {
	from := pi.Position()
	// to, _ := pi.Peek(-1)
	// fmt.Printf("gotextParser : %v\n", to)

	// Read until a templ expression opens or line ends.

	// FIXME gotextParser: finding a line that starts with } should be go text,
	// until we parse a "\n}\n\n" or "\n}\nEOF", or "{{".
	// start which can be found inline.
	// if the %{}% or {{}} expression parser fails, then it defaults to goTextParser
	// again so there shouldnt be issues misinterpreting go text as gotempl exp
	// alternative:  parse go expressions individually the same way as html elements - much more troublesome.
	var t Text
	if t.Value, ok, err = parse.StringUntil(goTemplOrNewLine).Parse(pi); err != nil || !ok {
		return
	}
	if isWhitespace(t.Value) {
		return t, false, nil
	}
	if _, ok = pi.Peek(1); !ok {
		err = parse.Error("gotextParser: unterminated text: expected templ expression open, or newline", from)
		return
	}
	t.Range = NewRange(from, pi.Position())

	// Elide any void element closing tags.
	if _, _, err = voidElementCloser.Parse(pi); err != nil {
		return
	}

	// Parse trailing whitespace.
	ws, _, err := parse.Whitespace.Parse(pi)
	if err != nil {
		return t, false, err
	}
	t.TrailingSpace, err = NewTrailingSpace(ws)
	if err != nil {
		return t, false, err
	}

	return t, true, nil
})
