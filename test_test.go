package grpcsteps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestT_Errorf(t *testing.T) {
	t.Parallel()

	tt := &testT{}

	tt.Errorf("error: %s", "t")

	expected := `error: t`

	assert.EqualError(t, tt.error, expected)
}
