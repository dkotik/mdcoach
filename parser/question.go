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

func IsQuestion(node ast.Node, reader text.Reader) bool {
	if node.Kind() != ast.KindListItem {
		return false
	}
	firstParagraph := node.FirstChild()
	if firstParagraph == nil || firstParagraph.Kind() != ast.KindParagraph {
		return false
	}
	lines := firstParagraph.Lines()
	lineCount := lines.Len()
	if lineCount < 1 {
		return false
	}
	lastLine := reader.Value(lines.At(lineCount - 1))
	if !reIsLineQuestion.Match(lastLine) {
		return false // no match
	}
	return true
}

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
	if !IsQuestion(node, reader) {
		return
	}
	node.SetAttributeString("question", true)
	questions := QuestionListFromContext(pc)
	pc.Set(questionsContextKey, append(
		questions,
		node,
	))
}

type QuestionExtractor func(pc parser.Context, question ast.Node)

type questionListItemExtractor struct {
	parser.BlockParser
	extractor QuestionExtractor
}

func NewReviewQuestionExtractor(extractor QuestionExtractor) parser.BlockParser {
	return &questionListItemExtractor{
		BlockParser: parser.NewListItemParser(),
		extractor:   extractor,
	}
}

func (p *questionListItemExtractor) Close(
	node ast.Node,
	reader text.Reader,
	pc parser.Context,
) {
	if !IsQuestion(node, reader) {
		return
	}
	p.extractor(pc, node)
}
