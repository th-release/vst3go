#!/usr/bin/env bash
set -euo pipefail

bundle_root="${1:-dist/windows}"
plugin_name="${2:-vst3go}"
output_dir="${bundle_root}/${plugin_name}.vst3"
output_dll="${output_dir}/Contents/x86_64-win/${plugin_name}.vst3"

mkdir -p "${output_dir}/Contents/x86_64-win"
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -buildmode=c-shared -o "$output_dll" ./cmd/vst3go-dll
cp "${output_dll}.h" "${output_dir}/Contents/x86_64-win/${plugin_name}.h"
echo "built Windows VST3 DLL at: $output_dll"
