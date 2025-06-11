package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseExpectations(t *testing.T) {
	tests := map[string]struct {
		spec    string
		want    []string
		errFunc assert.ErrorAssertionFunc
	}{
		"empty": {spec: "", want: nil, errFunc: assert.NoError},
		"one query": {
			spec:    "foo | bar",
			want:    []string{"foo | bar"},
			errFunc: assert.NoError,
		},
		"a couple of queries with comments": {
			spec: `# a comment
query1 | very |
  long
# comment in between
query2 | even |
  longer
# trailing comment`,
			want: []string{
				"query1 | very | long",
				"query2 | even | longer",
			},
			errFunc: assert.NoError,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			sourceCtx := SourceContext{
				Filename: name,
				Line:     1,
			}
			got, err := ParseExpectations(tt.spec, sourceCtx)
			tt.errFunc(t, err)
			assert.Len(t, got, len(tt.want))
			for i, query := range got {
				assert.Equal(t, tt.want[i], query.Source, "source of query %d should match", i)
			}
		})
	}
}
