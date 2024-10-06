package parser

import (
	"strings"

	"github.com/a-h/parse"
	"github.com/a-h/templ/parser/v2/goexpression"
)

var gogoCode = parse.Func(func(pi *parse.Input) (n Node, ok bool, err error) {
	src, _ := pi.Peek(-1)
	start := pi.Index()
	// Check the prefix first.
	if _, ok, err = parse.Or(parse.String("{{ "), parse.String("{{")).Parse(pi); err != nil || !ok {
		pi.Seek(start) // might be just text
		return
	}
	hasLineComment := peekPrefix(pi, "//")
	start1 := pi.Index()
	var s string
	if s, ok, err = parse.StringUntil(goTemplCommentEnd).Parse(pi); err != nil || (hasLineComment && !strings.Contains(s, "\n")) {
		err = parse.Error("Use /**/ syntax for gotempl comments", pi.Position())
		return
	}
	pi.Seek(start1)
	var r GoCode
	r.GoTempl = true
	pi2 := parse.NewInput(src)
	_, _, _ = parse.OptionalWhitespace.Parse(pi)
	commentStartPos := pi2.Position()
	_, _, _ = goTemplComment.Parse(pi2)
	commentEndPos := pi2.Position()
	_, _, _ = parse.OptionalWhitespace.Parse(pi2)

	if _, ok, _ = dblCloseBraceWithOptionalPadding.Parse(pi2); ok {
		// there is only a comment, nothing else
		commentStartPosIndex := commentStartPos.Index - start
		commentEndPosIndex := commentEndPos.Index - start
		if commentEndPosIndex-commentStartPosIndex > 0 && commentEndPosIndex <= len(src) {
			commentExpr := src[commentStartPosIndex:commentEndPosIndex]

			// There were only comments and the end of the gotempl expression was found
			// Return them so they can be printed in .templ files (but not in _templ.go)
			r.Expression = NewExpression(commentExpr, commentStartPos, commentEndPos)
		} else {
			// empty {{ }}, delete it
			return GoCode{}, false, nil
		}
		return r, true, nil
	}

	// Once we have a prefix, we must have an expression that returns a string, with optional err.
	l := pi.Position().Line
	if r.Expression, err = parseGo("go code", pi, goexpression.Expression); err != nil {
		return r, false, err
	}

	if l != pi.Position().Line {
		r.Multiline = true
	}

	// Clear any optional whitespace.
	_, _, _ = parse.OptionalWhitespace.Parse(pi)

	// }}
	if _, ok, err = dblCloseBraceWithOptionalPadding.Parse(pi); err != nil || !ok {
		err = parse.Error("go code: missing close braces", pi.Position())
		return
	}

	// Parse trailing whitespace but dont consume it
	wsStart := pi.Index()
	ws, _, err := parse.Whitespace.Parse(pi)
	if err != nil {
		return r, false, err
	}
	r.TrailingSpace, err = NewTrailingSpace(ws)
	if err != nil {
		return r, false, err
	}
	pi.Seek(wsStart)

	return r, true, nil
})
