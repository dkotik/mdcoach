package mdcoach

import (
	"fmt"

	"github.com/OneOfOne/xxhash"
	"github.com/dkotik/mdcoach/internal"
	"github.com/yuin/goldmark/v2/ast"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/text"
)

type documentIDInjector struct{}

// NewDocumentIDInjector returns a transformer that sets the document's id
// attribute to the xxhash of its deterministic AST dump.
func NewDocumentIDInjector() parser.ASTTransformer {
	return &documentIDInjector{}
}

// Transform hashes the AST dump and stores the hexadecimal hash in the
// document's id attribute.
func (*documentIDInjector) Transform(document *ast.Document, reader text.Reader, _ parser.Context) {
	var source []byte
	if reader != nil {
		source = reader.Source()
	}

	hash := xxhash.New64()
	internal.WriteAST(hash, document, source)
	document.SetAttribute(
		"id",
		text.NewMultiLineValueFromString(fmt.Sprintf("%x", hash.Sum(nil)), text.IdentityDecoder),
	)
}

// propertyValue formats a Goldmark node property for callers that need to
// extend the AST dump with deterministic property output.
func propertyValue(value any, source []byte) string {
	switch value := value.(type) {
	case text.SingleLineValue:
		return value.Value(source)
	case text.MultiLineValue:
		return value.Value(source)
	default:
		return fmt.Sprint(value)
	}
}
