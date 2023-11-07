package test

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/dkotik/mdcoach/document"
	mdcoachParser "github.com/dkotik/mdcoach/parser"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

func TestMetaDataLoading(t *testing.T) {
	ctx := parser.NewContext()
	_ = mdcoachParser.New().Parse(text.NewReader([]byte(`---
title: Testing Front Matter
author: Anonymous
year: 2023
tags: [markdown, goldmark]
description: |
  Testing parsing YAML front matter. TOML is also supported.
---

# Heading 1`)), parser.WithContext(ctx))

	metadata, err := document.NewMetadata(ctx)
	if err != nil {
		t.Fatal(err)
	}
	spew.Dump(metadata)
	// t.Fatal("impl")
}
