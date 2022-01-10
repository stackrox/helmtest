package framework

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFind(t *testing.T) {
	childTest := &Test{
		Name: "child test",
		Tests: []*Test{{
			Name: "child child test",
		}},
	}
	anotherChildTest := &Test{
		Name: "another child test",
		Tests: []*Test{{
			Name: "another child child test",
		}},
	}

	suite := &Test{
		Name: "root test",
		Tests: []*Test{
			childTest,
			anotherChildTest,
			{Name: "same name"},
			{Name: "same name"},
		},
	}

	rootQuery := []string{"root test"}
	testCases := map[string]struct {
		query          []string
		expectNotFound bool
	}{
		"with only root node":       {query: rootQuery},
		"with child test":           {query: []string{"root test", "child test"}},
		"with nested child":         {query: []string{"root test", "child test", "child child test"}},
		"with not existing nested":  {query: []string{"root test", "child test", "child child test", "does not exist"}, expectNotFound: true},
		"with not existing root":    {query: []string{"root does not exist"}, expectNotFound: true},
		"with another child":        {query: []string{"root test", "another child test"}},
		"with another nested child": {query: []string{"root test", "another child test", "another child child test"}},
		"with same name finds both": {query: []string{"root test", "same name"}},
	}

	for _, tt := range testCases {
		results := suite.find(tt.query)
		if tt.expectNotFound {
			assert.Nil(t, results)
		} else {
			require.NotEmpty(t, results)
			for _, result := range results {
				assert.Equal(t, tt.query[len(tt.query)-1], result.Name)
			}
		}
	}
}
