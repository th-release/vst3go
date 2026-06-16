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

if [[ "${1:-}" == "build" ]]; then
  args=("$@")
  output_file=""
  for ((index = 0; index < ${#args[@]}; index++)); do
    if [[ "${args[index]}" == "-o" && $((index + 1)) -lt ${#args[@]} ]]; then
      output_file="${args[index + 1]}"
      break
    fi
  done

  if [[ -z "$output_file" ]]; then
    echo "missing build output path" >&2
    exit 1
  fi

  mkdir -p "$(dirname "$output_file")"
  : > "$output_file"
  : > "${output_file}.h"
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

  cat >"$target_dir/mkdir" <<'EOF'
#!/bin/bash
exec /bin/mkdir "$@"
EOF
  chmod +x "$target_dir/mkdir"

  cat >"$target_dir/cp" <<'EOF'
#!/bin/bash
exec /bin/cp "$@"
EOF
  chmod +x "$target_dir/cp"
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

package_output_dir="$tmpdir/package-output"
PATH="$success_bin" bash "$repo_root/scripts/build_windows_vst3.sh" "$package_output_dir" >/dev/null
bash "$repo_root/scripts/check_windows_vst3.sh" "$package_output_dir" >/dev/null

if [[ ! -f "$package_output_dir/vst3go.vst3/Contents/x86_64-win/vst3go.vst3" ]]; then
  echo "builder did not create the Windows DLL output" >&2
  exit 1
fi

if [[ ! -f "$package_output_dir/vst3go.vst3/Contents/x86_64-win/vst3go.vst3.h" ]]; then
  echo "builder did not create the Windows header sidecar" >&2
  exit 1
fi

package_input_dir="$tmpdir/package-input"
package_output_copy="$tmpdir/package-output-copy"
mkdir -p "$package_input_dir" "$package_output_copy"
printf 'fake dll\n' > "$package_input_dir/vst3go.vst3"
printf 'fake header\n' > "$package_input_dir/vst3go.h"
printf 'fake loader\n' > "$package_input_dir/WebView2Loader.dll"

bash "$repo_root/scripts/package_windows_vst3.sh" "$package_input_dir/vst3go.vst3" "$package_output_copy" >/dev/null
bash "$repo_root/scripts/check_windows_vst3.sh" "$package_output_copy" >/dev/null

if [[ ! -f "$package_output_copy/vst3go.vst3/Contents/x86_64-win/WebView2Loader.dll" ]]; then
  echo "packager did not copy the WebView2 loader when present" >&2
  exit 1
fi

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
