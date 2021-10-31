package grpcsteps_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/cucumber/godog"

	"github.com/godogx/grpcsteps"
)

func runSuite(t *testing.T, c *grpcsteps.Client, paths ...string) {
	t.Helper()

	buf := bytes.NewBuffer(nil)

	suite := godog.TestSuite{
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			c.RegisterContext(ctx)
		},
		Options: &godog.Options{
			Format:    "pretty",
			Output:    buf,
			Paths:     paths,
			Strict:    true,
			Randomize: time.Now().UTC().UnixNano(),
		},
	}

	if status := suite.Run(); status != 0 {
		t.Fatal(buf.String())
	}
}
