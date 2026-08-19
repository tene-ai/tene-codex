#!/bin/sh
set -eu

root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
version=${1:-dev}
stage=${2:-"$root/dist/tene-codex-$version"}

mkdir -p "$stage"
for path in .codex-plugin skills hooks scripts references schemas templates LICENSE NOTICE README.md README-KR.md; do
  cp -R "$root/$path" "$stage/"
done

for target in darwin/arm64 darwin/amd64 linux/arm64 linux/amd64; do
  os=${target%/*}
  arch=${target#*/}
  output="$stage/bin/$os-$arch/tene-workflow"
  mkdir -p "$(dirname "$output")"
  CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" go build -trimpath -ldflags "-s -w -X main.version=$version" -o "$output" "$root/cmd/tene-workflow"
done

(
  cd "$stage"
  find . -type f ! -name checksums.txt -print0 | sort -z | xargs -0 shasum -a 256 > checksums.txt
)

echo "$stage"
