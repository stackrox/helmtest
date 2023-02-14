package integration_test

import (
	"testing"

	"github.com/stackrox/helmtest/pkg/framework"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
)

func TestHelmTestShouldSucceed(t *testing.T) {
	l := framework.NewLoader("testdata/helmtest", framework.WithAdditionalTestDirs("testdata/helmtest/additional_directory"))
	s, err := l.LoadSuite()
	require.NoError(t, err)

	chart, err := loader.Load("testdata/nginx-example")
	require.NoError(t, err)

	target := &framework.Target{
		Chart: chart,
		ReleaseOptions: chartutil.ReleaseOptions{
			Name:      "nginx-lb",
			Namespace: "loadbalancer",
			IsInstall: true,
		},
	}

	s.Run(t, target)
}
