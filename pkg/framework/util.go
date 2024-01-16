package framework

import (
	"k8s.io/kubectl/pkg/util/openapi"
	"os"

	"github.com/pkg/errors"
	yamlv3 "gopkg.in/yaml.v3"
)

// unmarshalYamlFromFileStrict unmarshals the contents of filename into out, relying on gopkg.in/yaml.v3 semantics.
// Any field that is not present in the output data type, as well as any duplicate keys within the same YAML object,
// will result in an error.
func unmarshalYamlFromFileStrict(filename string, out interface{}) error {
	contents, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	var root yamlv3.Node
	if err := yamlv3.Unmarshal(contents, &root); err != nil {
		return errors.Wrapf(err, "parsing YAML in file %s", filename)
	}
	injectFilename(&root, filename)
	if err := root.Decode(out); err != nil {
		return errors.Wrapf(err, "decoding YAML in file %s", filename)
	}
	return nil
}

type openAPIResourcesGetter struct {
	resources openapi.Resources
}

func (o openAPIResourcesGetter) OpenAPISchema() (openapi.Resources, error) {
	return o.resources, nil
}

// OpenAPIResourcesGetter returns resources wrapped in a simple implementation of openapi.OpenAPIResourcesGetter
// that returns a fixed set of resources.
func OpenAPIResourcesGetter(resources openapi.Resources) openapi.OpenAPIResourcesGetter {
	return openAPIResourcesGetter{resources: resources}
}
