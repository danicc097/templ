/*
Package testgousage.
*/
package testgousage

// A comment that should be rendered
/* Another comment that should be rendered
 */
// comment with var with space
// comment withvarwithout space
const _ab = "ab"

// "ab"-1-1-y
// nested 1
// nested 2
// nested 3
// "ab"-2-1-y
// nested 1
// nested 2
// nested 3
// "ab"-3-1-y
// nested 1
// nested 2
// nested 3
var abc = "abc"

func main() {
	const z = "z"
	for i := 0; i < 3; i++ {
		println("z")
	}
	for i := 0; i < 3; i++ {
		println("z")
	}
	_ = map[string][]A{
		"a": {{}, {}, {}},
}
}

/*
 * Fields: a, b, c || a, b, c
 */
// Fields: a, b, c
type A struct {
	a string
	b string
	c string
}
