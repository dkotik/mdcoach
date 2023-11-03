package picture

import "testing"

func newTestLocalProvider(t *testing.T) *LocalProvider {
	p, err := NewLocalProvider(
		WithDestinationPath(t.TempDir()),
	)
	if err != nil {
		t.Fatal(err)
	}
	return p
}
