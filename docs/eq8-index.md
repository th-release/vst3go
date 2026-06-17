# EQ8 Docs Index

Use this page as the entry point for the EQ8 example plugin workflow.

## Start Here

- [`docs/eq8-operator-guide.md`](eq8-operator-guide.md): one-page build, install, validate loop
- [`docs/eq8-build-quickstart.md`](eq8-build-quickstart.md): shorter build and package commands
- [`docs/eq8-host-validation.md`](eq8-host-validation.md): host-side validation checklist
- [`docs/cross-platform-development.md`](cross-platform-development.md): shared runtime and OS split

## Recommended Reading Order

1. read `docs/eq8-operator-guide.md`
2. use `docs/eq8-build-quickstart.md` when you only need package commands
3. use `docs/eq8-host-validation.md` when you are inside a real host
4. refer back to `docs/cross-platform-development.md` for host-agnostic workflow boundaries

## Code And Build Entry Points

- `example/eq8/README.md`: example-specific source layout and commands
- `example/eq8/web/`: Vite + TypeScript + React editor source
- `cmd/eq8-dylib/main_darwin.go`: macOS EQ8 bundle entrypoint
- `cmd/eq8-dll/main_windows.go`: Windows EQ8 bundle entrypoint
- `scripts/build_darwin_vst3.sh`: macOS build helper
- `scripts/build_windows_vst3.sh`: Windows build helper

## Validation Targets

- `go test -timeout=30000s ./...`
- `go build ./...`
- `cd example/eq8/web && npm run build`
- `just eq8-mac-release`
- `just eq8-win-release`

If you are planning a bigger cross-platform change, start with `docs/development-workflow.md` and then return here for the EQ8-specific path.
