/*
Package figure renders Goldmark Markdown paragraphs with only one child, which is an image,
as HTML figure tag with caption set to image title.
*/
package figure

import (
	"github.com/yuin/goldmark/ast"
)

// HasOnlyOneChildOfKind returns true if a given [ast.Node] holds only one child,
// which matches exactly the expected [ast.NodeKind].
func HasOnlyOneChildOfKind(n ast.Node, k ast.NodeKind) bool {
	if n.ChildCount() != 1 {
		return false
	}
	return n.FirstChild().Kind() == k
}
