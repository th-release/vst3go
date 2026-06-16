set shell := ["zsh", "-cu"]

default:
  @just --list

fmt:
  env XDG_CACHE_HOME=/tmp/treefmt-cache treefmt --allow-missing-formatter

fmt-check:
  env XDG_CACHE_HOME=/tmp/treefmt-cache treefmt --allow-missing-formatter --fail-on-change

check-formatted: fmt-check

lint:
  @if ! command -v golangci-lint >/dev/null 2>&1; then echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; exit 1; fi
  env GOCACHE=/tmp/gocache GOMODCACHE=/tmp/gomodcache GOLANGCI_LINT_CACHE=/tmp/golangci-lint-cache golangci-lint run ./pkg/...

lint-fix:
  @if ! command -v golangci-lint >/dev/null 2>&1; then echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; exit 1; fi
  env GOCACHE=/tmp/gocache GOMODCACHE=/tmp/gomodcache GOLANGCI_LINT_CACHE=/tmp/golangci-lint-cache golangci-lint run --fix ./pkg/...

test-go:
  env GOCACHE=/tmp/gocache GOMODCACHE=/tmp/gomodcache go test ./pkg/...

test: fmt-check test-go

windows-package plugin_dll bundle_root='dist/windows' plugin_name='vst3go':
  bash scripts/package_windows_vst3.sh "{{plugin_dll}}" "{{bundle_root}}" "{{plugin_name}}"

windows-build plugin_dll bundle_root='dist/windows' plugin_name='vst3go':
  just windows-package "{{plugin_dll}}" "{{bundle_root}}" "{{plugin_name}}"

windows-check-bundle bundle_root='dist/windows' plugin_name='vst3go':
  bash scripts/check_windows_vst3.sh "{{bundle_root}}" "{{plugin_name}}"

windows-build-dll bundle_root='dist/windows' plugin_name='vst3go':
  bash scripts/build_windows_vst3.sh "{{bundle_root}}" "{{plugin_name}}"

windows-preflight:
  bash scripts/preflight_windows_vst3.sh

windows-release bundle_root='dist/windows' plugin_name='vst3go':
  just windows-preflight
  just windows-build-dll "{{bundle_root}}" "{{plugin_name}}"
  just windows-check-bundle "{{bundle_root}}" "{{plugin_name}}"

fix:
    just lint-fix
    just fmt
