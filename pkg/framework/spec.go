package framework

import (
	"github.com/itchyny/gojq"
	"github.com/stackrox/helmtest/internal/parser"
)

// RawDict is an alias for map[string]interface{}, that is needed because `yaml.Unmarshal` and `json.Unmarshal` differ
// in that the latter will never produce int values, while the former may.
type RawDict map[string]interface{}

// Test defines a helmtest test. A Test can be regarded as the equivalent of the *testing.T scope of a Go unit test.
// Tests are scoped, and a test may either define concrete expectations, or contain an arbitrary number of nested tests.
// See README.md in this directory for a more detailed explanation.
type Test struct {
	// Public section - fields settable via YAML

	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	Values RawDict `json:"values,omitempty" yaml:"values,omitempty"`
	Set    RawDict `json:"set,omitempty" yaml:"set,omitempty"`

	Defs         string                   `json:"defs,omitempty" yaml:"defs,omitempty"`
	Release      *ReleaseSpec             `json:"release,omitempty" yaml:"release,omitempty"`
	Server       *ServerSpec              `json:"server,omitempty" yaml:"server,omitempty"`
	Capabilities *CapabilitiesSpec        `json:"capabilities,omitempty" yaml:"capabilities,omitempty"`
	Objects      []map[string]interface{} `json:"objects,omitempty" yaml:"objects,omitempty"`

	Expect      string `json:"expect,omitempty" yaml:"expect,omitempty"`
	ExpectError *bool  `json:"expectError,omitempty" yaml:"expectError,omitempty"`

	Tests []*Test `json:"tests,omitempty" yaml:"tests,omitempty"`

	// Private section - the following fields are never set in the YAML, they are always populated by initialize
	// or during YAML parsing.
	parent *Test

	funcDefs   []*gojq.FuncDef
	predicates []*parser.ParsedQuery

	defsSrcCtx   parser.SourceContext
	expectSrcCtx parser.SourceContext
}

// ReleaseSpec specifies how the release options for Helm will be constructed.
type ReleaseSpec struct {
	Name      string `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Revision  *int   `json:"revision,omitempty" yaml:"revision,omitempty"`
	IsInstall *bool  `json:"isInstall,omitempty" yaml:"isInstall,omitempty"`
	IsUpgrade *bool  `json:"isUpgrade,omitempty" yaml:"isUpgrade,omitempty"`
}

// ServerSpec specifies how the model of the server will be constructed.
type ServerSpec struct {
	// AvailableSchemas are the names of schemas that are available on the server (i.e., that rendered objects must
	// pass validation against, but not necessarily discoverable via `.Capabilities.APIVersions`).
	AvailableSchemas []string `json:"availableSchemas,omitempty" yaml:"availableSchemas,omitempty"`
	// VisibleSchemas are the names of schemas that are available on the server AND discoverable via
	// `.Capabilities.APIVersions`.
	VisibleSchemas []string `json:"visibleSchemas,omitempty" yaml:"visibleSchemas,omitempty"`

	// NoInherit indicates that server-side settings should *not* be inherited from the enclosing scope.
	NoInherit bool `json:"noInherit,omitempty" yaml:"noInherit,omitempty"`
}

// CapabilitiesSpec represents the `Capabilities` in Helm.
type CapabilitiesSpec struct {
	// KubeVersion represents the kubernetes version which is discoverable via `.Capabilities.KubeVersion`.
	KubeVersion *KubeVersion `json:"kubeVersion,omitempty" yaml:"kubeVersion,omitempty"`
}

// KubeVersion is the Kubernetes version.
type KubeVersion struct {
	Version string `json:"version,omitempty" yaml:"version,omitempty"` // i.e. v1.18
	Major   string `json:"major,omitempty" yaml:"major,omitempty"`     // i.e. 1
	Minor   string `json:"minor,omitempty" yaml:"minor,omitempty"`     // i.e. 18
}
