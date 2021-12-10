package framework

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

const (
	defaultTestFileGlobPattern = "*.test.yaml"
)

// Loader loads a test suite.
type Loader struct {
	globPattern string
	rootDir     string
}

// NewLoader returns a a loader and applies options to it.
func NewLoader(rootDir string, opts ...LoaderOpt) *Loader {
	loader := Loader{
		rootDir:     rootDir,
		globPattern: defaultTestFileGlobPattern,
	}

	for _, opt := range opts {
		opt(&loader)
	}
	return &loader
}

// LoaderOpts allows to set custom options.
type LoaderOpt func(loader *Loader)

// WithCustomFilePattern sets a custom file pattern to load test files.
func WithCustomFilePattern(pattern string) LoaderOpt {
	return func(loader *Loader) {
		loader.globPattern = pattern
	}
}

// LoadSuite loads a helmtest suite from the given directory.
func (loader *Loader) LoadSuite() (*Test, error) {
	var suite Test
	if err := unmarshalYamlFromFileStrict(filepath.Join(loader.rootDir, "suite.yaml"), &suite); err != nil && !os.IsNotExist(err) {
		return nil, errors.Wrap(err, "loading suite specification")
	}

	if suite.Name == "" {
		suite.Name = strings.TrimRight(loader.rootDir, "/")
	}

	// Locate test files, if any.
	testYAMLFiles, err := filepath.Glob(filepath.Join(loader.rootDir, loader.globPattern))
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

	if err := suite.initialize(); err != nil {
		return nil, err
	}

	return &suite, nil
}

// LoadSuiteWithFlavour loads a helmtest suite from the given directory with a specific flavour pre-configured.
func (loader *Loader) LoadSuiteWithFlavour(defaultFlavour string) (*Test, error) {
	var suite Test

	if err := unmarshalYamlFromFileStrict(filepath.Join(loader.rootDir, "suite.yaml"), &suite); err != nil && !os.IsNotExist(err) {
		return nil, errors.Wrap(err, "loading suite specification")
	}

	if suite.Flavour == "" {
		suite.Flavour = defaultFlavour
	}

	if suite.Name == "" {
		suite.Name = strings.TrimRight(loader.rootDir, "/")
	}

	// Locate test files, if any.
	testYAMLFiles, err := filepath.Glob(filepath.Join(loader.rootDir, loader.globPattern))
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

	if err := suite.initialize(); err != nil {
		return nil, err
	}

	return &suite, nil
}
