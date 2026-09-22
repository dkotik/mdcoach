package presentation

import (
	"bytes"
	"context"
	"testing"

	"github.com/sebdah/goldie/v2"
)

func TestNewPresentation(t *testing.T) {
	var output bytes.Buffer
	if err := New(
		context.Background(),
		&output,
		[]string{"../testdata/presentation.md"},
	); err != nil {
		t.Fatal(err)
	}

	goldie.New(t, goldie.WithNameSuffix(".html")).Assert(t, "presentation", output.Bytes())
}
