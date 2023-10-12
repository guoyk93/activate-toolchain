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

func (u sourceAdoptium) Resolve(ctx context.Context, version *semver.Version, os, arch string) (out string, err error) {
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

	// key is 'semver'
	releaseVersions := map[string]ReleaseVersion{}

	// fetch all release versions
	{
		var (
			page     int
			pageSize = 20
		)
		for {
			var (
				res  *resty.Response
				data ResponseReleaseVersions
			)
			if res, err = client.R().SetQueryParams(map[string]string{
				"architecture": remapArch(arch),
				"heap_size":    "normal",
				"image_type":   "jdk",
				"os":           remapOS(os),
				"page":         strconv.Itoa(page),
				"page_size":    strconv.Itoa(pageSize),
				"project":      "jdk",
				"release_type": "ga",
				"sort_method":  "DEFAULT",
				"sort_order":   "DESC",
				"vendor":       "eclipse",
				"jvm_impl":     "hotspot",
				"version":      fmt.Sprintf("[%d, %d)", version.Major(), version.Major()+1),
			}).SetResult(&data).Get("https://api.adoptium.net/v3/info/release_versions"); err != nil {
				return
			}
			if res.IsError() {
				err = errors.New("failed fetching release versions: " + res.String())
				return
			}

			for _, item := range data.Versions {
				releaseVersions[item.Semver] = item
			}

			if len(data.Versions) < pageSize {
				break
			}

			page++
		}
	}

	var bestSemver string

	// calculate the best semver
	{
		var semvers []string

		for sv := range releaseVersions {
			semvers = append(semvers, sv)
		}

		if bestSemver, err = activate_toolchain.ResolveVersion(semvers, version.Original()); err != nil {
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
		if res, err = client.R().SetQueryParams(map[string]string{
			"architecture": remapArch(arch),
			"heap_size":    "normal",
			"image_type":   "jdk",
			"os":           remapOS(os),
			"page":         "0",
			"page_size":    "20",
			"project":      "jdk",
			"release_type": "ga",
			"sort_method":  "DEFAULT",
			"sort_order":   "DESC",
			"vendor":       "eclipse",
			"jvm_impl":     "hotspot",
			"semver":       "true",
			"version":      bestSemver,
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

	out = fmt.Sprintf(
		"https://api.adoptium.net/v3/binary/version/%s/%s/%s/jdk/hotspot/normal/eclipse?project=jdk",
		url.PathEscape(releaseName),
		remapOS(os),
		remapArch(arch),
	)

	return
}

func init() {
	addSource(sourceAdoptium{})
}
