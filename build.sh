#!/usr/bin/env bash

goOS=("linux" "darwin" "freebsd" "linux" "netbsd" "openbsd" "windows")
goArch=("amd64" "arm64" "ppc64" "ppc64le" "riscv64" "s390x")

for os in ${goOS[@]}; do
    for arch in ${goArch[@]}; do
        if [ "$os" = "windows" ]; then
            GOOS=$os GOARCH=$arch go build -o "dist/bottle-$os-$arch.exe"
            continue
        fi
        GOOS=$os GOARCH=$arch go build -o dist/bottle-$os-$arch
    done
done