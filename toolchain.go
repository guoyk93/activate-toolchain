package activate_toolchain

import (
	"context"
	"sync"
)

// Toolchain is an interface abstracting a toolchain.
type Toolchain interface {
	// Support returns whether the toolchain supports the spec.
	Support(spec Spec) bool

	// Activate download, install and activate the toolchain
	Activate(ctx context.Context, spec Spec) (script string, err error)
}

var (
	toolchains     []Toolchain
	toolchainsLock sync.Locker = &sync.Mutex{}
)

// AddToolchain adds a toolchain to the list.
func AddToolchain(t Toolchain) {
	toolchainsLock.Lock()
	defer toolchainsLock.Unlock()
	toolchains = append(toolchains, t)
}

// FindToolchain finds a toolchain by spec.
func FindToolchain(spec Spec) (Toolchain, bool) {
	toolchainsLock.Lock()
	defer toolchainsLock.Unlock()
	for _, t := range toolchains {
		if t.Support(spec) {
			return t, true
		}
	}
	return nil, false
}
