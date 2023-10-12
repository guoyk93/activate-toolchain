package activate_toolchain

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestDetectFastestURL(t *testing.T) {
	urls := []string{
		"https://mirrors.cloud.tencent.com/nodejs-release/index.json",
		"https://mirrors.aliyun.com/nodejs-release/index.json",
		"https://nodejs.org/download/release/index.json",
	}
	fastest, err := DetectFastestURL(context.Background(), urls)
	require.NoError(t, err)
	require.Contains(t, urls, fastest)
}

func TestAdvancedFetchFile(t *testing.T) {
	urls := []string{
		"https://mirrors.cloud.tencent.com/nodejs-release/index.json",
		"https://mirrors.aliyun.com/nodejs-release/index.json",
		"https://nodejs.org/download/release/index.json",
	}
	file := filepath.Join(os.TempDir(), "test-fetch-file-from-candidate-urls.json")
	err := AdvancedFetchFile(context.Background(), urls, file)
	require.NoError(t, err)
	buf, err := os.ReadFile(file)
	require.NoError(t, err)
	require.True(t, bytes.HasPrefix(buf, []byte("[")))
}
