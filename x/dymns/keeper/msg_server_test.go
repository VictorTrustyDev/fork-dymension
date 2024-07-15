package keeper

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	require.Equal(
		t,
		time.Now().Unix(), time.Now().UTC().Unix(),
		"if mis-match, 100% sure will causes AppHash",
	)
}
