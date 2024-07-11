package keeper_test

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func requireErrorContains(t *testing.T, err error, contains string) {
	require.Error(t, err)
	require.NotEmpty(t, contains, "mis-configured test")
	require.Contains(t, err.Error(), contains)
}

func requireErrorFContains(t *testing.T, f func() error, contains string) {
	requireErrorContains(t, f(), contains)
}
