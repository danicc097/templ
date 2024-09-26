package parser

import (
	"io"

	"github.com/a-h/parse"
	"github.com/a-h/templ/parser/v2/goexpression"
)

// {{ if p.Type == "test" && p.thing }}
//
//	println(var)
//
// {{ end }}
type GoIfExpression struct {
	Expression Expression
	Then       []Node
	ElseIfs    []GoElseIfExpression
	Else       []Node
}

type GoElseIfExpression struct {
	Expression Expression
	Then       []Node
}

func (n GoIfExpression) ChildNodes() []Node {
	var nodes []Node
	nodes = append(nodes, n.Then...)
	nodes = append(nodes, n.Else...)
	for _, elseIf := range n.ElseIfs {
		nodes = append(nodes, elseIf.Then...)
	}
	return nodes
}

func (ie GoIfExpression) IsNode() bool { return true }

func (ie GoIfExpression) Write(w io.Writer, indent int) error {
	if err := writeIndent(w, indent, "{{ if ", ie.Expression.Value, " }}\n"); err != nil {
		return err
	}
	if err := writeNodesIndented(w, indent+1, ie.Then); err != nil {
		return err
	}
	for _, elif := range ie.ElseIfs {
		if err := elif.Write(w, indent); err != nil {
			return err
		}
	}
	if len(ie.Else) > 0 {
		if err := writeIndent(w, indent, "{{ else }}\n"); err != nil {
			return err
		}
		if err := writeNodesIndented(w, indent+1, ie.Else); err != nil {
			return err
		}
	}
	if err := writeIndent(w, indent, "{{ end }}"); err != nil {
		return err
	}
	return nil
}

// until {{ else if, {{ else or {{ end }}
var (
	gountilElseIfElseOrEnd = parse.Any(StripType(goelseIfExpression), StripType(goelseExpression), StripType(goTemplExpressionEnd))
	goIfBlockEnd           = parse.All(parse.OptionalWhitespace, parse.String("}}"), parse.NewLine)
	gountilIfBlockEnd      = parse.StringUntil(goIfBlockEnd)
)

var goIfExpression parse.Parser[Node] = goIfExpressionParser{}

type goIfExpressionParser struct{}

func (goIfExpressionParser) Parse(pi *parse.Input) (n Node, ok bool, err error) {
	var r GoIfExpression
	start := pi.Index()

	if _, ok, err = parse.All(
		parse.OptionalWhitespace,
		parse.String("{{"),
		parse.OptionalWhitespace,
	).Parse(pi); err != nil || !ok {
		pi.Seek(start)
		return r, false, nil
	}
	if !peekPrefix(pi, "if ") {
		pi.Seek(start)
		return r, false, nil
	}

	if r.Expression, err = parseGotemplGo("if", pi, goexpression.If); err != nil {
		return r, false, err
	}
	r.Expression.GoTempl = true

	// now eat up to first }}\n after if cond in the actual input
	if _, ok, err = goIfBlockEnd.Parse(pi); err != nil || !ok {
		err = parse.Error(`if: expected closing "}}" but was not found`, pi.Position())
		return r, false, err
	}

	// we may also have nested {{ if ... }} ... {{ end }} and we want to stop at the first elseif/else/end found after parsing those nodes
	np := newGoTemplateNodeParser(gountilElseIfElseOrEnd, "else expression or closing {{end}}")
	var thenNodes Nodes
	if thenNodes, ok, err = np.Parse(pi); err != nil || !ok {
		err = parse.Error("if: expected nodes, but none were found", pi.Position())
		return r, false, err
	}
	r.Then = thenNodes.Nodes

	if r.ElseIfs, _, err = parse.ZeroOrMore(goelseIfExpression).Parse(pi); err != nil {
		return r, false, err
	}

	var elseNodes Nodes
	if elseNodes, _, err = goelseExpression.Parse(pi); err != nil {
		return r, false, err
	}
	r.Else = elseNodes.Nodes

	if _, ok, err = goTemplExpressionEnd.Parse(pi); err != nil {
		err = parse.Error("if: missing `{{ end }}`", pi.Position())
		return r, false, err
	}

	return r, true, nil
}

var goelseIfExpression parse.Parser[GoElseIfExpression] = goelseIfExpressionParser{}

type goelseIfExpressionParser struct{}

func (goelseIfExpressionParser) Parse(pi *parse.Input) (r GoElseIfExpression, ok bool, err error) {
	start := pi.Index()

	if _, ok, err = parse.All(
		parse.OptionalWhitespace,
		parse.String("{{"),
		parse.OptionalWhitespace,
	).Parse(pi); err != nil || !ok {
		pi.Seek(start)
		return r, false, nil
	}

	if !peekPrefix(pi, "else if ") {
		pi.Seek(start)
		return r, false, nil
	}

	if r.Expression, err = parseGotemplGo("else if", pi, goexpression.If); err != nil {
		return r, false, err
	}
	r.Expression.GoTempl = true

	if _, ok, err = goIfBlockEnd.Parse(pi); err != nil || !ok {
		err = parse.Error(`else if: expected closing "}}" but was not found`, pi.Position())
		return r, false, err
	}

	// same possibility as if
	np := newGoTemplateNodeParser(gountilElseIfElseOrEnd, "else expression or closing brace")
	var thenNodes Nodes
	if thenNodes, ok, err = np.Parse(pi); err != nil || !ok {
		err = parse.Error("else if: expected nodes, but none were found", pi.Position())
		return r, false, err
	}
	r.Then = thenNodes.Nodes

	return r, true, nil
}

func (e GoElseIfExpression) Write(w io.Writer, indent int) error {
	if err := writeIndent(w, indent, "{{ else if ", e.Expression.Value, " }}\n"); err != nil {
		return err
	}
	return writeNodesIndented(w, indent+1, e.Then)
}

var goelseExpression parse.Parser[Nodes] = goelseExpressionParser{}

type goelseExpressionParser struct{}

func (goelseExpressionParser) Parse(pi *parse.Input) (r Nodes, ok bool, err error) {
	start := pi.Index()

	if _, ok, err = goTemplExpressionElse.Parse(pi); err != nil || !ok {
		pi.Seek(start)
		return r, false, nil
	}

	if r, ok, err = newGoTemplateNodeParser(goTemplExpressionEnd, "else expression closing {{end}}").Parse(pi); err != nil || !ok {
		pi.Seek(start)
		return r, false, err
	}

	return r, true, nil
}
