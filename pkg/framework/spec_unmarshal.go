package framework

import (
	"fmt"
	"github.com/stackrox/helmtest/internal/parser"
	yamlv3 "gopkg.in/yaml.v3"
	"strings"
)

const (
	filenamePragma = `#!helmtest-filename:`
)

func injectFilename(node *yamlv3.Node, filename string) {
	node.FootComment += fmt.Sprintf("\n%s%s", filenamePragma, filename)
	for _, child := range node.Content {
		injectFilename(child, filename)
	}
}

func (t *Test) UnmarshalYAML(node *yamlv3.Node) error {
	type testNoMethods Test
	if err := node.Decode((*testNoMethods)(t)); err != nil {
		return err
	}

	if node.Kind != yamlv3.MappingNode {
		return nil // weird but ok
	}

	// We've moved the filename into a fake footer comment, to allow extracting it here.
	var filename string
	for _, line := range strings.Split(node.FootComment, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, filenamePragma) {
			filename = line[len(filenamePragma):]
			break
		}
	}

	for i := 0; i < len(node.Content); i += 2 {
		valueNode := node.Content[i+1]
		valueSrcContext := parser.SourceContext{
			Filename: filename,
			Line:     valueNode.Line,
		}
		// In literal and folded style, the value only begins on the preceding line.
		if valueNode.Style == yamlv3.LiteralStyle || valueNode.Style == yamlv3.FoldedStyle {
			valueSrcContext.Line++
		}

		keyNode := node.Content[i]
		if keyNode.Value == "defs" {
			t.defsSrcCtx = valueSrcContext
		} else if keyNode.Value == "expect" {
			t.expectSrcCtx = valueSrcContext
		}
	}
	return nil
}
