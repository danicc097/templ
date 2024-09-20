package parser

import (
	"github.com/a-h/parse"
)

// ) {
// TODO: now must find {{end}}
var goExpressionFuncEnd = parse.All(parse.Rune(')'), openBraceWithOptionalPadding)

// Template

var goTemplateParser = parse.Func(func(pi *parse.Input) (r GoTemplate, ok bool, err error) {
	start := pi.Position()

	// gotempl FuncName(...) {
	var te templateExpression
	if te, ok, err = templateExpressionParser.Parse(pi); err != nil || !ok {
		return
	}
	r.Expression = te.Expression

	// Once we're in a template, we should expect some template whitespace, if/switch/for,
	// or node string expressions etc.
	var nodes Nodes
	nodes, ok, err = newTemplateNodeParser(closeBraceWithOptionalPadding, "template closing brace").Parse(pi)
	if err != nil {
		return
	}

	if !ok {
		err = parse.Error("templ: expected nodes in templ body, but found none", pi.Position())
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
		err = parse.Error("template: missing closing brace", pi.Position())
		return
	}

	r.Range = NewRange(start, pi.Position())

	return r, true, nil
})

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
		err = parse.Error("malformed gotempl expression, expected `gotempl functionName() {`", pi.PositionAt(start))
		return
	}

	return r, true, nil
})
