#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "${script_dir}/.." && pwd)"
tmpdir="$(mktemp -d "${TMPDIR:-/tmp}/vst3go-win-smoke.XXXXXX")"
trap 'rm -rf "$tmpdir"' EXIT

make_fake_wrappers() {
  local target_dir="$1"
  mkdir -p "$target_dir"

  cat >"$target_dir/bash" <<'EOF'
#!/bin/bash
exec /bin/bash "$@"
EOF
  chmod +x "$target_dir/bash"

  cat >"$target_dir/go" <<'EOF'
#!/bin/bash
set -euo pipefail

if [[ "${1:-}" == "env" && "${2:-}" == "GOHOSTOS" ]]; then
  printf 'linux\n'
  exit 0
fi

echo "unexpected go invocation: $*" >&2
exit 1
EOF
  chmod +x "$target_dir/go"

  cat >"$target_dir/mktemp" <<'EOF'
#!/bin/bash
exec /usr/bin/mktemp "$@"
EOF
  chmod +x "$target_dir/mktemp"

  cat >"$target_dir/rm" <<'EOF'
#!/bin/bash
exec /bin/rm "$@"
EOF
  chmod +x "$target_dir/rm"

  cat >"$target_dir/dirname" <<'EOF'
#!/bin/bash
exec /usr/bin/dirname "$@"
EOF
  chmod +x "$target_dir/dirname"
}

make_fake_compiler() {
  local target_dir="$1"
  mkdir -p "$target_dir"

  cat >"$target_dir/x86_64-w64-mingw32-gcc" <<'EOF'
#!/bin/bash
set -euo pipefail

args=("$@")
for arg in "${args[@]}"; do
  if [[ "$arg" == "-E" ]]; then
    while IFS= read -r _; do
      :
    done
    exit 0
  fi
done

output_file=""
for ((index = 0; index < ${#args[@]}; index++)); do
  if [[ "${args[index]}" == "-o" && $((index + 1)) -lt ${#args[@]} ]]; then
    output_file="${args[index + 1]}"
    break
  fi
done

if [[ -n "$output_file" ]]; then
  : > "$output_file"
fi
exit 0
EOF
  chmod +x "$target_dir/x86_64-w64-mingw32-gcc"
}

success_bin="$tmpdir/success-bin"
make_fake_wrappers "$success_bin"
make_fake_compiler "$success_bin"

selected_cc="$(PATH="$success_bin" bash "$repo_root/scripts/select_windows_cc.sh")"
expected_cc="$success_bin/x86_64-w64-mingw32-gcc"
if [[ "$selected_cc" != "$expected_cc" ]]; then
  echo "unexpected compiler selection: $selected_cc" >&2
  echo "expected: $expected_cc" >&2
  exit 1
fi

PATH="$success_bin" bash "$repo_root/scripts/preflight_windows_vst3.sh" >/dev/null

failure_bin="$tmpdir/failure-bin"
make_fake_wrappers "$failure_bin"

if output="$(PATH="$failure_bin" bash "$repo_root/scripts/select_windows_cc.sh" 2>&1)"; then
  echo "expected selector failure, but it succeeded" >&2
  exit 1
fi

if [[ "$output" != *"no Windows-capable C compiler found"* ]]; then
  echo "selector failure output did not explain the problem" >&2
  printf '%s\n' "$output" >&2
  exit 1
fi

if [[ "$output" != *"x86_64-w64-mingw32-gcc: not found"* ]]; then
  echo "selector failure output did not list tried candidates" >&2
  printf '%s\n' "$output" >&2
  exit 1
fi

echo "Windows toolchain smoke test passed."
