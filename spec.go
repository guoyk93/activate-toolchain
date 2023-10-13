package activate_toolchain

import (
	"errors"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"runtime"
	"strings"
)

// Spec is a toolchain spec.
// It is used to identify a toolchain.
type Spec struct {
	// Name is the name of the toolchain.
	Name string
	// VersionRaw is the original version string, without any prefixing 'v'.
	VersionRaw string
	// Version is the version requested
	Version *semver.Version
	// VersionConstraints is the version constraints with tide.
	VersionConstraints *semver.Constraints
	// OS is the target operating system.
	// follows go convention, see https://go.dev/doc/install/source#environment for values.
	OS string
	// Arch is the target architecture.
	// follows go convention, see https://go.dev/doc/install/source#environment for values.
	Arch string
}

// VersionHasMinor returns whether the version has minor.
func (s Spec) VersionHasMinor() bool {
	return s.Version.Minor() > 0 ||
		s.Version.Patch() > 0 ||
		strings.HasPrefix(s.VersionRaw, fmt.Sprintf("%d.", s.Version.Major()))
}

// VersionHasPatch returns whether the version has patch.
func (s Spec) VersionHasPatch() bool {
	return s.Version.Patch() > 0 ||
		strings.HasPrefix(s.VersionRaw, fmt.Sprintf("%d.%d.", s.Version.Major(), s.Version.Minor()))
}

func (s Spec) VersionedName() string {
	return s.Name + "-" + s.VersionRaw
}

// ParseSpec parses a spec from string.
func ParseSpec(s string) (spec Spec, err error) {
	splits := strings.Split(s, "@")
	if len(splits) != 2 {
		err = errors.New("invalid spec '" + s + "': must be in format of 'NAME@VERSION'")
		return
	}

	spec.Name = strings.TrimSpace(splits[0])
	spec.VersionRaw = strings.TrimPrefix(
		strings.TrimSpace(splits[1]),
		"v",
	)

	if spec.Version, err = semver.NewVersion(spec.VersionRaw); err != nil {
		err = errors.New("invalid spec '" + s + "': " + err.Error())
		return
	}

	if spec.VersionConstraints, err = semver.NewConstraint("~" + spec.VersionRaw); err != nil {
		err = errors.New("invalid spec '" + s + "': " + err.Error())
		return
	}

	spec.OS = runtime.GOOS
	spec.Arch = runtime.GOARCH

	return
}
