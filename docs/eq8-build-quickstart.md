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

## What To Check Next

After packaging, open the bundle in a real VST3 host and verify:

- the editor opens
- parameter changes update the processor
- state restore behaves as expected
- reopen and rescan flows remain stable

If you are only changing Go DSP or editor logic, stay in the shared work loop until that stabilizes.

