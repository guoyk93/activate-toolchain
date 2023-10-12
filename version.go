package activate_toolchain

import (
	"errors"
	"github.com/Masterminds/semver/v3"
	"sort"
)

// ResolveVersion resolves the best match version from a list of versions.
func ResolveVersion(versions []string, target string) (result string, err error) {
	var c *semver.Constraints
	if c, err = semver.NewConstraint("~" + target); err != nil {
		return
	}

	var svs semver.Collection

	for _, version := range versions {
		var v *semver.Version
		if v, err = semver.NewVersion(version); err != nil {
			err = nil
			continue
		}
		if c.Check(v) {
			svs = append(svs, v)
		}
	}

	sort.Sort(sort.Reverse(svs))

	if len(svs) == 0 {
		err = errors.New("no matching version for " + target)
		return
	}

	result = svs[0].Original()
	return
}
