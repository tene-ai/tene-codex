#!/bin/sh
set -eu
root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
tmp=$(mktemp -d "${TMPDIR:-/tmp}/tene-release-smoke.XXXXXX")
trap 'rm -rf "$tmp"' EXIT HUP INT TERM
stage="$tmp/tene-codex-0.1.0"
PACKAGE_TARGETS="$(uname -s | tr '[:upper:]' '[:lower:]')/$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')" "$root/scripts/package-plugin.sh" 0.1.0 "$stage" >/dev/null
python3 -m json.tool "$stage/.codex-plugin/plugin.json" >/dev/null
python3 -m json.tool "$stage/sbom.spdx.json" >/dev/null
PLUGIN_ROOT="$stage" "$stage/scripts/tene-workflow" version --json >/dev/null
test "$(PLUGIN_ROOT="$stage" "$stage/scripts/tene-workflow" route --active false --text 'use \$tene-qa' --json | python3 -c 'import json,sys;print(json.load(sys.stdin)["result"]["selected_skill"])')" = tene-qa
test "$(PLUGIN_ROOT="$stage" "$stage/scripts/tene-workflow" route --active true --phase qa --text 'Playwright UX QA 테스트해줘' --json | python3 -c 'import json,sys;print(json.load(sys.stdin)["result"]["selected_skill"])')" = tene-qa
platform="$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')"
binary="$stage/bin/$platform/tene-workflow";cp "$binary" "$tmp/binary.backup";printf x >> "$binary"
if PLUGIN_ROOT="$stage" "$stage/scripts/tene-workflow" version >/dev/null 2>&1; then echo "tampered binary was accepted" >&2;exit 1;fi
mv "$tmp/binary.backup" "$binary"
project="$tmp/project";mkdir -p "$project";PLUGIN_ROOT="$stage" "$stage/scripts/tene-workflow" --root "$project" init --name preserve >/dev/null
cp -R "$stage" "$tmp/update";PLUGIN_ROOT="$tmp/update" "$tmp/update/scripts/tene-workflow" --root "$project" status --json >/dev/null
rm -rf "$tmp/update" "$stage"
test -f "$project/.tene-workflow/project.json"
echo "release smoke passed"
