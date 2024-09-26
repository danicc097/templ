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

// FIXME: should respect indent for Text nodes inside gotempl.
// we will apply gofumpt to output anyway
func main() {
	const z = "z"
	for i := 0; i < 3; i++ {
		println("z")
	}
	for i := 0; i < 3; i++ {
		println("z")
	}
}
