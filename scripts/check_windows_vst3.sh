#!/usr/bin/env bash
set -euo pipefail

if [[ $# -lt 1 ]]; then
  echo "usage: $0 <bundle-root> [plugin-name]" >&2
  exit 1
fi

bundle_root="$1"
plugin_name="${2:-vst3go}"
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "${script_dir}/.." && pwd)"

if [[ "$bundle_root" != /* ]]; then
  bundle_root="${repo_root}/${bundle_root}"
fi

layout_bundle="${bundle_root}/${plugin_name}.vst3"
layout_binary="${layout_bundle}/Contents/x86_64-win/${plugin_name}.vst3"
layout_header="${layout_bundle}/Contents/x86_64-win/${plugin_name}.h"

if [[ ! -f "$layout_binary" ]]; then
  echo "missing Windows VST3 binary: $layout_binary" >&2
  exit 1
fi

if [[ ! -f "$layout_header" ]]; then
  echo "missing Windows VST3 header sidecar: $layout_header" >&2
  exit 1
fi

echo "Windows VST3 bundle looks valid: $layout_bundle"
