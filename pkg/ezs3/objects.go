package ezs3

import (
	"context"
	"encoding/xml"
	"errors"
	"github.com/go-resty/resty/v2"
	"strconv"
	"strings"
)

type ListBucketContent struct {
	Key  string `xml:"Key"`
	ETag string `xml:"ETag"`
	Size int64  `xml:"Size"`
}

type ListBucketCommonPrefix struct {
	Prefix string `xml:"Prefix"`
}

type ListBucketResult struct {
	XMLName xml.Name `xml:"ListBucketResult"`

	Name           string                   `xml:"Name"`
	Prefix         string                   `xml:"Prefix"`
	Marker         string                   `xml:"Marker"`
	NextMarker     string                   `xml:"NextMarker"`
	IsTruncated    bool                     `xml:"IsTruncated"`
	Contents       []ListBucketContent      `xml:"Contents"`
	CommonPrefixes []ListBucketCommonPrefix `xml:"CommonPrefixes"`
}

type ListObjectsOptions struct {
	Endpoint  string
	Prefix    string
	Marker    string
	Delimiter string
	MaxKeys   int64
}

func (o ListObjectsOptions) AsMap() map[string]string {
	args := map[string]string{}
	if o.Prefix != "" {
		args["prefix"] = o.Prefix
	}
	if o.Marker != "" {
		args["marker"] = o.Marker
	}
	if o.Delimiter != "" {
		args["delimiter"] = o.Delimiter
	}
	if o.MaxKeys > 0 {
		args["max-keys"] = strconv.FormatInt(o.MaxKeys, 10)
	}
	return args
}

// ListObjects list objects from s3 compatible storage
func ListObjects(ctx context.Context, opts ListObjectsOptions) (contents []ListBucketContent, prefixes []ListBucketCommonPrefix, err error) {
	if !strings.HasSuffix(opts.Endpoint, "/") {
		opts.Endpoint = opts.Endpoint + "/"
	}

	client := resty.New()

	for {
		var (
			res  *resty.Response
			data ListBucketResult
		)
		if res, err = client.R().
			SetContext(ctx).
			SetQueryParams(opts.AsMap()).
			SetResult(&data).
			Get(opts.Endpoint); err != nil {
			return
		}
		if res.IsError() {
			err = errors.New("failed to list objects: " + res.String())
			return
		}
		for _, item := range data.Contents {
			contents = append(contents, item)
		}
		for _, item := range data.CommonPrefixes {
			prefixes = append(prefixes, item)
		}

		if !data.IsTruncated {
			break
		}

		opts.Marker = data.NextMarker
	}

	return

}
