package activate_toolchain

import "context"

// Toolchains is a list of all toolchains.
var Toolchains []Toolchain

// Toolchain is an interface abstracting a toolchain.
type Toolchain interface {
	// Support returns whether the toolchain supports the spec.
	Support(spec Spec) bool

	// Activate download, install and activate the toolchain
	Activate(ctx context.Context, spec Spec) (script string, err error)
}
