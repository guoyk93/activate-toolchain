package main

import (
	"context"
	"errors"
	"github.com/guoyk93/activate-toolchain"
	"log"
	"os"
	"strings"

	_ "github.com/guoyk93/activate-toolchain/toolchains/jdk"
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
		var spec activate_toolchain.Spec

		if spec, err = activate_toolchain.ParseSpec(arg); err != nil {
			return
		}

		for _, toolchain := range activate_toolchain.Toolchains {
			if !toolchain.Support(spec) {
				continue
			}

			var script string
			if script, err = toolchain.Activate(ctx, spec); err != nil {
				return
			}

			scripts = append(scripts, script)

			continue argLoop
		}

		err = errors.New("no supported toolchain: " + arg)
		return
	}

	if _, err = os.Stdout.WriteString(strings.Join(scripts, "\n\n")); err != nil {
		return
	}
	_ = os.Stdout.Sync()
}
