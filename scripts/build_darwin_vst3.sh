#!/usr/bin/env bash
set -euo pipefail

bundle_root="${1:-dist/macos}"
plugin_name="${2:-vst3go}"
entrypoint="${3:-./cmd/vst3go-dylib}"
output_dir="${bundle_root}"
output_dylib="${output_dir}/${plugin_name}.dylib"
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "${script_dir}/.." && pwd)"

if [[ "$bundle_root" != /* ]]; then
  bundle_root="${repo_root}/${bundle_root}"
  output_dir="${bundle_root}"
  output_dylib="${output_dir}/${plugin_name}.dylib"
fi

host_os="$(go env GOHOSTOS)"
if [[ "$host_os" != "darwin" ]]; then
  echo "macOS builds require a macOS host with the Apple SDK/toolchain installed." >&2
  exit 1
fi

mkdir -p "$output_dir"
(
  cd "$repo_root"
  GOOS=darwin CGO_ENABLED=1 go build -buildmode=c-shared -o "$output_dylib" "$entrypoint"
)
cp "${output_dylib}.h" "${output_dir}/${plugin_name}.h"
echo "built macOS VST3 dylib at: $output_dylib"
