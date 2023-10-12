package jdk

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/PuerkitoBio/goquery"
	"github.com/guoyk93/activate-toolchain"
	"log"
	"net/url"
	"strconv"
)

type sourceUSTC struct {
}

func (u sourceUSTC) Name() string {
	return "ustc.edu.cn"
}

func (u sourceUSTC) Primary() bool {
	return false
}

func (u sourceUSTC) Resolve(ctx context.Context, version *semver.Version, os, arch string) (out string, err error) {
	if version.Original() != strconv.Itoa(int(version.Major())) {
		err = errors.New("only support major version")
		return
	}

	baseURL := fmt.Sprintf(
		"https://mirrors.ustc.edu.cn/adoptium/releases/temurin%d-binaries/LatestRelease/",
		version.Major(),
	)

	if err = activate_toolchain.FetchQueryHTML(
		ctx,
		baseURL,
		"a",
		func(i int, s *goquery.Selection) {
			href, _ := s.Attr("href")

			if !validateAdoptiumFilename(href, os, arch) {
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
