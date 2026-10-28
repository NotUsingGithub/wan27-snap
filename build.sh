#!/bin/sh
# Build static Linux binaries (amd64, arm64) and pack into .snap packages using mksquashfs.
set -e
V=0.1.0
for a in amd64 arm64; do
  mkdir -p build/$a/meta build/$a/bin
  sed "s/ARCH/$a/" snap.yaml > build/$a/meta/snap.yaml
  CGO_ENABLED=0 GOOS=linux GOARCH=$a go build -ldflags="-s -w" -o build/$a/bin/wan27 .
  chmod 755 build/$a/bin/wan27
  cp build/$a/bin/wan27 build/wan27-linux-$a
done

if command -v mksquashfs >/dev/null 2>&1; then
  for a in amd64 arm64; do
    mksquashfs build/$a build/wan27_${V}_${a}.snap -noappend -comp xz -all-root -no-xattrs -no-fragments
  done
  ls -la build/*.snap
fi
