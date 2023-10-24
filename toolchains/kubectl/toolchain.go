package kubectl

import (
	"context"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/guoyk93/activate-toolchain"
	"github.com/guoyk93/activate-toolchain/pkg/ezs3"
	"github.com/guoyk93/activate-toolchain/pkg/ezscript"
	"strings"
)

const endpointURL = "https://cdn.dl.k8s.io"

type toolchain struct {
}

func (t *toolchain) resolveURL(ctx context.Context, spec activate_toolchain.Spec) (out string, err error) {
	var (
		prefixes []ezs3.ListBucketCommonPrefix
	)

	if _, prefixes, err = ezs3.ListObjects(ctx, ezs3.ListObjectsOptions{
		Endpoint:  endpointURL,
		Prefix:    "release/",
		Delimiter: "/",
	}); err != nil {
		return
	}

	var best ezs3.ListBucketCommonPrefix

	if best, err = activate_toolchain.FindBestVersionedItem(
		spec.VersionConstraints,
		prefixes,
		func(v ezs3.ListBucketCommonPrefix) (version *semver.Version, err error) {
			return semver.NewVersion(
				strings.TrimSuffix(
					strings.TrimPrefix(v.Prefix, "release/"),
					"/",
				),
			)
		},
	); err != nil {
		return
	}

	out = fmt.Sprintf(
		"%s/%s/kubernetes-client-%s-%s.tar.gz",
		endpointURL,
		strings.TrimSuffix(strings.TrimPrefix(best.Prefix, "/"), "/"),
		spec.OS,
		spec.Arch,
	)

	return
}

func (t *toolchain) Support(spec activate_toolchain.Spec) bool {
	return spec.Name == "kubectl"
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
		Name:           spec.VersionedName(),
		File:           spec.VersionedName() + ".tar.gz",
		DirectoryLevel: 2,
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
