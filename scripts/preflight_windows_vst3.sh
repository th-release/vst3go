#!/usr/bin/env bash
set -euo pipefail

host_os="$(go env GOHOSTOS)"
cc_bin="${CC:-cc}"

if [[ "$host_os" == "windows" ]]; then
  exit 0
fi

if ! printf '#include <windows.h>\n' | "$cc_bin" -E -x c - >/dev/null 2>&1; then
  echo "Windows build toolchain is not ready: need a compiler that can find <windows.h> (for example MinGW-w64)." >&2
  echo "Set CC to a Windows-capable compiler before running the Windows build again." >&2
  exit 1
fi

if ! printf '#include <WebView2.h>\n' | "$cc_bin" -E -x c - >/dev/null 2>&1; then
  echo "Windows build toolchain is missing the WebView2 SDK headers." >&2
  echo "Make sure the compiler can include <WebView2.h> before running the Windows build again." >&2
  exit 1
fi

tmpdir="$(mktemp -d "${TMPDIR:-/tmp}/vst3go-win-preflight.XXXXXX")"
trap 'rm -rf "$tmpdir"' EXIT

printf 'int main(void) { return 0; }\n' > "$tmpdir/check.c"
if ! "$cc_bin" "$tmpdir/check.c" -lWebView2Loader -o "$tmpdir/check" >/dev/null 2>&1; then
  echo "Windows build toolchain cannot link against WebView2Loader." >&2
  echo "Make sure the loader library is installed and discoverable before running the Windows build again." >&2
  exit 1
fi

echo "Windows build toolchain preflight passed."
