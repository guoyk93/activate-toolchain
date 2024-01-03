package pnpm

import (
	"context"
	"errors"
	"github.com/Masterminds/semver/v3"
	"github.com/go-resty/resty/v2"
	"github.com/yankeguo/activate-toolchain"
	"github.com/yankeguo/activate-toolchain/pkg/ezscript"
	"os"
	"path/filepath"
)

var (
	remapOS = map[string]string{
		"darwin":  "macos",
		"windows": "win",
	}
	remapArch = map[string]string{
		"amd64": "x64",
		"386":   "x86",
		"arm":   "armv7l",
	}
)

type PackageVersion struct {
	Version string `json:"version"`
	Dist    struct {
		Tarball string `json:"tarball"`
	} `json:"dist"`
}

type PackageVersionsResponse struct {
	Versions map[string]PackageVersion `json:"versions"`
}

type toolchain struct {
}

func (t *toolchain) Support(spec activate_toolchain.Spec) bool {
	return spec.Name == "pnpm"
}

func (t *toolchain) resolveURL(ctx context.Context, spec activate_toolchain.Spec) (out string, err error) {
	os, arch := spec.ConvertPlatform(remapOS, remapArch)

	var versions []PackageVersion

	{
		client := resty.New()

		var (
			data PackageVersionsResponse
			res  *resty.Response
		)

		if res, err = client.R().
			SetContext(ctx).
			SetPathParam("os", os).
			SetPathParam("arch", arch).
			SetResult(&data).
			Get("https://registry.npmjs.org/@pnpm/{os}-{arch}"); err != nil {
			return
		}
		if res.StatusCode() != 200 {
			err = errors.New("unexpected status code: " + res.Status())
			return
		}

		for _, v := range data.Versions {
			versions = append(versions, v)
		}
	}

	var best PackageVersion

	if best, err = activate_toolchain.FindBestVersionedItem(
		spec.VersionConstraints,
		versions,
		func(i PackageVersion) (*semver.Version, error) {
			return semver.NewVersion(i.Version)
		},
	); err != nil {
		return
	}

	out = best.Dist.Tarball

	return
}

func (t *toolchain) Activate(ctx context.Context, spec activate_toolchain.Spec) (script string, err error) {
	var dir string

	if dir, err = activate_toolchain.InstallArchive(ctx, activate_toolchain.InstallArchiveOptions{
		ProvideURLs: func() (urls []string, err error) {
			var u string
			if u, err = t.resolveURL(ctx, spec); err != nil {
				return
			}
			urls = append(urls, u)
			return
		},
		Name: spec.VersionedName(),
		File: spec.VersionedName() + ".tar.gz",
	}); err != nil {
		return
	}

	// make it executable
	exe := filepath.Join(dir, "package", "pnpm")

	if err = os.Chmod(exe, 0755); err != nil {
		return
	}

	// create bin dir and create a symlink
	dirBin := filepath.Join(dir, "bin")

	if err = os.MkdirAll(dirBin, 0755); err != nil {
		return
	}

	os.RemoveAll(filepath.Join(dirBin, "pnpm"))

	if err = os.Symlink(filepath.Join("..", "package", "pnpm"), filepath.Join(dirBin, "pnpm")); err != nil {
		return
	}

	return ezscript.Render(
		map[string]any{
			"dir_bin": dirBin,
		},
		`{{setEnv "PNPM_HOME" (filepathJoin "$HOME" ".pnpm")}}`,
		`{{addEnv "PATH" "$PNPM_HOME" .dir_bin}}`,
	)
}

func init() {
	activate_toolchain.AddToolchain(&toolchain{})
}
