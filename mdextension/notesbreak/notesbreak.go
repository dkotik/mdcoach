/*
Package notesbreak provides Goldmark Markdown extension for detecting and rendering document end notes.
*/
package notesbreak

import (
	"github.com/yuin/goldmark/ast"
)

var KindNotesBreak = ast.NewNodeKind("NotesBreak")

type NotesBreak struct {
	*ast.ThematicBreak
}

func (n *NotesBreak) Kind() ast.NodeKind {
	return KindNotesBreak
}
