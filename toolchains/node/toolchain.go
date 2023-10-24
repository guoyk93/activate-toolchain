package node

import (
	"context"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/guoyk93/activate-toolchain"
	"github.com/guoyk93/activate-toolchain/pkg/ezscript"
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

	remapOS   = map[string]string{}
	remapArch = map[string]string{
		"amd64": "x64",
		"386":   "x86",
		"arm":   "armv7l",
	}
)

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
	os, arch := spec.ConvertPlatform(remapOS, remapArch)

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

	return ezscript.Render(
		map[string]any{
			"dir": dir,
		},
		`{{addEnv "PATH" (filepathJoin .dir "bin")}}`,
	)
}

func init() {
	activate_toolchain.AddToolchain(&toolchain{})
}
