package jdk

import (
	"context"
	"github.com/yankeguo/activate-toolchain"
	"sync"
)

// Source is a source of JDK
type Source interface {
	// Name returns the name of the source
	Name() string

	// Primary returns if the source is primary
	Primary() bool

	// ResolveDownloadURL resolves the JDK download url
	ResolveDownloadURL(ctx context.Context, spec activate_toolchain.Spec) (downloadURL string, err error)
}

var (
	sources     []Source
	sourcesLock sync.Locker = &sync.Mutex{}
)

func addSource(src Source) {
	sourcesLock.Lock()
	defer sourcesLock.Unlock()

	sources = append(sources, src)
}
