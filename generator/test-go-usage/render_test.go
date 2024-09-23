package testgousage

import (
	"bytes"
	"context"
	_ "embed"
	goformat "go/format"
	"testing"

	"mvdan.cc/gofumpt/format"

	"github.com/google/go-cmp/cmp"
)

//go:embed expected.go
var expected string

func Test(t *testing.T) {
	component := TestComponent()
	buf := bytes.Buffer{}
	component.Render(context.Background(), &buf)
	src, err := goformat.Source(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	srcFormatted, err := format.Source(src, format.Options{})
	if err != nil {
		t.Fatal(err)
	}
	expectedFormatted, err := format.Source([]byte(expected), format.Options{})
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(string(expectedFormatted), string(srcFormatted)); diff != "" {
		t.Errorf(" mismatch (-want +got):\n%s", diff)
	}
}
