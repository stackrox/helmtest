package framework

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoader(t *testing.T) {
	const testdataPath = "testdata/suite"

	tests := map[string]struct {
		opts          []LoaderOpt
		expectedFunc  func(*testing.T, *Test)
		additionalDir string
	}{
		"With root dir": {
			expectedFunc: func(t *testing.T, helmTest *Test) {
				assert.Len(t, helmTest.Tests, 2)
			},
		},
		"Loader loads test hierarchy": {
			expectedFunc: func(t *testing.T, test *Test) {
				require.Len(t, test.Tests[1].Tests, 1)
				childTest := test.findFirst([]string{testdataPath, "helm.test.yaml", "test in helm.test.yaml", "with overwrites"})
				assert.Equal(t, "with overwrites", childTest.Name)
				assert.EqualValues(t, map[string]interface{}{"testValue": "value overwrite"}, childTest.Values)
			},
		},
		"Loader loads additional dir": {
			additionalDir: "testdata/additional_dir",
			expectedFunc: func(t *testing.T, test *Test) {
				childTest := test.findFirst([]string{testdataPath, "additional.test.yaml"})
				require.NotNil(t, test)
				assert.Equal(t, "additional.test.yaml", childTest.Name)
			},
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			var opts []LoaderOpt
			if tt.additionalDir != "" {
				opts = append(opts, WithAdditionalTestDirs(tt.additionalDir))
			}

			loader := NewLoader(testdataPath, opts...)
			helmTests, err := loader.LoadSuite()
			require.NoError(t, err)

			tt.expectedFunc(t, helmTests)
		})
	}
}
