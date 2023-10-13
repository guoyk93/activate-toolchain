package activate_toolchain

import (
	"errors"
	"github.com/Masterminds/semver/v3"
	"sort"
)

// FindBestVersionedItem finds the best match version item from a list of items.
func FindBestVersionedItem[T any](c *semver.Constraints, items []T, fn func(v T) (version *semver.Version, err error)) (matched T, err error) {
	var versions semver.Collection

	for _, item := range items {
		var version *semver.Version
		if version, err = fn(item); err != nil {
			return
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
func FindBestVersion(c *semver.Constraints, versions semver.Collection) (idx int, version *semver.Version, err error) {
	var matched semver.Collection
	for _, v := range versions {
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
		if v.Equal(version) {
			idx = i
			return
		}
	}

	return
}
