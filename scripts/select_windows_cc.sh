#!/usr/bin/env bash
set -euo pipefail

host_os="$(go env GOHOSTOS)"
if [[ "$host_os" == "windows" ]]; then
  printf '%s\n' "${CC:-cc}"
  exit 0
fi

tmpdir="$(mktemp -d "${TMPDIR:-/tmp}/vst3go-win-cc.XXXXXX")"
trap 'rm -rf "$tmpdir"' EXIT

printf 'int main(void) { return 0; }\n' > "$tmpdir/check.c"

candidate_list=()
if [[ -n "${CC:-}" ]]; then
  candidate_list+=("$CC")
fi
candidate_list+=("x86_64-w64-mingw32-gcc" "clang" "cc")

for candidate in "${candidate_list[@]}"; do
  candidate_path="$(command -v "$candidate" 2>/dev/null || true)"
  if [[ -z "$candidate_path" ]]; then
    continue
  fi

  if ! printf '#include <windows.h>\n' | "$candidate_path" -E -x c - >/dev/null 2>&1; then
    continue
  fi

  if ! printf '#include <WebView2.h>\n' | "$candidate_path" -E -x c - >/dev/null 2>&1; then
    continue
  fi

  if "$candidate_path" "$tmpdir/check.c" -lWebView2Loader -o "$tmpdir/check" >/dev/null 2>&1; then
    printf '%s\n' "$candidate_path"
    exit 0
  fi
done

echo "no Windows-capable C compiler found" >&2
exit 1
