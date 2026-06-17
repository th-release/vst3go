#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "${script_dir}/.." && pwd)"

report_path="${1:-}"
bundle_root="${2:-dist/windows}"
plugin_name="${3:-vst3go}"

if [[ -z "$report_path" ]]; then
  report_path="${repo_root}/docs/windows-validation-$(date +%F).md"
fi

report_init_output="$(bash "${script_dir}/init_windows_validation_report.sh" "$report_path")"
printf '%s\n' "$report_init_output"

run_log_dir="$(dirname "$report_path")/windows-validation-logs-$(date +%s)"
mkdir -p "$run_log_dir"

run_step() {
  local label="$1"
  shift
  local log_file="$run_log_dir/${label}.log"
  {
    printf '## %s\n\n' "$label"
    printf '$ %s\n\n' "$*"
    if "$@"; then
      printf 'result: success\n'
    else
      local exit_code=$?
      printf 'result: failure (%d)\n' "$exit_code"
      return "$exit_code"
    fi
  } | tee "$log_file"
}

{
  printf '\n## Execution Logs\n\n'
  printf -- '- Log directory: %s\n' "$run_log_dir"
} >>"$report_path"

run_step go-test env GOCACHE="${repo_root}/.gocache" go test -timeout=30000s ./...
run_step go-build env GOCACHE="${repo_root}/.gocache" go build ./...
run_step windows-preflight bash "${repo_root}/scripts/preflight_windows_vst3.sh"
run_step windows-build-dll just windows-build-dll "${bundle_root}" "${plugin_name}"
run_step windows-check-bundle just windows-check-bundle "${bundle_root}" "${plugin_name}"

if [[ "${VST3GO_SKIP_RACE:-0}" != "1" ]]; then
  run_step go-test-race env GOCACHE="${repo_root}/.gocache" go test -race -p 1 -timeout=30000s ./...
fi

{
  printf '\n## Log Files\n\n'
  find "$run_log_dir" -maxdepth 1 -type f -name '*.log' | sort | while read -r log_file; do
    printf -- '- %s\n' "$log_file"
  done
} >>"$report_path"

printf 'Windows validation run prepared: %s\n' "$report_path"
