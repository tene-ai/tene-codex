#!/bin/sh
set -eu

root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
version=${1:-dev}
stage=${2:-"$root/dist/tene-codex-$version"}

if [ -e "$stage" ]; then
  echo "refusing to package over existing stage: $stage" >&2
  exit 2
fi
mkdir -p "$stage"
for path in .codex-plugin skills hooks scripts references schemas templates LICENSE NOTICE README.md README-KR.md CHANGELOG.md PRIVACY.md TERMS.md SECURITY.md SUPPORT.md; do
  cp -R "$root/$path" "$stage/"
done

for target in ${PACKAGE_TARGETS:-"darwin/arm64 darwin/amd64 linux/arm64 linux/amd64"}; do
  os=${target%/*}
  arch=${target#*/}
  output="$stage/bin/$os-$arch/tene-workflow"
  mkdir -p "$(dirname "$output")"
  CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" go build -trimpath -ldflags "-s -w -X main.version=$version" -o "$output" "$root/cmd/tene-workflow"
done

python3 "$root/scripts/generate-sbom.py" "$stage" "$version"

(
  cd "$stage"
  find . -type f ! -name checksums.txt -print0 | sort -z | xargs -0 shasum -a 256 > checksums.txt
)

echo "$stage"
