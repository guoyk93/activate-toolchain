package main

import (
	"context"
	"errors"
	"github.com/guoyk93/activate-toolchain"
	"log"
	"os"
	"runtime"
	"strings"

	_ "github.com/guoyk93/activate-toolchain/toolchains/node"
)

func main() {
	var err error
	defer func() {
		if err == nil {
			return
		}
		log.Println("exited with error:", err.Error())
		os.Exit(1)
	}()

	ctx := context.Background()

	var scripts []string

argLoop:
	for _, arg := range os.Args[1:] {
		splits := strings.Split(arg, "@")
		if len(splits) != 2 {
			err = errors.New("invalid argument: " + arg + ", must be 'NAME@VERSION'")
			return
		}
		name := strings.TrimSpace(splits[0])
		version := strings.TrimSpace(splits[1])

		if strings.HasPrefix(version, "v") {
			version = version[1:]
		}

		for _, toolchain := range activate_toolchain.Toolchains {
			if toolchain.Name() != name {
				continue
			}
			var script string
			if script, err = toolchain.Activate(ctx, version, runtime.GOOS, runtime.GOARCH); err != nil {
				return
			}

			scripts = append(scripts, script)

			continue argLoop
		}

		err = errors.New("toolchain not supported: " + name)
		return
	}

	if _, err = os.Stdout.WriteString(strings.Join(scripts, "\n\n")); err != nil {
		return
	}
	_ = os.Stdout.Sync()
}
