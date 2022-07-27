#!/bin/bash

VERSION=${1:?Version required}

function build {
    GOOS=$1
    GOARCH=$2
    CGO_ENABLED=0
    NAME="pukcab-${VERSION}_${GOOS}_${GOARCH}.tar.gz"

    rm -f pukcab
    go build -ldflags="-s -w" -trimpath -buildmode=exe
    tar -czf ${NAME} pukcab
    rm -f pukcab
    mv ${NAME} ../../artifacts/
}

rm -rf artifacts
mkdir -p artifacts
cd cmd/pukcab

for ARCH in 'amd64' 'arm64'; do
    for OS in 'linux' 'freebsd' 'openbsd' 'netbsd' 'darwin'; do
        build ${OS} ${ARCH}
    done
done
