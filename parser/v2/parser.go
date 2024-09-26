package parser

import (
	"unicode"

	"github.com/a-h/parse"
)

// adapted from go/src/unicode/tables.go
var _White_Space_No_Newline = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x0009, 0x0009, 1},    // Horizontal Tab (U+0009)
		{0x0020, 0x0020, 1},    // Space (U+0020)
		{0x00A0, 0x1680, 5600}, // No-Break Space (U+00A0) to Ogham Space Mark (U+1680)
		{0x2000, 0x200A, 1},    // Various spaces (U+2000 to U+200A)
		{0x202F, 0x205F, 48},   // Narrow No-Break Space (U+202F) to Medium Mathematical Space (U+205F)
		{0x3000, 0x3000, 1},    // Ideographic Space (U+3000)
	},
	LatinOffset: 2,
}

var whitespaceExceptNewline parse.Parser[string] = parse.StringFrom(parse.OneOrMore(parse.RuneInRanges(_White_Space_No_Newline)))

// ) {
var expressionFuncEnd = parse.All(parse.Rune(')'), openBraceWithOptionalPadding)

// Template

var template = parse.Func(func(pi *parse.Input) (r HTMLTemplate, ok bool, err error) {
	start := pi.Position()

	// templ FuncName(p Person, other Other) {
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
