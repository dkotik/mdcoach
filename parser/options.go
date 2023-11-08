package parser

import (
	"errors"

	"github.com/yuin/goldmark/parser"
)

type options struct {
	listItemParser parser.BlockParser
}
type Option func(*options) error

func WithListItemParser(p parser.BlockParser) Option {
	return func(o *options) error {
		if p == nil {
			return errors.New("cannot use a <nil> list item parser")
		}
		if o.listItemParser != nil {
			return errors.New("list item parser is already set")
		}
		o.listItemParser = p
		return nil
	}
}

func WithQuestionExtractor(extractor QuestionExtractor) Option {
	return func(o *options) error {
		if extractor == nil {
			return errors.New("cannot use a <nil> question extractor")
		}
		return WithListItemParser(NewReviewQuestionExtractor(extractor))(o)
	}
}

func withDefaultListItemParser() Option {
	return func(o *options) error {
		if o.listItemParser != nil {
			return nil
		}
		// util.Prioritizparser.Nd(NewListItemParser(), 400),
		// util.Prioritized(NewReviewQuestionParser(), 400),
		return WithListItemParser(NewReviewQuestionParser())(o)
	}
}
