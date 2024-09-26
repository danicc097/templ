package parser

import (
	"fmt"
	"os"
	"strconv"

	"github.com/a-h/parse"
)

var untilGoTemplOrNewLine = parse.StringUntil(parse.Any(parse.String("{{"), openGotemplStringExpr, parse.String("\r\n"), parse.Rune('\n')))

var gotextParser = parse.Func(func(pi *parse.Input) (n Node, ok bool, err error) {
	// src, _ := pi.Peek(-1)
	from := pi.Position()

	t := Text{GoTempl: true}
	if t.LeadingSpaceLit, _, err = parse.OptionalWhitespace.Parse(pi); err != nil {
		pi.Seek(from.Index)
		return
	}
	fmt.Fprintf(os.Stderr, "t.LeadingSpaceLit: %s\n", strconv.Quote(t.LeadingSpaceLit))
	// Read until a templ expression opens or line ends.

	if t.Value, ok, err = untilGoTemplOrNewLine.Parse(pi); err != nil || !ok {
		pi.Seek(from.Index)
		return
	}
	if isWhitespace(t.Value) {
		return Whitespace{GoTempl: true, Value: t.Value}, false, nil
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
	wsStart := pi.Index()
	t.TrailingSpaceLit, _, err = parse.Whitespace.Parse(pi)
	if err != nil {
		return t, false, err
	}
	t.TrailingSpace, err = NewTrailingSpace(t.TrailingSpaceLit)
	if err != nil {
		return t, false, err
	}
	pi.Seek(wsStart) // leave whitespace for the next parser so bare text spacing is not formatted by templ

	return t, true, nil
})
