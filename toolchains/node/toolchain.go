package node

import (
	"context"
	"fmt"
	"github.com/Masterminds/semver/v3"
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

func convertPlatform(spec activate_toolchain.Spec) (os string, arch string) {
	os, arch = spec.OS, spec.Arch

	switch arch {
	case "amd64":
		arch = "x64"
	case "386":
		arch = "x86"
	case "arm":
		arch = "armv7l"
	}

	return
}

type VersionItem struct {
	Version string   `json:"version"`
	Date    string   `json:"date"`
	Files   []string `json:"files"`
}

type toolchain struct{}

func (t *toolchain) resolveBestVersion(ctx context.Context, spec activate_toolchain.Spec) (version string, err error) {
	var data []VersionItem
	if err = activate_toolchain.FetchJSON(ctx, indexURL, &data); err != nil {
		return
	}

	var versions []string

	{
		versionItems := make(map[string]VersionItem)

		for _, item := range data {
			versionItems[item.Version] = item
		}

		for name := range versionItems {
			versions = append(versions, name)
		}
	}

	if version, err = activate_toolchain.FindBestVersionedItem(
		spec.VersionConstraints,
		versions,
		semver.NewVersion,
	); err != nil {
		return
	}

	return
}

func (t *toolchain) Support(spec activate_toolchain.Spec) bool {
	return spec.Name == "node"
}

func (t *toolchain) Activate(ctx context.Context, spec activate_toolchain.Spec) (script string, err error) {
	os, arch := convertPlatform(spec)

	var dir string

	if dir, err = activate_toolchain.InstallArchive(ctx, activate_toolchain.InstallArchiveOptions{
		ProvideURLs: func() (urls []string, err error) {
			var version string
			if version, err = t.resolveBestVersion(ctx, spec); err != nil {
				return
			}

			for _, url := range baseURLs {
				urls = append(urls, fmt.Sprintf("%s/%s/node-%s-%s-%s.tar.gz", url, version, version, os, arch))
			}

			return
		},
		Name:           spec.VersionedName(),
		File:           spec.VersionedName() + ".tar.gz",
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
	activate_toolchain.AddToolchain(&toolchain{})
}
