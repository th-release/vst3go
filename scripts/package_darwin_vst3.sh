#!/usr/bin/env bash
set -euo pipefail

if [[ $# -lt 1 ]]; then
  echo "usage: $0 <plugin-dylib> [bundle-root] [plugin-name]" >&2
  exit 1
fi

plugin_dylib="$1"
bundle_root="${2:-dist/macos}"
plugin_name="${3:-vst3go}"
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "${script_dir}/.." && pwd)"

if [[ "$plugin_dylib" != /* ]]; then
  plugin_dylib="${repo_root}/${plugin_dylib}"
fi
if [[ "$bundle_root" != /* ]]; then
  bundle_root="${repo_root}/${bundle_root}"
fi

if [[ ! -f "$plugin_dylib" ]]; then
  echo "plugin dylib not found: $plugin_dylib" >&2
  exit 1
fi

bundle_dir="${bundle_root}/${plugin_name}.vst3"
contents_dir="${bundle_dir}/Contents"
macos_dir="${contents_dir}/MacOS"
binary_path="${macos_dir}/${plugin_name}"
header_source="$(dirname "$plugin_dylib")/${plugin_name}.h"
header_path="${macos_dir}/${plugin_name}.h"
plist_path="${contents_dir}/Info.plist"

mkdir -p "$macos_dir"
cp "$plugin_dylib" "$binary_path"

if [[ ! -f "$header_source" ]]; then
  echo "missing macOS VST3 header sidecar: $header_source" >&2
  exit 1
fi

cp "$header_source" "$header_path"

cat >"$plist_path" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>CFBundleDevelopmentRegion</key>
  <string>en</string>
  <key>CFBundleExecutable</key>
  <string>${plugin_name}</string>
  <key>CFBundleIdentifier</key>
  <string>com.example.${plugin_name}</string>
  <key>CFBundleInfoDictionaryVersion</key>
  <string>6.0</string>
  <key>CFBundleName</key>
  <string>${plugin_name}</string>
  <key>CFBundlePackageType</key>
  <string>BNDL</string>
  <key>CFBundleShortVersionString</key>
  <string>0.1.0</string>
  <key>CFBundleVersion</key>
  <string>1</string>
  <key>CSResourcesFileMapped</key>
  <true/>
</dict>
</plist>
EOF

echo "packaged macOS VST3 bundle at: $bundle_dir"
