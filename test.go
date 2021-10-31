package grpcsteps

import "fmt"

type testT struct {
	error error
}

func (t *testT) Errorf(format string, args ...interface{}) {
	t.error = fmt.Errorf(format, args...) // nolint: goerr113
}
