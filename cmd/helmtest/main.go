package main

import (
	"flag"
	"os"
	"testing"

	"github.com/stackrox/helmtest/pkg/framework"

	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
)

func mainCmd() error {
	var chartPath string
	var help bool

	var releaseName string
	var namespace string
	var upgrade bool

	flag.BoolVar(&help, "help", false, "print help information and exit")

	flag.StringVar(&chartPath, "chart", "", "archive file or directory containing the chart to test")

	flag.BoolVar(&upgrade, "upgrade", false, "if set, render chart as if performing an upgrade instead of an installation")
	flag.StringVar(&releaseName, "release-name", "helmtest-release", "the name of the Helm release")
	flag.StringVar(&namespace, "namespace", "default", "the namespace into which to simulate installing")

	testing.Init()

	flag.Parse()

	if help {
		flag.Usage()
		return nil
	}

	if chartPath == "" {
		return errors.New("no chart specified")
	}

	suiteDirs := flag.Args()
	if len(suiteDirs) == 0 {
		return errors.New("no test suites specified")
	}

	st, err := os.Stat(chartPath)
	if err != nil {
		return errors.Wrap(err, "loading chart")
	}
	var chartToTest *chart.Chart
	switch {
	case st.IsDir():
		chartToTest, err = loader.LoadDir(chartPath)
	case st.Mode().IsRegular():
		chartToTest, err = loader.LoadFile(chartPath)
	default:
		return errors.Errorf("invalid chart %q: neither a directory nor a regular file", chartPath)
	}

	if err != nil {
		return errors.Wrapf(err, "loading chart %q", chartPath)
	}

	target := &framework.Target{
		ReleaseOptions: chartutil.ReleaseOptions{
			Name:      releaseName,
			Namespace: namespace,
			Revision:  1,
			IsUpgrade: upgrade,
			IsInstall: !upgrade,
		},
		Chart: chartToTest,
	}

	var suites []*framework.Test
	for _, suiteDir := range suiteDirs {
		topLevelTests, err := framework.LoadSuite(suiteDir)
		if err != nil {
			return errors.Wrapf(err, "loading suite %q", suiteDir)
		}
		suites = append(suites, topLevelTests...)
	}

	tests := make([]testing.InternalTest, 0, len(suites))
	for _, suite := range suites {
		s := suite
		tests = append(tests, testing.InternalTest{
			Name: suite.Name,
			F: func(t *testing.T) {
				t.Parallel()
				s.DoRun(t, target)
			},
		})
	}

	testing.Main(func(string, string) (bool, error) { return true, nil }, tests, nil, nil)
	return nil
}

func main() {
	if err := mainCmd(); err != nil {
		panic(err)
	}
}
