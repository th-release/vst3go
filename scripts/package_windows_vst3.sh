#!/usr/bin/env bash
set -euo pipefail

if [[ $# -lt 1 ]]; then
  echo "usage: $0 <plugin-dll> [bundle-root] [plugin-name]" >&2
  exit 1
fi

plugin_dll="$1"
bundle_root="${2:-dist/windows}"
plugin_name="${3:-vst3go}"

if [[ ! -f "$plugin_dll" ]]; then
  echo "plugin DLL not found: $plugin_dll" >&2
  exit 1
fi

bundle_dir="${bundle_root}/${plugin_name}.vst3"
binary_dir="${bundle_dir}/Contents/x86_64-win"
binary_path="${binary_dir}/${plugin_name}.vst3"
header_source="$(dirname "$plugin_dll")/${plugin_name}.h"
header_path="${binary_dir}/${plugin_name}.h"

mkdir -p "$binary_dir"
cp "$plugin_dll" "$binary_path"

if [[ -f "$header_source" ]]; then
  cp "$header_source" "$header_path"
fi

loader_source="$(dirname "$plugin_dll")/WebView2Loader.dll"
if [[ -f "$loader_source" ]]; then
  cp "$loader_source" "$binary_dir/WebView2Loader.dll"
fi

echo "packaged Windows VST3 bundle at: $bundle_dir"
