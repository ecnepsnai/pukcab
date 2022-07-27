#!/bin/bash

VERSION=${1:?Version required}

function build {
    export GOOS=$1
    export GOARCH=$2
    export CGO_ENABLED=0
    export NAME="pukcab-${VERSION}_${GOOS}_${GOARCH}.tar.gz"

    rm -f pukcab
    go build -ldflags="-s -w" -trimpath -buildmode=exe
    tar -czf ${NAME} pukcab
    rm -f pukcab
    mv ${NAME} ../../artifacts/
    echo ${NAME}
}

rm -rf artifacts
mkdir -p artifacts
cd cmd/pukcab

for ARCH in 'amd64' 'arm64'; do
    for OS in 'linux' 'freebsd' 'openbsd' 'netbsd' 'darwin'; do
        build ${OS} ${ARCH}
    done
done
