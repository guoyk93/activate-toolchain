package activate_toolchain

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseSpec(t *testing.T) {
	spec, err := ParseSpec("node@v0.10.26")
	require.NoError(t, err)
	require.Equal(t, "node", spec.Name)
	require.Equal(t, "node-0.10.26", spec.VersionedName())
}
