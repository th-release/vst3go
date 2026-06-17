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
failed_step=""
completed_steps=()

append_report_footer() {
  {
    printf '\n## Execution Logs\n\n'
    printf -- '- Log directory: %s\n' "$run_log_dir"
    if [[ -n "$failed_step" ]]; then
      printf -- '- Failed step: %s\n' "$failed_step"
    else
      printf -- '- Failed step: none\n'
    fi

    printf '\n## Log Files\n\n'
    find "$run_log_dir" -maxdepth 1 -type f -name '*.log' | sort | while read -r log_file; do
      printf -- '- %s\n' "$log_file"
    done

    printf '\n## Runner Status\n\n'
    if [[ -n "$failed_step" ]]; then
      printf -- '- Result: blocked\n'
    else
      printf -- '- Result: pass\n'
    fi
    if [[ ${#completed_steps[@]} -gt 0 ]]; then
      printf -- '- Completed steps: %s\n' "${completed_steps[*]}"
    fi
  } >>"$report_path"
}

trap append_report_footer EXIT

run_step() {
  local label="$1"
  shift
  local log_file="$run_log_dir/${label}.log"
  local status=0
  local output=""
  printf '## %s\n\n' "$label" | tee "$log_file"
  printf '$ %s\n\n' "$*" | tee -a "$log_file"
  set +e
  output="$("$@" 2>&1)"
  status=$?
  set -e
  if [[ -n "$output" ]]; then
    printf '%s\n' "$output" | tee -a "$log_file"
  fi
  if [[ $status -eq 0 ]]; then
    printf 'result: success\n' | tee -a "$log_file"
  else
    printf 'result: failure (%d)\n' "$status" | tee -a "$log_file"
  fi
  if [[ $status -eq 0 ]]; then
    completed_steps+=("$label")
  else
    failed_step="$label"
  fi
  return "$status"
}

if ! run_step go-test env GOCACHE="${repo_root}/.gocache" go test -timeout=30000s ./...; then
  exit 1
fi

if ! run_step go-build env GOCACHE="${repo_root}/.gocache" go build ./...; then
  exit 1
fi

if ! run_step windows-preflight bash "${repo_root}/scripts/preflight_windows_vst3.sh"; then
  exit 1
fi

if ! run_step windows-build-dll bash "${repo_root}/scripts/build_windows_vst3.sh" "${bundle_root}" "${plugin_name}"; then
  exit 1
fi

if ! run_step windows-check-bundle bash "${repo_root}/scripts/check_windows_vst3.sh" "${bundle_root}" "${plugin_name}"; then
  exit 1
fi

if [[ "${VST3GO_SKIP_RACE:-0}" != "1" ]] && ! run_step go-test-race env GOCACHE="${repo_root}/.gocache" go test -race -p 1 -timeout=30000s ./...; then
  exit 1
fi

if [[ -z "$failed_step" ]]; then
  printf 'Windows validation run prepared: %s\n' "$report_path"
else
  printf 'Windows validation run blocked at %s: %s\n' "$failed_step" "$report_path" >&2
  exit 1
fi
