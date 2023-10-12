package activate_toolchain

import "context"

// Toolchains is a list of all toolchains.
var Toolchains []Toolchain

// Toolchain is an interface abstracting a toolchain.
type Toolchain interface {
	// Name returns the name of the toolchain.
	Name() string

	// Activate download, install and activate the toolchain with specified version, os and arch.
	// 'targetVersion' is the version to resolve, can be major version, major.minor version, or full version.
	// 'os' and 'arch' are the target operating system and architecture, follows go convention, see https://go.dev/doc/install/source#environment for values.
	// returns the shell eval script string to activate the toolchain, or error if any.
	Activate(ctx context.Context, targetVersion, os, arch string) (script string, err error)
}
