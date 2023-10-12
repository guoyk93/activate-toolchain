package node

import (
	"context"
	"fmt"
	"github.com/guoyk93/activate-toolchain"
)

const (
	indexURL = "https://nodejs.org/download/release/index.json"
)

var (
	baseURLs = []string{
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

func remapArch(arch string) string {
	if v := map[string]string{
		"amd64": "x64",
		"386":   "x86",
		"arm":   "armv7l",
	}[arch]; v != "" {
		return v
	}
	return arch
}

type toolchain struct {
}

func (t *toolchain) resolveVersion(ctx context.Context, targetVersion string) (version string, err error) {
	var data []VersionItem
	if err = activate_toolchain.FetchJSON(ctx, indexURL, &data); err != nil {
		return
	}

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

func (t *toolchain) Name() string {
	return "node"
}

func (t *toolchain) Activate(ctx context.Context, targetVersion, os, arch string) (script string, err error) {
	var dir string

	if dir, err = activate_toolchain.InstallArchive(ctx, activate_toolchain.InstallArchiveOptions{
		ProvideURLs: func() (urls []string, err error) {
			var version string
			if version, err = t.resolveVersion(ctx, targetVersion); err != nil {
				return
			}

			for _, url := range baseURLs {
				urls = append(urls, fmt.Sprintf("%s/%s/node-%s-%s-%s.tar.gz", url, version, version, os, remapArch(arch)))
			}

			return
		},
		Name:           "node-" + targetVersion,
		File:           "node-" + targetVersion + ".tar.gz",
		DirectoryLevel: 1,
	}); err != nil {
		return
	}

	script = fmt.Sprintf(`
export PATH="%s/bin:$PATH";
`, dir)

	return
}

func init() {
	activate_toolchain.Toolchains = append(activate_toolchain.Toolchains, &toolchain{})
}
