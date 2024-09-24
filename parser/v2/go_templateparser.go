package parser

import (
	"fmt"

	"github.com/a-h/parse"
)

/**
 *
 *
 * Adhoc Go parser.
 *
 *
 */

// Go Template

var gotemplate = parse.Func(func(pi *parse.Input) (r GoTemplate, ok bool, err error) {
	start := pi.Position()

	// gotempl FuncName(p Person, other Other) {
	var te gotemplateExpression
	if te, ok, err = gotemplateExpressionParser.Parse(pi); err != nil || !ok {
		return
	}
	r.Expression = te.Expression

	// Once we're in a gotemplate, we should expect some gotemplate whitespace, if/switch/for,
	// or node string expressions etc.
	var nodes Nodes
	nodes, ok, err = newGoTemplateNodeParser(lastCloseBraceWithOptionalPadding, "gotemplate closing brace").Parse(pi)
	if err != nil {
		return
	}
	if !ok {
		err = parse.Error("gotempl: expected nodes in gotempl body, but found none", pi.Position())
		return
	}
	r.Children = nodes.Nodes

	// Eat any whitespace.
	_, _, err = parse.OptionalWhitespace.Parse(pi)
	if err != nil {
		return
	}

	// Try for }
	if _, ok, err = closeBraceWithOptionalPadding.Parse(pi); err != nil || !ok {
		err = parse.Error("gotemplate: missing closing brace", pi.Position())
		return
	}

	r.Range = NewRange(start, pi.Position())

	return r, true, nil
})

/**
 *
 *
 *
 *
 *
 *
 *
 *
 *
 */

// TemplateExpression.

// TemplateExpression.
// gotempl Func(p Parameter) {
// gotempl (data Data) Func(p Parameter) {
// gotempl (data []string) Func(p Parameter) {
type gotemplateExpression struct {
	Expression Expression
}

var gotemplateExpressionParser = parse.Func(func(pi *parse.Input) (r gotemplateExpression, ok bool, err error) {
	start := pi.Index()

	if !peekPrefix(pi, "gotempl ") {
		return r, false, nil
	}

	// Once we have the prefix, everything to the brace is Go.
	// e.g.
	// gotempl (x []string) Test() {
	// becomes:
	// func (x []string) Test() templ.Component {
	if _, r.Expression, err = parseGoTemplFuncDecl(pi); err != nil {
		return r, false, err
	}

	// Eat " {\n".
	if _, ok, err = parse.All(openBraceWithOptionalPadding, parse.StringFrom(parse.Optional(parse.NewLine))).Parse(pi); err != nil || !ok {
		err = parse.Error("gotempl: malformed gotempl expression, expected `gotempl functionName() {`", pi.PositionAt(start))
		return
	}

	return r, true, nil
})

const (
	gotemplUnterminatedMissingEnd = `missing end (expected '{{end}}') - https://templ.guide/syntax-and-usage/statements#incomplete-statements`
)

// Template node (element, call, if, switch, for, whitespace etc.)
func newGoTemplateNodeParser[TUntil any](until parse.Parser[TUntil], untilName string) gotemplateNodeParser[TUntil] {
	return gotemplateNodeParser[TUntil]{
		until:     until,
		untilName: untilName,
	}
}

type gotemplateNodeParser[TUntil any] struct {
	until     parse.Parser[TUntil]
	untilName string
}

var gotemplateNodeSkipParsers = []parse.Parser[Node]{}

// {{ end }}
var goTemplExpressionEnd = parse.All(
	parse.OptionalWhitespace,
	parse.String("{{"),
	parse.OptionalWhitespace,
	parse.String("end"),
	parse.OptionalWhitespace,
	parse.String("}}"),
)

var gotemplateNodeParsers = [...]parse.Parser[Node]{
	// goComment, part of output in gotempl
	// goifExpression,           // TODO:
	goForExpression, // {{for ...}}...{{end}}
	// goswitchExpression,       // maybe not worth it
	gostringExpression,     // %{ "abc" }%, %{ fmt.Sprintf("abc") }%
	callTemplateExpression, // {! TemplateName(a, b, c) }
	templElementExpression, // @TemplateName(a, b, c) { <div>Children</div> }
	childrenExpression,     // { children... }
	goCode,                 // {{ myval := x.myval }}
	gowhitespaceExpression,
	gotextParser, // match anything, assume they're valid go code fragments
}

func (p gotemplateNodeParser[T]) Parse(pi *parse.Input) (op Nodes, ok bool, err error) {
outer:
	for {
		// Check if we've reached the end.
		if p.until != nil {
			start := pi.Index()
			_, ok, err = p.until.Parse(pi)
			if err != nil {
				return
			}
			if ok {
				// end reached for a gotempl ...(...) {}
				pi.Seek(start)
				return op, true, nil
			}
		}

		// Skip any nodes that we don't care about.
		for _, p := range gotemplateNodeSkipParsers {
			_, matched, err := p.Parse(pi)
			if err != nil {
				return Nodes{}, false, err
			}
			if matched {
				continue outer
			}
		}

		// Attempt to parse a node.
		// Loop through the parsers and try to parse a node.
		var matched bool
		for _, p := range gotemplateNodeParsers {
			var node Node
			node, matched, err = p.Parse(pi)
			if err != nil {
				return Nodes{}, false, err
			}
			if matched {
				op.Nodes = append(op.Nodes, node)
				break
			}
		}
		if matched {
			continue
		}

		if p.until == nil {
			// In this case, we're just reading as many nodes as we can until we can't read any more.
			// If we've reached here, we couldn't find a node.
			// The element parser checks the final node returned to make sure it's the expected close tag.
			break
		}

		err = UntilNotFoundError{
			ParseError: parse.Error(fmt.Sprintf("%v not found", p.untilName), pi.Position()),
		}
		return
	}

	return op, true, nil
}
