package parser

import (
	"io"

	"github.com/a-h/parse"
	"github.com/a-h/templ/parser/v2/goexpression"
)

// {{ for i, v := range p.Addresses }}
//
//	{! Address(v) }
//
// {{ end }}
type GoForExpression struct {
	Expression Expression
	Children   []Node
}

func (fe GoForExpression) ChildNodes() []Node {
	return fe.Children
}
func (fe GoForExpression) IsNode() bool { return true }
func (fe GoForExpression) Write(w io.Writer, indent int) error {
	if err := writeIndent(w, indent, " {{ for ", fe.Expression.Value, " }}\n"); err != nil {
		return err
	}
	if err := writeNodesIndented(w, indent+1, fe.Children); err != nil {
		return err
	}
	if err := writeIndent(w, indent, "{{ end }}"); err != nil {
		return err
	}
	return nil
}

var goForExpression parse.Parser[Node] = goForExpressionParser{}

type goForExpressionParser struct{}

func (goForExpressionParser) Parse(pi *parse.Input) (n Node, ok bool, err error) {
	var r GoForExpression
	start := pi.Index()

	_, _, _ = parse.OptionalWhitespace.Parse(pi)

	// Detect the `{{ for ` syntax, allowing for optional spaces around `{{`.
	if _, ok, err = parse.All(
		parse.OptionalWhitespace,
		parse.String("{{"),
		parse.OptionalWhitespace,
	).Parse(pi); err != nil || !ok {
		pi.Seek(start)
		return r, false, nil
	}

	if !peekPrefix(pi, "for ") {
		// not a for loop
		pi.Seek(start)
		return r, false, nil
	}

	if r.Expression, err = parseGo("for", pi, goexpression.For); err != nil {
		return r, false, err
	}

	if _, ok, err = parse.All(parse.OptionalWhitespace, parse.String("}}")).Parse(pi); err != nil || !ok {
		err = parse.Error(`for: expected closing "}}" but was not found`, pi.Position())
		return r, false, err
	}

	// Parse the body of the `for` loop (everything until `{{ end }}`).
	// There may be other gotempl statements in between.
	tnp := newGoTemplateNodeParser(goTemplExpressionEnd, "for expression closing {{end}}")
	var nodes Nodes
	if nodes, ok, err = tnp.Parse(pi); err != nil || !ok {
		err = parse.Error("for: expected nodes, but none were found", pi.Position())
		return
	}
	r.Children = nodes.Nodes

	if _, ok, err = goTemplExpressionEnd.Parse(pi); err != nil || !ok {
		err = parse.Error("for: missing `{{ end }}`", pi.Position())
		return
	}

	return r, true, nil
}
