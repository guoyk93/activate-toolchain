package node

import (
	"context"
	"fmt"
	"github.com/guoyk93/activate-toolchain"
	"github.com/guoyk93/rg"
)

const (
	indexURL = "https://nodejs.org/download/release/index.json"
)

var (
	downloadURLs = []string{
		"https://mirrors.cloud.tencent.com/nodejs-release",
		"https://mirrors.aliyun.com/nodejs-release",
		"https://nodejs.org/download/release",
	}
)

type VersionItem struct {
	Version string   `json:"version"`
	Date    string   `json:"date"`
	Files   []string `json:"files"`
}

func createFilenameSuffix(os string, arch string) string {
	if arch == "amd64" {
		arch = "x64"
	}
	if arch == "386" {
		arch = "x86"
	}
	if arch == "arm" {
		arch = "armv7l"
	}
	return "-" + os + "-" + arch + ".tar.gz"
}

type toolchain struct {
}

func (t *toolchain) resolveVersion(ctx context.Context, targetVersion string) (version string, err error) {
	var data []VersionItem
	rg.Must0(activate_toolchain.FetchJSON(ctx, indexURL, &data))

	var (
		versions     []string
		versionItems = make(map[string]VersionItem)
	)

	{
		for _, item := range data {
			versionItems[item.Version] = item
		}

		for name := range versionItems {
			versions = append(versions, name)
		}
	}

	if version, err = activate_toolchain.ResolveVersion(versions, targetVersion); err != nil {
		return
	}

	return
}

func (t *toolchain) resolveFilename(version, os, arch string) (rel string) {
	return "node-" + version + createFilenameSuffix(os, arch)
}

func (t *toolchain) Name() string {
	return "node"
}

func (t *toolchain) Activate(ctx context.Context, targetVersion, os, arch string) (script string, err error) {
	defer rg.Guard(&err)

	version := rg.Must(t.resolveVersion(ctx, targetVersion))
	filename := t.resolveFilename(version, os, arch)

	var urls []string
	for _, url := range downloadURLs {
		urls = append(urls, url+"/"+version+"/"+filename)
	}

	var dir string
	if dir, err = activate_toolchain.InstallArchive(ctx, activate_toolchain.InstallArchiveOptions{
		URLs:           urls,
		Filename:       filename,
		Name:           "node-" + targetVersion,
		StripDirectory: true,
	}); err != nil {
		return
	}

	script = fmt.Sprintf(
		`
export PATH="%s/bin:$PATH";
echo "node %s activated";
`,
		dir,
		version,
	)

	return
}

func init() {
	activate_toolchain.Toolchains = append(activate_toolchain.Toolchains, &toolchain{})
}
