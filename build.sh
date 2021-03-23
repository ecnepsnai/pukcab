#!/bin/bash

function build {
    GOOS=$1
    GOARCH=$2
    NAME="pukcab_${GOOS}_${GOARCH}.tar.gz"

    rm -f pukcab
    go build -ldflags="-s -w"
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
