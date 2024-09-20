package testgousage

import (
	"context"
	_ "embed"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

//go:embed expected.go
var expected string

func Test(t *testing.T) {
	component := TestComponent()
	component.Render(context.Background(), os.Stdout)
	// TODO: use gofmt printer on both and cmp diff

	// use cmp diff to compare:
	// cmp.Diff(component, expected)

	if diff := cmp.Diff(component, expected); diff != "" {
		t.Error(diff)
	}
}
