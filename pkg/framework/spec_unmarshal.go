package framework

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/stackrox/helmtest/internal/parser"
	yamlv3 "gopkg.in/yaml.v3"
	"sigs.k8s.io/yaml"
)

const (
	filenamePragma = `#!helmtest-filename:`
)

// injectFilename stores the filename into a "fake" foot comment of the respective node. This allows us to reconstruct
// the filename in UnmarshalYAML, even though it isn't stored in the node itself (and there is no context that would
// allow us to pass other values).
func injectFilename(node *yamlv3.Node, filename string) {
	node.FootComment += fmt.Sprintf("\n%s%s", filenamePragma, filename)
	for _, child := range node.Content {
		injectFilename(child, filename)
	}
}

// UnmarshalYAML unmarshals a test spec from YAML, preserving line information that we're interested in.
func (t *Test) UnmarshalYAML(node *yamlv3.Node) error {
	// Create an alias of this type without a custom UnmarshalYAML method.
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

	// Check for mapping keys of interest.
	for i := 0; i < len(node.Content); i += 2 {
		valueNode := node.Content[i+1]
		valueSrcContext := parser.SourceContext{
			Filename: filename,
			Line:     valueNode.Line - 1, // Node.Line is one-based, but SourceContext.Line is zero-based
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

// UnmarshalYAML unmarshals a RawDict from YAML, making sure that the resulting type matches the result of
// json.Unmarshal on the equivalent JSON.
func (d *RawDict) UnmarshalYAML(node *yamlv3.Node) error {
	yamlBytes, err := yamlv3.Marshal(node)
	if err != nil {
		return err
	}
	jsonBytes, err := yaml.YAMLToJSON(yamlBytes)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonBytes, d)
}
