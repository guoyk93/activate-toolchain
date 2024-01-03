package jdk

import (
	"context"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/yankeguo/activate-toolchain"
	"log"
	"net/url"
)

type sourceUSTC struct {
}

func (u sourceUSTC) Name() string {
	return "ustc.edu.cn"
}

func (u sourceUSTC) Primary() bool {
	return false
}

func (u sourceUSTC) ResolveDownloadURL(ctx context.Context, spec activate_toolchain.Spec) (out string, err error) {
	if spec.VersionHasMinor() {
		err = errors.New("sourceUSTC: only support major version")
		return
	}

	baseURL := fmt.Sprintf(
		"https://mirrors.ustc.edu.cn/adoptium/releases/temurin%d-binaries/LatestRelease/",
		spec.Version.Major(),
	)

	if err = activate_toolchain.FetchQueryHTML(
		ctx,
		baseURL,
		"a",
		func(i int, s *goquery.Selection) {
			href, _ := s.Attr("href")

			if !validateAdoptiumFilename(href, spec) {
				return
			}

			log.Println(href)

			out, _ = url.JoinPath(baseURL, href)
		},
	); err != nil {
		return
	}

	if out == "" {
		err = fmt.Errorf("no release found")
		return
	}

	return
}

func init() {
	addSource(sourceUSTC{})
}
