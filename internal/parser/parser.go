package parser

import (
	"bufio"
	"strings"
	"unicode"

	"github.com/itchyny/gojq"
	"github.com/pkg/errors"
)

// ParseExpectations parses an "expect" section. The expect section consists of several jq filters, one per line.
// In order to allow longer filter expressions, a filter expression may be continued on the next line. This is indicated
// by having the continuation line start with any whitespace character.
func ParseExpectations(spec string, sctx SourceContext) ([]*ParsedQuery, error) {
	var queries []*ParsedQuery
	scanner := bufio.NewScanner(strings.NewReader(spec))
	current := ""
	scanned := true
	if sctx.IsZero() {
		sctx = SourceContext{
			Filename: "<expectations block>",
			Line:     0,
		}
	}
	currentSCtx := sctx
	for ; scanned; sctx.Line++ {
		scanned = scanner.Scan()
		var next string
		if scanned {
			line := strings.TrimRightFunc(scanner.Text(), unicode.IsSpace)
			trimmed := strings.TrimLeftFunc(line, unicode.IsSpace)
			if len(trimmed) < len(line) {
				// Continuation line.
				if current == "" {
					return nil, errors.Errorf("unexpected continuation at %s", sctx)
				}
				current += " " + trimmed
				continue
			}
			next = line
		}

		if current != "" && !strings.HasPrefix(current, "#") {
			query, err := ParseQuery(current, currentSCtx)
			if err != nil {
				return nil, errors.Wrapf(err, "parsing query ending at %s", sctx)
			}
			queries = append(queries, query)
		}
		current = next
		currentSCtx = sctx
	}
	if err := scanner.Err(); err != nil {
		return nil, errors.Wrap(err, "parsing expectations")
	}

	if current != "" && !strings.HasPrefix(current, "#") {
		query, err := ParseQuery(current, currentSCtx)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing query ending at %s", sctx)
		}
		queries = append(queries, query)
	}

	return queries, nil
}

// ParseQuery parses a single query.
func ParseQuery(src string, sctx SourceContext) (*ParsedQuery, error) {
	query, err := gojq.Parse(src)
	if err != nil {
		return nil, err
	}
	if err := postProcessQuery(query); err != nil {
		return nil, err
	}
	return &ParsedQuery{
		Query:     query,
		Source:    src,
		SourceCtx: sctx,
	}, nil
}
