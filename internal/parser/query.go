package parser

import "github.com/itchyny/gojq"

// ParsedQuery is a parsed query, with some extra metadata to aid diagnosing test failures.
type ParsedQuery struct {
	*gojq.Query
	Source    string
	SourceCtx SourceContext
}
