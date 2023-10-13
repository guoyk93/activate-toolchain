package jdk

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/go-resty/resty/v2"
	"github.com/guoyk93/activate-toolchain"
	"log"
	"net/url"
	"strconv"
)

type sourceAdoptium struct {
}

func (u sourceAdoptium) Name() string {
	return "adoptium"
}

func (u sourceAdoptium) Primary() bool {
	return true
}

func (u sourceAdoptium) ResolveDownloadURL(ctx context.Context, spec activate_toolchain.Spec) (downloadURL string, err error) {
	os, arch := convertPlatform(spec)

	client := resty.New()

	type ReleaseVersion struct {
		OpenJDKVersion string `json:"openjdk_version"`
		Semver         string `json:"semver"`
	}

	type ResponseReleaseVersions struct {
		Versions []ReleaseVersion `json:"versions"`
	}

	type ResponseReleaseNames struct {
		Releases []string `json:"releases"`
	}

	// find the best release version
	var bestReleaseVersion ReleaseVersion
	{
		var items []ReleaseVersion

		var (
			page     int
			pageSize = 20
		)

		for {
			var (
				res  *resty.Response
				data ResponseReleaseVersions
			)
			if res, err = client.R().SetContext(ctx).SetQueryParams(map[string]string{
				"architecture": arch,
				"heap_size":    "normal",
				"image_type":   "jdk",
				"os":           os,
				"page":         strconv.Itoa(page),
				"page_size":    strconv.Itoa(pageSize),
				"project":      "jdk",
				"release_type": "ga",
				"sort_method":  "DEFAULT",
				"sort_order":   "DESC",
				"vendor":       "eclipse",
				"jvm_impl":     "hotspot",
				"version":      fmt.Sprintf("[%d, %d)", spec.Version.Major(), spec.Version.Major()+1),
			}).SetResult(&data).Get("https://api.adoptium.net/v3/info/release_versions"); err != nil {
				return
			}
			if res.IsError() {
				err = errors.New("failed fetching release versions: " + res.String())
				return
			}

			for _, item := range data.Versions {
				items = append(items, item)
			}

			if len(data.Versions) < pageSize {
				break
			}

			page++
		}

		if bestReleaseVersion, err = activate_toolchain.FindBestVersionedItem(
			spec.VersionConstraints,
			items,
			func(v ReleaseVersion) (version *semver.Version, err error) {
				return semver.NewVersion(v.Semver)
			}); err != nil {
			return
		}
	}

	var releaseName string

	// get release name from semver
	{

		var (
			res  *resty.Response
			data ResponseReleaseNames
		)
		if res, err = client.R().SetContext(ctx).SetQueryParams(map[string]string{
			"architecture": arch,
			"heap_size":    "normal",
			"image_type":   "jdk",
			"os":           os,
			"page":         "0",
			"page_size":    "20",
			"project":      "jdk",
			"release_type": "ga",
			"sort_method":  "DEFAULT",
			"sort_order":   "DESC",
			"vendor":       "eclipse",
			"jvm_impl":     "hotspot",
			"semver":       "true",
			"version":      bestReleaseVersion.Semver,
		}).SetResult(&data).Get("https://api.adoptium.net/v3/info/release_names"); err != nil {
			return
		}
		if res.IsError() {
			err = errors.New("failed fetching release names: " + res.String())
			return
		}

		if len(data.Releases) == 0 {
			err = errors.New("no release found")
			return
		}

		releaseName = data.Releases[0]
	}

	log.Println("found matched release:", releaseName)

	downloadURL = fmt.Sprintf(
		"https://api.adoptium.net/v3/binary/version/%s/%s/%s/jdk/hotspot/normal/eclipse?project=jdk",
		url.PathEscape(releaseName),
		os,
		arch,
	)

	return
}

func init() {
	addSource(sourceAdoptium{})
}
