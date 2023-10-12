package activate_toolchain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// TryFetch tries to fetch url, returns delay and error
func TryFetch(ctx context.Context, url string) (delay time.Duration, err error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var req *http.Request
	if req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil); err != nil {
		return
	}

	start := time.Now()

	var res *http.Response
	if res, err = http.DefaultClient.Do(req); err != nil {
		return
	}
	defer res.Body.Close()

	delay = time.Since(start)

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status code: %d while fetching %s", res.StatusCode, url)
		return
	}
	return
}

// DetectFastestURL detects the fastest url from urls
func DetectFastestURL(ctx context.Context, urls []string) (fastest string, err error) {
	var (
		delays  = map[string]time.Duration{}
		delaysL = &sync.Mutex{}
		wg      = &sync.WaitGroup{}
	)

	for _, _url := range urls {
		url := _url
		wg.Add(1)

		go func() {
			defer wg.Done()
			if dur, err := TryFetch(ctx, url); err == nil {
				delaysL.Lock()
				delays[url] = dur
				delaysL.Unlock()
			} else {
				log.Println("failed to fetch:", url, err)
			}
		}()
	}

	wg.Wait()

	for url, delay := range delays {
		if fastest == "" || delay < delays[fastest] {
			fastest = url
		}
	}

	if fastest == "" {
		err = errors.New("no fastest url available, all failed")
		return
	}

	return
}

func FetchFile(ctx context.Context, url string, localFile string) (err error) {
	var f *os.File
	if f, err = os.OpenFile(localFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644); err != nil {
		return
	}
	defer f.Close()

	var req *http.Request
	if req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil); err != nil {
		return
	}

	var res *http.Response
	if res, err = http.DefaultClient.Do(req); err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status code: %d while fetching %s", res.StatusCode, url)
		return
	}

	if _, err = io.Copy(f, res.Body); err != nil {
		return
	}

	return
}

// AdvancedFetchFile detect fastest url from candidate urls and fetch file to local path
func AdvancedFetchFile(ctx context.Context, urls []string, localFile string) (err error) {
	if len(urls) == 0 {
		err = errors.New("no urls provided")
		return
	}

	var fastest string
	if fastest, err = DetectFastestURL(ctx, urls); err != nil {
		return
	}

	log.Println("fastest url:", fastest)

	tmpPath := localFile + ".tmp"

	if err = FetchFile(ctx, fastest, tmpPath); err != nil {
		return
	}

	_ = os.RemoveAll(localFile)

	if err = os.Rename(tmpPath, localFile); err != nil {
		return
	}

	return
}

// FetchJSON fetch json from url and unmarshal to v
func FetchJSON(ctx context.Context, url string, v interface{}) (err error) {
	var req *http.Request
	if req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil); err != nil {
		return
	}
	var res *http.Response
	if res, err = http.DefaultClient.Do(req); err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status code: %d while fetching %s", res.StatusCode, url)
		return
	}
	if err = json.NewDecoder(res.Body).Decode(v); err != nil {
		err = fmt.Errorf("failed to parse json: %w while fetching %s", err, url)
		return
	}
	return
}

// FetchQueryHTML fetches html from url, and find elements by selector, and call fn for each element
func FetchQueryHTML(ctx context.Context, url string, sel string, fn func(i int, s *goquery.Selection)) (err error) {
	var req *http.Request
	if req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil); err != nil {
		return
	}
	var res *http.Response
	if res, err = http.DefaultClient.Do(req); err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status code: %d while fetching %s", res.StatusCode, url)
		return
	}
	var doc *goquery.Document
	if doc, err = goquery.NewDocumentFromReader(res.Body); err != nil {
		err = fmt.Errorf("failed to parse html: %w while fetching %s", err, url)
		return
	}
	doc.Find(sel).Each(fn)
	return
}
