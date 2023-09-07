package executor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExecute(t *testing.T) {
	res, err := RunCmd("top -b -n1 | grep average")
	require.Nil(t, err, "Error is not nil!")
	require.Contains(t, res, "load average", "Grep didn't work due to execution of top command")
}
