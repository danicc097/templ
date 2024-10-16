package parser

import (
	"github.com/a-h/parse"
)

var (
	goSingleLineCommentStart = parse.String("//")
	goSingleLineCommentEnd   = parse.Any(parse.NewLine, parse.EOF[string]())
)

type goSingleLineCommentParser struct {
	ParseUntilDoubleCloseBrace bool
}

var (
	goSingleLineComment      = goSingleLineCommentParser{}
	goTemplSingleLineComment = goSingleLineCommentParser{
		ParseUntilDoubleCloseBrace: true,
	}
)

func (p goSingleLineCommentParser) Parse(pi *parse.Input) (n Node, ok bool, err error) {
	// Comment start.
	var c GoComment
	if _, ok, err = goSingleLineCommentStart.Parse(pi); err != nil || !ok {
		return
	}
	// Once we've got the comment start sequence, parse anything until the end
	// sequence as the comment contents.

	if !p.ParseUntilDoubleCloseBrace {
		if c.Contents, ok, err = parse.StringUntil(goSingleLineCommentEnd).Parse(pi); err != nil || !ok {
			err = parse.Error("expected end comment literal '\n' not found", pi.Position())
			return
		}
	} else {
		// try to close if its a line comment in gotempl comment
		_, _, _ = parse.StringUntil(dblCloseBraceWithOptionalPadding).Parse(pi)
	}
	// Move past the end element.
	_, _, _ = goSingleLineCommentEnd.Parse(pi)
	// Return the comment.
	c.Multiline = false
	return c, true, nil
}

var (
	goMultiLineCommentStart = parse.String("/*")
	goMultiLineCommentEnd   = parse.String("*/")
)

type goMultiLineCommentParser struct{}

var goMultiLineComment = goMultiLineCommentParser{}

func (p goMultiLineCommentParser) Parse(pi *parse.Input) (n Node, ok bool, err error) {
	// Comment start.
	start := pi.Position()
	var c GoComment
	if _, ok, err = goMultiLineCommentStart.Parse(pi); err != nil || !ok {
		return
	}

	// Once we've got the comment start sequence, parse anything until the end
	// sequence as the comment contents.
	if c.Contents, ok, err = parse.StringUntil(goMultiLineCommentEnd).Parse(pi); err != nil || !ok {
		err = parse.Error("expected end comment literal '*/' not found", start)
		return
	}
	// Move past the end element.
	_, _, _ = goMultiLineCommentEnd.Parse(pi)
	// Return the comment.
	c.Multiline = true
	return c, true, nil
}

var (
	goComment      = parse.Any[Node](goSingleLineComment, goMultiLineComment)
	goTemplComment = parse.Any[Node](goTemplSingleLineComment, goMultiLineComment)
)
