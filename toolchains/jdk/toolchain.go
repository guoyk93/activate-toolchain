package jdk

import (
	"context"
	"fmt"
	"github.com/guoyk93/activate-toolchain"
	"log"
	"path/filepath"
	"strings"
)

func convertPlatform(spec activate_toolchain.Spec) (os, arch string) {
	os, arch = spec.OS, spec.Arch
	switch os {
	case "darwin":
		os = "mac"
	}
	switch arch {
	case "amd64":
		arch = "x64"
	case "386":
		arch = "x86"
	case "arm64":
		arch = "aarch64"
	}
	return
}

func validateAdoptiumFilename(file string, spec activate_toolchain.Spec) bool {
	os, arch := convertPlatform(spec)
	return strings.HasSuffix(file, ".tar.gz") &&
		strings.Contains(file, fmt.Sprintf("jdk_%s_%s_hotspot_", arch, os))
}

type toolchain struct {
}

func (t *toolchain) Support(spec activate_toolchain.Spec) bool {
	return spec.Name == "jdk"
}

func (t *toolchain) Activate(ctx context.Context, spec activate_toolchain.Spec) (script string, err error) {
	var dirPath string
	if spec.OS == "darwin" {
		dirPath = filepath.Join("Contents", "Home")
	}

	usePrimary := true

	opts := activate_toolchain.InstallArchiveOptions{
		ProvideURLs: func() (urls []string, err error) {
			for _, src := range sources {
				if src.Primary() != usePrimary {
					continue
				}
				if u, err1 := src.ResolveDownloadURL(ctx, spec); err1 == nil {
					urls = append(urls, u)
				} else {
					log.Println("failed to resolve source:", src.Name(), ":", err1.Error())
				}
			}
			return
		},
		Name:           spec.VersionedName(),
		File:           spec.VersionedName() + ".tar.gz",
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
