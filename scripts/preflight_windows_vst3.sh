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

echo "Windows build toolchain preflight passed."
