# Cross-Platform Development

This repository is designed so you can develop the shared plugin/runtime code on any host OS, while keeping final host-specific packaging isolated to the platform that owns it.

## The Core Idea

Keep these layers separate:

- **Shared Go runtime**: `pkg/plugin`, `pkg/framework/*`, `pkg/vst3`, `pkg/midi`
- **Downstream plugin implementation**: your plugin package, processor, and metadata
- **Web editor source**: `example/eq8/web/` or your own `web/` tree
- **Platform shell**: macOS `WKWebView`, Windows `WebView2`, and bundle/entrypoint glue

That lets you work on behavior, UI, and state contracts without needing to rewrite the plugin for each OS.

## What You Can Build Anywhere

These steps are safe on macOS, Windows, and Linux:

- `go test -timeout=30000s ./...`
- `go build ./...`
- `cd example/eq8/web && npm run build`
- editing the Go snapshot/state contracts
- editing the React editor source and rebuilding the browser shell

This is the best loop for most day-to-day plugin work.

## What Still Needs A Native Platform Toolchain

Final VST3 packaging is still platform-specific:

- **macOS**: requires the macOS SDK/toolchain and the native editor host glue
- **Windows**: requires MinGW-w64 or another Windows-capable CGO toolchain plus WebView2 headers/loader
- **Linux**: good for shared Go development and web editor work, but not a final VST3 packaging target in this repository yet

So the rule is:

- develop the shared code anywhere
- package and host-test on the target OS

## Recommended Development Loop

1. change the Go runtime or editor model
2. run `go test -timeout=30000s ./...`
3. run `go build ./...`
4. rebuild the web editor when UI changes land
5. package on the target OS
6. verify in a real VST3 host

For the EQ8 example specifically:

- `go test -timeout=30000s ./...`
- `cd example/eq8/web && npm run build`
- inspect the generated editor shell in `example/eq8/web/editor/`
- then validate the bundle in the target host

## Platform Roles

### macOS

- Hosts the editor with `WKWebView`
- Uses the same snapshot contract as the shared Go runtime
- Should keep platform code thin and avoid duplicating editor state
- Packages as a `.vst3` bundle with `Contents/MacOS/<plugin-name>` and `Contents/Info.plist`

### Windows

- Hosts the editor with `WebView2`
- Uses the same snapshot and message contract
- Packages as a `.vst3` bundle with the Windows layout described in `docs/windows-build.md`

### Linux

- Best used for Go-level runtime work and web editor iteration
- Can validate shared state, parameter flow, and the browser editor build
- Should not be treated as the final host packaging reference unless the repository adds explicit Linux packaging support later

## For Downstream Plugin Projects

A clean downstream layout usually looks like this:

- `plugin/` for metadata and processor implementation
- `web/` for the React editor source
- `cmd/<plugin>-dll` or a platform-specific entrypoint package
- `dist/` for packaged bundles and validation artifacts

If you follow that split, each OS only needs a thin packaging layer while the audio/editor logic stays shared.

For a concrete example, `example/eq8` now has:

- `cmd/vst3go-dylib` for the shared runtime bundle
- `cmd/eq8-dylib` for the EQ8-specific macOS bundle
- the same Go runtime, editor snapshot contract, and React web source
