package framework

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// LoadSuite loads a helmtest suite from the given directory.
func LoadSuite(rootDir string) ([]*Test, error) {
	var suite Test
	if err := unmarshalYamlFromFileStrict(filepath.Join(rootDir, "suite.yaml"), &suite); err != nil && !os.IsNotExist(err) {
		return nil, errors.Wrap(err, "loading suite specification")
	}

	if suite.Name == "" {
		suite.Name = strings.TrimRight(rootDir, "/")
	}

	// Locate `.test.yaml` files, if any.
	testYAMLFiles, err := filepath.Glob(filepath.Join(rootDir, "*.test.yaml"))
	if err != nil {
		return nil, errors.Wrap(err, "globbing for .test.yaml files")
	}

	for _, file := range testYAMLFiles {
		test := Test{
			parent: &suite,
		}
		if err := unmarshalYamlFromFileStrict(file, &test); err != nil {
			return nil, errors.Wrapf(err, "loading test specification from file %s", file)
		}
		if test.Name == "" {
			test.Name = filepath.Base(file)
		}
		suite.Tests = append(suite.Tests, &test)
	}

	tests, err := suite.instantiate(nil)
	if err != nil {
		return nil, errors.Wrap(err, "instantiating suite")
	}

	for _, test := range tests {
		if err := test.initialize(); err != nil {
			return nil, err
		}
	}

	return tests, nil
}
