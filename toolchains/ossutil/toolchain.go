package ossutil

import (
	"context"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/go-resty/resty/v2"
	"github.com/guoyk93/activate-toolchain"
	"log"
	"net/http"
	"strings"
)

var (
	remapOS = map[string]string{
		"darwin": "mac",
	}
	remapArch = map[string]string{}
)

type toolchain struct {
}

func (t *toolchain) resolveURL(ctx context.Context, spec activate_toolchain.Spec) (u string, err error) {
	client := resty.New()

	var version string
	{
		var res *resty.Response
		if res, err = client.R().SetContext(ctx).Get("https://gosspublic.alicdn.com/ossutil/version.txt"); err != nil {
			return
		}

		body := res.String()

		if res.StatusCode() != http.StatusOK {
			err = fmt.Errorf("unexpected status code: %d: %s", res.StatusCode(), body)
			return
		}

		splits := strings.Split(body, ":")
		if len(splits) != 2 {
			err = fmt.Errorf("unexpected response: %s", body)
			return
		}

		version = strings.TrimSpace(splits[1])

		log.Println(version)

		var sv *semver.Version
		if sv, err = semver.NewVersion(version); err != nil {
			return
		}

		if !spec.VersionConstraints.Check(sv) {
			err = fmt.Errorf("unsupported version: %s", sv)
			return
		}
	}

	os, arch := spec.ConvertPlatform(remapOS, remapArch)

	u = fmt.Sprintf(
		"https://gosspublic.alicdn.com/ossutil/%s/ossutil-%s-%s-%s.zip",
		strings.TrimPrefix(version, "v"),
		version,
		os,
		arch,
	)

	return
}

func (t *toolchain) Support(spec activate_toolchain.Spec) bool {
	return spec.Name == "ossutil"
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
		File:           spec.VersionedName() + ".zip",
		DirectoryLevel: 1,
	}); err != nil {
		return
	}

	script = fmt.Sprintf(`
export PATH="%s:$PATH";
`, dir)
	return
}

func init() {
	activate_toolchain.AddToolchain(&toolchain{})
}
