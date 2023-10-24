package activate_toolchain

import (
	"errors"
	"github.com/Masterminds/semver/v3"
	"sort"
)

// FindBestVersionedItem finds the best match version item from a list of items.
// If there is a match, the returned matched item is the best match item.
// The returned error is nil if there is a match, otherwise it is not nil.
// nil value and failed value in input items will be ignored.
func FindBestVersionedItem[T any](c *semver.Constraints, items []T, fn func(v T) (version *semver.Version, err error)) (matched T, err error) {
	var versions semver.Collection

	for _, item := range items {
		var version *semver.Version
		if version, err = fn(item); err != nil {
			version = nil
		}
		versions = append(versions, version)
	}

	var idx int
	if idx, _, err = FindBestVersion(c, versions); err != nil {
		return
	}

	matched = items[idx]
	return
}

// FindBestVersion finds the best match version from a list of versions.
// The returned index is the index of the version in the original list.
// The returned version is the best match version.
// The returned error is nil if there is a match, otherwise it is not nil.
// nil value in input versions will be ignored.
func FindBestVersion(c *semver.Constraints, versions semver.Collection) (idx int, version *semver.Version, err error) {
	var matched semver.Collection

	for _, v := range versions {
		if v == nil {
			continue
		}
		if c.Check(v) {
			matched = append(matched, v)
		}
	}

	sort.Sort(sort.Reverse(matched))

	if len(matched) == 0 {
		err = errors.New("no matching version")
		return
	}

	version = matched[0]

	for i, v := range versions {
		if v == nil {
			continue
		}
		if v.Equal(version) {
			idx = i
			return
		}
	}

	return
}
