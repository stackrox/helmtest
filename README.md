Package `helmtest`
======

The `helmtest` package allows you to declaratively specify test suites for Helm charts. It specifically
seeks to address two inconveniences of the "normal" Go unit test-based approach:
- it allows testing a multitude of different configurations via a hierarchical, YAML-based specification
  of test cases.
- it makes writing assertions about the generated Kubernets objects easy, by using `jq` filters as the
  assertion language.
  
Format of a test
=========
A `helmtest` test is generally defined in a YAML file according to the format specified in `spec.go`.
Tests are organized in a hierarchical fashion, in the sense that a test may contain one or more
sub-tests. Tests with no sub-tests are called "leaf tests", and other tests are called "non-leaf tests".
A Helm chart is only rendered and checked against expectations in leaf tests; in such a setting,
the leaf test inherits certain properties from its non-leaf ancestors.

The general schema of a test is as follows:
```yaml
name: "string" # the name of the test (optional but strongly recommended). Auto-generated if left empty.
release:  # Overrides for the Helm release properties. These are applied in root-to-leaf order.
  name: "string"  # override for the Helm release name
  namespace: "string"  # override for the Helm release namespace
  revision: int # override for the Helm release revision
  isInstall: bool # override for the "IsInstall" property of the release options
  isUpgrade: bool # override for the "IsUpgrade" property of the release options
server:
  visibleSchema: # openAPI schema which is visible to helm, i.e. to check API resource availability
  # all valid schemas are:
  - kubernetes-1.20.2
  - openshift-3.11.0
  - openshift-4.1.0
  - com.coreos
  availableSchemas: [] # openAPI schema to validate against, i.e. to validate if rendered objects could be applied
values:  # values as consumed by Helm via the `-f` CLI flag.
  key: value
set:  # alternative format for Helm values, as consumed via the `--set` CLI flag.
  nes.ted.key: value
defs: |
  # Sequence of jq "def" statements. Each def statement must be terminated with a semicolon (;). Defined functions
  # are only visible in this and descendant scopes, but not in ancestor scopes.
  def name: .metadata.name;

expectError: bool # indicates whether we can tolerate an error. If unset, inherit from the parent test, or
                  # default to `false` at the root level.
expect: |
  # Sequence of jq filters, one per line (or spanning multiple lines, where each continuation line must begin with a
  # space).
  # See the below section on the world model and special functions.
  .objects[] | select(.metadata?.name? // "" == "")
    | assertNotExist  # continuation line
tests: []  # a list of sub-tests. Determines whether the test is a leaf test or non-leaf test.
```

A comprehensive set of hierarchically organized tests to be run against a Helm chart is called a "suite". Each suite
is defined in a set of YAML files located in a single directory on the filesystem (a directory may hold at most one
suite). The properties of the top-level test in the suite (such as a common set of expectations or Helm values to be
inherited by all tests) can be specified in a `suite.yaml` file within this directory. The `suite.yaml` file may be
absent, in which case there are no values, definitions, expectations etc. shared among all the tests in the suite. In
addition to the tests specified in the `tests:` stanza of the `suite.yaml` file (if any), the tests of the suite are
additionally read from files with the extension `.test.yaml` in the suite directory. Note that any combination of
defining tests in the `suite.yaml` and in individual files may be used, these tests will then be combined. In
particular, it is possible to define arbitrary test suites either with only `.test.yaml` files, with only a `suite.yaml`
file, or with combinations thereof.

Inheritance
----------------
For most fields in the test specification, a test will inherit the value from its parent test (which might use an
inherited value as well, etc.). If an explicit value is given, this value
- overrides the values from the parent for the following fields: `expectError` and the individual sub-fields of
  `release`.
- is merged with the values from the parent for the following fields: `values`, `set` (in such a way that the values
  from the child take precedence).
- is added to the values from the parents for the following fields: `expect`, `defs`.

World model
============

As stated above, expectations are written as `jq` filters (using `gojq` as an evaluation engine). Generally, a filter
that evaluates to a "falsy" value is treated as a violation. In contrast to normal JS/`jq` semantics, an empty list,
object, or string will also be treated as "falsy". The input to those filters (i.e., the `.` value at the start of
each `jq` pipeline) is a JSON object containing anything that is relevant to the test execution, referred to as the
"world". See [the world model documentation](./docs/world-model.md) for an explanation of what it contains.

Special functions
===============

See the [documentation on functions](./docs/functions.md) for an overview of what functions are available in filters,
beyond the ones known from `jq`.

Debugging
===============

**Run a single test:**

To get the name of a single test:
1. Run the whole test suite
2. Look up the test to execute
3. Copy the complete name from the CLI (or logs)

```
$ go test -run "TestWithHelmtest/testdata/helmtest/some_values.test.yaml/some_example_test"
```

**Print rendered manifests and values:**

```
- name: "some example test"
  expect: |
    .helm[]| toyaml | print                  ## Print all Helm values
    .secrets["secret-name"] | toyaml | print ## Print a specific object
    .objects[]| toyaml | print               ## Print the complete rendered Helm chart
```
