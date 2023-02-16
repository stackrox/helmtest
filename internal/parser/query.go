package parser

import "github.com/itchyny/gojq"

// ParsedQuery is a parsed query, with some extra metadata to aid diagnosing test failures.
type ParsedQuery struct {
	*gojq.Query
	Source    string
	SourceCtx SourceContext
}

// Copy returns deep copy of ParsedQuery
func (q *ParsedQuery) Copy() *ParsedQuery {
	query, _ := gojq.Parse(q.String())
	return &ParsedQuery{
		Query:     query,
		Source:    q.Source,
		SourceCtx: q.SourceCtx,
	}
}
