package jdk

import (
	"context"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/guoyk93/activate-toolchain"
	"log"
	"path/filepath"
	"strings"
)

func remapOS(os string) string {
	if v := map[string]string{
		"darwin": "mac",
	}[os]; v != "" {
		return v
	}
	return os
}

func remapArch(arch string) string {
	if v := map[string]string{
		"amd64": "x64",
		"386":   "x86",
		"arm64": "aarch64",
	}[arch]; v != "" {
		return v
	}
	return arch
}

func validateAdoptiumFilename(file string, os string, arch string) bool {
	return strings.HasSuffix(file, ".tar.gz") &&
		strings.Contains(
			file,
			fmt.Sprintf(
				"jdk_%s_%s_hotspot_",
				remapArch(arch),
				remapOS(os),
			))
}

type toolchain struct {
}

func (t *toolchain) Name() string {
	return "jdk"
}

func (t *toolchain) Activate(ctx context.Context, targetVersion, os, arch string) (script string, err error) {
	var dirPath string
	if os == "darwin" {
		dirPath = filepath.Join("Contents", "Home")
	}

	var version *semver.Version
	if version, err = semver.NewVersion(targetVersion); err != nil {
		return
	}

	usePrimary := true

	opts := activate_toolchain.InstallArchiveOptions{
		ProvideURLs: func() (urls []string, err error) {
			for _, src := range sources {
				if src.Primary() != usePrimary {
					continue
				}
				if u, e := src.Resolve(ctx, version, os, arch); e == nil {
					urls = append(urls, u)
				} else {
					log.Println("failed to resolve source:", src.Name(), ":", e)
				}
			}
			return
		},
		Name:           "jdk-" + targetVersion,
		File:           "jdk-" + targetVersion + ".tar.gz",
		DirectoryLevel: 1,
		DirectoryPath:  dirPath,
	}

	var dir string

	if dir, err = activate_toolchain.InstallArchive(ctx, opts); err != nil {
		log.Println("trying to use secondary sources")
		usePrimary = false
		if dir, err = activate_toolchain.InstallArchive(ctx, opts); err != nil {
			return
		}
	}

	script = fmt.Sprintf(`
export JAVA_HOME="%s";
export PATH="$JAVA_HOME/bin:$PATH";
`, dir)

	return
}

func init() {
	activate_toolchain.Toolchains = append(activate_toolchain.Toolchains, &toolchain{})
}
