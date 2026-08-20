#!/bin/sh
set -eu
root=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
tmp=$(mktemp -d "${TMPDIR:-/tmp}/tene-release-smoke.XXXXXX")
trap 'rm -rf "$tmp"' EXIT HUP INT TERM
stage="$tmp/tene-codex-0.1.1"
PACKAGE_TARGETS="$(uname -s | tr '[:upper:]' '[:lower:]')/$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')" "$root/scripts/package-plugin.sh" 0.1.1 "$stage" >/dev/null
python3 -m json.tool "$stage/.codex-plugin/plugin.json" >/dev/null
python3 -m json.tool "$stage/sbom.spdx.json" >/dev/null
echo "stage package-manifest-sbom passed"
PLUGIN_ROOT="$stage" "$stage/scripts/tene-workflow" version --json >/dev/null
echo "stage bundled-cli passed"
test "$(PLUGIN_ROOT="$stage" "$stage/scripts/tene-workflow" route --active false --text 'use \$tene:qa' --json | python3 -c 'import json,sys;print(json.load(sys.stdin)["result"]["selected_skill"])')" = qa
test "$(PLUGIN_ROOT="$stage" "$stage/scripts/tene-workflow" route --active true --phase qa --text 'Playwright UX QA 테스트해줘' --json | python3 -c 'import json,sys;print(json.load(sys.stdin)["result"]["selected_skill"])')" = qa
echo "stage routing-explicit-implicit passed"
platform="$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')"
binary="$stage/bin/$platform/tene-workflow";cp "$binary" "$tmp/binary.backup";printf x >> "$binary"
if PLUGIN_ROOT="$stage" "$stage/scripts/tene-workflow" version >/dev/null 2>&1; then echo "tampered binary was accepted" >&2;exit 1;fi
mv "$tmp/binary.backup" "$binary"
echo "stage tampered-binary-rejection passed"
project="$tmp/project";mkdir -p "$project";PLUGIN_ROOT="$stage" "$stage/scripts/tene-workflow" --root "$project" init --name preserve >/dev/null
PLUGIN_ROOT="$stage" python3 "$root/scripts/portable-workflow-smoke.py" --cli "$stage/scripts/tene-workflow" --workspace "$tmp/portable-matrix" >/dev/null
echo "stage portable-workflow-matrix passed"
cp -R "$stage" "$tmp/update";PLUGIN_ROOT="$tmp/update" "$tmp/update/scripts/tene-workflow" --root "$project" status --json >/dev/null
rm -rf "$tmp/update" "$stage"
echo "stage update-remove-simulation passed"
test -f "$project/.tene-workflow/project.json"
echo "stage project-state-preservation passed"
echo "release smoke passed"
