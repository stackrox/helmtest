package framework

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"text/template"

	"github.com/stackrox/helmtest/internal/parser"
	helmUtil "github.com/stackrox/helmtest/internal/rox-imported/helmutil"
	"github.com/stackrox/helmtest/internal/rox-imported/pointers"

	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/chartutil"
)

// applySetOptions takes the values specified in the `set` stanza and merges them into the otherwise defined values.
func (t *Test) applySetOptions() error {
	for keyPathStr, val := range t.Set {
		vals, err := helmUtil.ValuesForKVPair(keyPathStr, val)
		if err != nil {
			return errors.Wrap(err, "in 'set'")
		}
		t.Values = chartutil.CoalesceTables(vals, t.Values)
	}
	t.Set = nil // no longer used, but make sure this is idempotent.

	return nil
}

// parseDefs parses the `Defs` section into a slice of `*gojq.FuncDef`s, and populates the `funcDefs` field.
func (t *Test) parseDefs() error {
	defsStr := strings.TrimSpace(t.Defs)
	if defsStr == "" {
		return nil
	}
	if !strings.HasSuffix(defsStr, ";") {
		return errors.New("definitions block must end with a semicolon")
	}
	parsedDefs, err := parser.ParseQuery(defsStr)
	if err != nil {
		return errors.Wrap(err, "parsing definitions")
	}
	t.funcDefs = parsedDefs.FuncDefs

	return nil
}

// parsePredicates parses the `Expect` section into a slice of `*gojq.Query` objects, and populates the `predicates`
// field.
func (t *Test) parsePredicates() error {
	expectStr := strings.TrimSpace(t.Expect)
	if expectStr == "" {
		return nil
	}

	predicates, err := parser.ParseExpectations(expectStr)
	if err != nil {
		return errors.Wrap(err, "parsing expectations")
	}

	t.predicates = predicates

	return nil
}

func (t *Test) instantiate(parentParamVals map[string]interface{}) ([]*Test, error) {
	allParamVals := []map[string]interface{}{parentParamVals}
	if t.Params != nil {
		var err error
		allParamVals, err = expandParameters(t.Params)
		if err != nil {
			return nil, errors.Wrap(err, "expanding parameters")
		}
		for _, paramVals := range allParamVals {
			for k, v := range parentParamVals {
				if _, ok := paramVals[k]; !ok {
					paramVals[k] = v
				}
			}
		}
		t.Params = nil
	}

	var instantiations []*Test
	for _, paramVals := range allParamVals {
		instantiation := *t
		if err := instantiation.applyParamValues(paramVals); err != nil {
			return nil, errors.Wrap(err, "applying parameter values")
		}
		instantiations = append(instantiations, &instantiation)
	}
	return instantiations, nil
}

func (t *Test) applyParamValues(paramVals map[string]interface{}) error {
	t.paramValues = paramVals
	if t.Params != nil {
		return errors.New("test still has parameters set")
	}
	if strings.Contains(t.Name, "{{") {
		nameTpl, err := template.New("").Parse(t.Name)
		if err != nil {
			return errors.Wrap(err, "parsing templatized name")
		}
		var buf strings.Builder
		if err := nameTpl.Execute(&buf, paramVals); err != nil {
			return errors.Wrap(err, "executing name template")
		}
		t.Name = buf.String()
	}
	if t.ValuesTemplate != nil {
		valuesTpl, err := template.New("").Parse(*t.ValuesTemplate)
		if err != nil {
			return errors.Wrap(err, "parsing templatized values")
		}
		var buf bytes.Buffer
		if err := valuesTpl.Execute(&buf, paramVals); err != nil {
			return errors.Wrap(err, "executing values template")
		}
		vals, err := chartutil.ReadValues(buf.Bytes())
		if err != nil {
			return errors.Wrap(err, "parsing rendered values template")
		}
		t.Values = chartutil.CoalesceTables(vals, t.Values)
		t.ValuesTemplate = nil
	}

	var instantiatedChildren []*Test
	for _, child := range t.Tests {
		instantiations, err := child.instantiate(paramVals)
		if err != nil {
			return errors.Wrap(err, "instantiating child")
		}
		instantiatedChildren = append(instantiatedChildren, instantiations...)
	}
	t.Tests = instantiatedChildren
	return nil
}

// initialize initializes the test, parsing some string-based values into their semantic counterparts. It also
// recursively initializes the sub-tests. initialize assumes that a name as well as the parent pointer has been set, and
// that the parent is fully initialized.
func (t *Test) initialize() error {
	if err := t.parseDefs(); err != nil {
		return err
	}
	if err := t.applySetOptions(); err != nil {
		return err
	}

	if t.ExpectError == nil {
		if t.parent != nil {
			t.ExpectError = t.parent.ExpectError
		} else {
			t.ExpectError = pointers.Bool(false)
		}
	}

	if err := t.parsePredicates(); err != nil {
		return errors.Wrap(err, "parsing predicates")
	}

	for i, subTest := range t.Tests {
		subTest.parent = t
		if subTest.Name == "" {
			subTest.Name = fmt.Sprintf("#%d", i)
		}
		if err := subTest.initialize(); err != nil {
			return errors.Wrapf(err, "initializing %q", subTest.Name)
		}
	}

	return nil
}

// Run runs a test against the given target.
func (t *Test) Run(testingT *testing.T, tgt *Target) {
	testingT.Run(t.Name, func(testingT *testing.T) {
		testingT.Parallel()
		t.DoRun(testingT, tgt)
	})
}

// DoRun runs a test directly, without an intermediate `testingT.Run` invocation.
func (t *Test) DoRun(testingT *testing.T, tgt *Target) {
	if len(t.Tests) > 0 {
		// non-leaf case
		for _, subTest := range t.Tests {
			subTest.Run(testingT, tgt)
		}
		return
	}

	// leaf case
	runner := &runner{
		t:    testingT,
		test: t,
		tgt:  tgt,
	}
	runner.Run()
}

// forEachScopeBottomUp runs the given doFn function for each test in the hierarchy, starting with the current
// test and ending at the root (suite).
func (t *Test) forEachScopeBottomUp(doFn func(t *Test)) {
	doFn(t)
	if t.parent == nil {
		return
	}
	t.parent.forEachScopeBottomUp(doFn)
}

// forEachScopeTopDown runs the given doFn function for each test in the hierarchy, starting with the root (suite)
// and ending at the current test.
func (t *Test) forEachScopeTopDown(doFn func(t *Test)) {
	if t.parent != nil {
		t.parent.forEachScopeTopDown(doFn)
	}
	doFn(t)
}
