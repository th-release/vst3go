#!/usr/bin/env bash
set -euo pipefail

bundle_root="${1:-dist/windows}"
plugin_name="${2:-vst3go}"
output_dir="${bundle_root}/${plugin_name}.vst3"
output_dll="${output_dir}/Contents/x86_64-win/${plugin_name}.vst3"

host_os="$(go env GOHOSTOS)"
cc_bin="${CC:-cc}"

if [[ "$host_os" != "windows" ]]; then
  if ! printf '#include <windows.h>\n' | "$cc_bin" -E -x c - >/dev/null 2>&1; then
    echo "Windows build toolchain is not ready: need a compiler that can find <windows.h> (for example MinGW-w64)." >&2
    echo "Set CC to a Windows-capable compiler before running this script again." >&2
    exit 1
  fi
fi

mkdir -p "${output_dir}/Contents/x86_64-win"
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -buildmode=c-shared -o "$output_dll" ./cmd/vst3go-dll
cp "${output_dll}.h" "${output_dir}/Contents/x86_64-win/${plugin_name}.h"
echo "built Windows VST3 DLL at: $output_dll"
