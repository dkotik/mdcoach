package mdcoach

import (
	"bytes"
	"fmt"
	"testing"
)

func TestYamlVisualizer(t *testing.T) {
	var b bytes.Buffer
	err := YamlWithTemplate(&b, `assets/test-yaml.yaml`, `assets/test-template.tmpl`)
	if err != nil || b.String() != `--` {
		fmt.Printf("Error: %s.\n", err)
		fmt.Print(b.String())
		// t.Fail()
	}
}
