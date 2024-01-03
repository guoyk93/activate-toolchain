package maven

import (
	"context"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/PuerkitoBio/goquery"
	"github.com/yankeguo/activate-toolchain"
	"github.com/yankeguo/activate-toolchain/pkg/ezscript"
	"log"
	"regexp"
	"strings"
)

var (
	urlsArchive = []string{
		"https://archive.apache.org/dist/maven",
	}
	urlsLatest = []string{
		"https://dlcdn.apache.org/maven",
		"https://mirrors.cloud.tencent.com/apache/maven",
		"https://mirrors.aliyun.com/apache/maven",
		"https://mirrors.ustc.edu.cn/apache/maven",
	}
)

type toolchain struct {
}

func (t toolchain) Support(spec activate_toolchain.Spec) bool {
	return spec.Name == "maven"
}

func (t toolchain) resolveURL(ctx context.Context, spec activate_toolchain.Spec) (downloadURL string, err error) {
	var (
		forceArchiveSearch bool
		archiveSearched    bool
	)

retry:

	var baseURL string

	if spec.VersionHasPatch() || forceArchiveSearch {
		archiveSearched = true
		if baseURL, err = activate_toolchain.DetectFastestURL(ctx, urlsArchive); err != nil {
			return
		}
	} else {
		if baseURL, err = activate_toolchain.DetectFastestURL(ctx, urlsLatest); err != nil {
			return
		}
	}

	var (
		versions []string

		regexpHref = regexp.MustCompile(`^\d+\.\d+\.\d+.*/$`)
	)

	if err = activate_toolchain.FetchQueryHTML(
		ctx,
		fmt.Sprintf("%s/maven-%d", baseURL, spec.Version.Major()),
		"a",
		func(i int, s *goquery.Selection) {
			href, _ := s.Attr("href")
			if !regexpHref.MatchString(href) {
				return
			}
			versions = append(versions, strings.TrimSuffix(href, "/"))
		},
	); err != nil {
		return
	}

	if len(versions) == 0 {
		err = fmt.Errorf("no version found in page")
		return
	}

	var bestVersion string

	if bestVersion, err = activate_toolchain.FindBestVersionedItem(spec.VersionConstraints, versions, semver.NewVersion); err != nil {

		// if no matching version found in latest site, force try archive site
		if !archiveSearched {
			forceArchiveSearch = true
			log.Println("no matching version found in latest site, force try archive site")
			goto retry
		}

		return
	}

	downloadURL = fmt.Sprintf(
		"%s/maven-%d/%s/binaries/apache-maven-%s-bin.tar.gz",
		baseURL,
		spec.Version.Major(),
		bestVersion,
		bestVersion,
	)
	return
}

func (t toolchain) Activate(ctx context.Context, spec activate_toolchain.Spec) (script string, err error) {
	var dir string

	if dir, err = activate_toolchain.InstallArchive(ctx, activate_toolchain.InstallArchiveOptions{
		ProvideURLs: func() (urls []string, err error) {
			var downloadURL string
			if downloadURL, err = t.resolveURL(ctx, spec); err != nil {
				return
			}
			urls = append(urls, downloadURL)
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
	activate_toolchain.AddToolchain(toolchain{})
}
