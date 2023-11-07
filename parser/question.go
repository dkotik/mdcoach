package parser

import (
	"regexp"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

var (
	questionsContextKey = parser.NewContextKey()
	reIsLineQuestion    = regexp.MustCompile(`\w\?\s*$`)
)

func QuestionListFromContext(ctx parser.Context) []ast.Node {
	questions, _ := ctx.Get(questionsContextKey).([]ast.Node)
	return questions
}

type questionListItemParser struct {
	parser.BlockParser
}

func NewReviewQuestionParser() parser.BlockParser {
	return &questionListItemParser{
		BlockParser: parser.NewListItemParser(),
	}
}

func (p *questionListItemParser) Close(
	node ast.Node,
	reader text.Reader,
	pc parser.Context,
) {
	// detect question
	if node.Kind() != ast.KindListItem {
		return
	}
	firstParagraph := node.FirstChild()
	if firstParagraph == nil || firstParagraph.Kind() != ast.KindParagraph {
		return
	}
	lines := firstParagraph.Lines()
	lineCount := lines.Len()
	if lineCount < 1 {
		return
	}
	lastLine := reader.Value(lines.At(lineCount - 1))
	if !reIsLineQuestion.Match(lastLine) {
		return // no match
	}

	questions := QuestionListFromContext(pc)
	pc.Set(questionsContextKey, append(
		questions,
		node,
	))
	// spew.Dump(lastLine)
	// spew.Dump(node.ChildCount())

	// add question to context
}

/*
import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type listItemParser struct {
}

var defaultListItemParser = &listItemParser{}

// NewListItemParser returns a new BlockParser that
// parses list items.
func NewListItemParser() BlockParser {
	return defaultListItemParser
}

func (b *listItemParser) Trigger() []byte {
	return []byte{'-', '+', '*', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
}

func (b *listItemParser) Open(parent ast.Node, reader text.Reader, pc Context) (ast.Node, State) {
	list, lok := parent.(*ast.List)
	if !lok { // list item must be a child of a list
		return nil, NoChildren
	}
	offset := lastOffset(list)
	line, _ := reader.PeekLine()
	match, typ := matchesListItem(line, false)
	if typ == notList {
		return nil, NoChildren
	}
	if match[1]-offset > 3 {
		return nil, NoChildren
	}

	pc.Set(emptyListItemWithBlankLines, nil)

	itemOffset := calcListOffset(line, match)
	node := ast.NewListItem(match[3] + itemOffset)
	if match[4] < 0 || util.IsBlank(line[match[4]:match[5]]) {
		return node, NoChildren
	}

	pos, padding := util.IndentPosition(line[match[4]:], match[4], itemOffset)
	child := match[3] + pos
	reader.AdvanceAndSetPadding(child, padding)
	return node, HasChildren
}

func (b *listItemParser) Continue(node ast.Node, reader text.Reader, pc Context) State {
	line, _ := reader.PeekLine()
	if util.IsBlank(line) {
		reader.Advance(len(line) - 1)
		return Continue | HasChildren
	}

	offset := lastOffset(node.Parent())
	isEmpty := node.ChildCount() == 0
	indent, _ := util.IndentWidth(line, reader.LineOffset())
	if (isEmpty || indent < offset) && indent < 4 {
		_, typ := matchesListItem(line, true)
		// new list item found
		if typ != notList {
			pc.Set(skipListParserKey, listItemFlagValue)
			return Close
		}
		if !isEmpty {
			return Close
		}
	}
	pos, padding := util.IndentPosition(line, reader.LineOffset(), offset)
	reader.AdvanceAndSetPadding(pos, padding)

	return Continue | HasChildren
}
*/
