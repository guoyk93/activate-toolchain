package jdk

import (
	"context"
	"github.com/Masterminds/semver/v3"
	"sync"
)

type Source interface {
	Name() string

	Primary() bool

	Resolve(ctx context.Context, version *semver.Version, os, arch string) (out string, err error)
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
