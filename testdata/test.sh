#!/bin/bash

set -eu

cd "$(dirname "${0}")"

go run ../cmd/activate-toolchain

eval "$(go run ../cmd/activate-toolchain)"

which java

java -version

which mvn

mvn --version

which node

node --version

which ossutil

ossutil --version

which pnpm

pnpm --version