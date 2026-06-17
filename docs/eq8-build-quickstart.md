# EQ8 Build Quickstart

Use this when you want the EQ8 example plugin itself, not just the shared runtime.

## Shared Work

These commands work on any host OS:

```bash
go test -timeout=30000s ./...
go build ./...
cd example/eq8/web && npm run build
```

The first command validates the shared runtime and example plugin logic.  
The second confirms the repo still compiles.  
The third regenerates the browser editor shell from the Vite + React source.

## macOS Bundle

On macOS, build and package the EQ8 bundle with:

```bash
just eq8-mac-release
```

Expected output:

- a `.vst3` bundle under `dist/macos/eq8.vst3`
- the binary in `Contents/MacOS/eq8`
- the generated header sidecar in `Contents/MacOS/eq8.h`
- `Info.plist` in `Contents/Info.plist`

Bundle shape:

```text
dist/macos/eq8.vst3/
└─ Contents/
   ├─ Info.plist
   └─ MacOS/
      ├─ eq8
      └─ eq8.h
```

Minimal `Info.plist` fields:

```xml
<key>CFBundleExecutable</key>
<string>eq8</string>
<key>CFBundleIdentifier</key>
<string>com.example.eq8</string>
<key>CFBundlePackageType</key>
<string>BNDL</string>
<key>CFBundleVersion</key>
<string>1</string>
```

## Windows Bundle

On Windows, build and package the EQ8 bundle with:

```bash
just eq8-win-release
```

Expected output:

- a `.vst3` bundle under `dist/windows/eq8.vst3`
- the binary in `Contents/x86_64-win/eq8.vst3`
- the generated header sidecar in `Contents/x86_64-win/eq8.h`
- `WebView2Loader.dll` if it is available beside the DLL output

Bundle shape:

```text
dist/windows/eq8.vst3/
└─ Contents/
   └─ x86_64-win/
      ├─ eq8.vst3
      ├─ eq8.h
      └─ WebView2Loader.dll
```

## Output Summary

| Host | Build Command | Bundle Root | Binary Path |
| --- | --- | --- | --- |
| macOS | `just eq8-mac-release` | `dist/macos/eq8.vst3` | `Contents/MacOS/eq8` |
| Windows | `just eq8-win-release` | `dist/windows/eq8.vst3` | `Contents/x86_64-win/eq8.vst3` |

## Install Locations

After packaging, copy the `.vst3` bundle into the standard VST3 folder for the host OS:

- **macOS**: `~/Library/Audio/Plug-Ins/VST3/`
- **Windows**: `%APPDATA%\\VST3\\` or the host's configured VST3 directory

The host should see the `.vst3` bundle folder itself, not the raw binary inside it.

## What To Check Next

After packaging, open the bundle in a real VST3 host and verify:

- the editor opens
- parameter changes update the processor
- state restore behaves as expected
- reopen and rescan flows remain stable

If you are only changing Go DSP or editor logic, stay in the shared work loop until that stabilizes.

For a one-page operator workflow, see `docs/eq8-operator-guide.md`.

## Quick Checklist

- [ ] `go test -timeout=30000s ./...`
- [ ] `go build ./...`
- [ ] `cd example/eq8/web && npm run build`
- [ ] `just eq8-mac-release` or `just eq8-win-release`
- [ ] copy the `.vst3` bundle into the host VST3 directory
- [ ] open the plugin in a real VST3 host
- [ ] verify parameter edits and state restore
