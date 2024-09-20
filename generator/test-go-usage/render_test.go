package testgousage

import (
	"bytes"
	"context"
	_ "embed"
	"testing"

	"github.com/google/go-cmp/cmp"
)

//go:embed expected.go
var expected string

func Test(t *testing.T) {
	component := TestComponent()
	buf := bytes.Buffer{}
	component.Render(context.Background(), &buf)
	// TODO: use gofmt printer on both and cmp diff
	if diff := cmp.Diff(buf.String(), expected); diff != "" {
		t.Error(diff)
	}
}
