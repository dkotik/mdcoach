package mdcoach

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/dkotik/mdcoach/internal"
	"github.com/sebdah/goldie/v2"
)

func TestStrikethroughExtension(t *testing.T) {
	source := []byte("This is ~~struck~~ text.\n")
	tree := NewParser().Parse(source)

	var output bytes.Buffer
	if err := NewRenderer(nil).Render(&output, source, tree); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "<del>struck</del>") {
		t.Fatalf("rendered HTML does not contain strikethrough: %s", output.String())
	}
}

func TestTableExtension(t *testing.T) {
	source := []byte("| Name | Score |\n| --- | ---: |\n| Ada | 10 |\n")
	tree := NewParser().Parse(source)

	var output bytes.Buffer
	if err := NewRenderer(nil).Render(&output, source, tree); err != nil {
		t.Fatal(err)
	}

	rendered := output.String()
	for _, expected := range []string{"<table>", "<th>Name</th>", "<td>Ada</td>", `<td style="text-align:right">10</td>`} {
		if !strings.Contains(rendered, expected) {
			t.Errorf("rendered HTML does not contain %q: %s", expected, rendered)
		}
	}
}

func TestTaskListExtension(t *testing.T) {
	source := []byte("- [x] Completed task\n- [ ] Pending task\n")
	tree := NewParser().Parse(source)

	var output bytes.Buffer
	if err := NewRenderer(nil).Render(&output, source, tree); err != nil {
		t.Fatal(err)
	}

	rendered := output.String()
	for _, expected := range []string{
		`<input checked="" disabled="" type="checkbox">`,
		`<input disabled="" type="checkbox">`,
		"Completed task",
		"Pending task",
	} {
		if !strings.Contains(rendered, expected) {
			t.Errorf("rendered HTML does not contain %q: %s", expected, rendered)
		}
	}
}

func TestPresentationAST(t *testing.T) {
	source, err := os.ReadFile("testdata/presentation-1.md")
	if err != nil {
		t.Fatal(err)
	}

	tree := NewParser().Parse(source)
	var dump bytes.Buffer
	internal.WriteAST(&dump, tree, source)

	for child := range tree.Children() {
		switch child.Kind() {
		case SlideKind:
			for _, attr := range child.Attributes() {
				switch attr.Name {
				case "class":
					t.Error("slide contains class attribute:", attr.Value.Str(source))
				}
			}
		}
	}

	goldie.New(t).Assert(t, "presentation_ast", dump.Bytes())
}
