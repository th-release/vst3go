# EQ8 Operator Guide

Use this when you want one concise workflow for building, installing, and validating the EQ8 example plugin.

## 1. Build The Shared Pieces

Run these anywhere:

```bash
go test -timeout=30000s ./...
go build ./...
cd example/eq8/web && npm run build
```

This validates the runtime, compiles the repo, and regenerates the browser editor shell.

## 2. Package For Your OS

### macOS

```bash
just eq8-mac-release
```

Result:

- `dist/macos/eq8.vst3`
- `Contents/MacOS/eq8`
- `Contents/MacOS/eq8.h`
- `Contents/Info.plist`

### Windows

```bash
just eq8-win-release
```

Result:

- `dist/windows/eq8.vst3`
- `Contents/x86_64-win/eq8.vst3`
- `Contents/x86_64-win/eq8.h`
- `WebView2Loader.dll` when present

## 3. Install It

- macOS: copy `dist/macos/eq8.vst3` to `~/Library/Audio/Plug-Ins/VST3/`
- Windows: copy `dist/windows/eq8.vst3` to `%APPDATA%\\VST3\\` or the host VST3 folder

The host should scan the `.vst3` bundle folder itself.

## 4. Scan In The Host

1. close the host if it is already open
2. copy the `.vst3` bundle into the VST3 folder
3. relaunch the host or trigger a plugin rescan
4. confirm the plugin appears as `EQ8 Example`
5. insert it on a track or test slot

If the plugin does not appear:

- recheck the bundle folder name
- recheck the internal binary path
- make sure the host is pointed at the right VST3 directory
- rebuild the bundle after any code change that affects the entrypoint or package name

## 5. Validate In A Host

Use the same checks on both platforms:

- plugin discovery succeeds
- editor opens
- band controls and global controls work
- parameter automation reflects back into the UI
- save and restore work
- repeated close/reopen cycles stay stable

For platform-specific details, use:

- `docs/eq8-build-quickstart.md`
- `docs/eq8-host-validation.md`

## 6. Keep A Tight Loop

When editing DSP or the editor source:

1. run `go test -timeout=30000s ./...`
2. run `go build ./...`
3. run `cd example/eq8/web && npm run build`
4. package only when behavior is stable
5. validate in a real host
