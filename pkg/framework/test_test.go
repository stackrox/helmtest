package framework

import (
	"github.com/stretchr/testify/assert"
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
		},
	}

	rootQuery := []string{"root test"}
	testCases := map[string]struct{
	query            []string
	expectNotFound   bool
}{
		"with only root node": {query: rootQuery},
		"with child test": {query: append(rootQuery, "child test")},
		"with nested child": {query: append(rootQuery, "child test", "child child test")},
		"with not existing nested": {query: append(rootQuery, "child test", "child child test", "does not exist"), expectNotFound: true},
		"with not existing root": {query: []string{"root does not exist"}, expectNotFound: true},
		"with another child": {query: append(rootQuery, "another child test")},
		"with another nested child": {query: append(rootQuery, "another child test", "another child child test")},
	}

	for _, tt := range testCases {
		r := suite.find(tt.query)
		if tt.expectNotFound {
			assert.Nil(t, r)
		} else {
			assert.Equal(t, tt.query[len(tt.query)-1], r.Name)
		}
	}
}
