package grpcsteps_test

import (
	"bytes"
	"fmt"
	"time"

	"github.com/cucumber/godog"
)

type suiteT interface {
	Errorf(format string, args ...interface{})
}

type suiteOption func(ts *godog.TestSuite)

func runSuite(t suiteT, opts ...suiteOption) {
	buf := bytes.NewBuffer(nil)

	suite := godog.TestSuite{
		Options: &godog.Options{
			Format:    "pretty",
			Output:    buf,
			Strict:    true,
			Randomize: time.Now().UTC().UnixNano(),
		},
	}

	for _, o := range opts {
		o(&suite)
	}

	if status := suite.Run(); status != 0 {
		t.Errorf(buf.String())
	}
}

func initSuite(init func(ctx *godog.TestSuiteContext)) suiteOption {
	return func(ts *godog.TestSuite) {
		ts.TestSuiteInitializer = init
	}
}

func afterSuite(fn func()) suiteOption {
	return initSuite(func(tsc *godog.TestSuiteContext) {
		tsc.AfterSuite(fn)
	})
}

func initScenario(initializers ...func(ctx *godog.ScenarioContext)) suiteOption {
	return func(ts *godog.TestSuite) {
		ts.ScenarioInitializer = func(ctx *godog.ScenarioContext) {
			for _, i := range initializers {
				i(ctx)
			}
		}
	}
}

func featureFiles(paths ...string) suiteOption {
	return func(ts *godog.TestSuite) {
		ts.Options.Paths = paths
	}
}

func noColors() suiteOption {
	return func(ts *godog.TestSuite) {
		ts.Options.NoColors = true
	}
}

type testT struct {
	error error
}

func (t *testT) Errorf(format string, args ...interface{}) {
	t.error = fmt.Errorf(format, args...)
}
