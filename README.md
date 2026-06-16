# VST3Go

`vst3go` is the low-level Go interface and runtime layer for VST3.

This fork is intentionally different from the original `justyntemme/vst3go` project. It is being kept deliberately narrow as a clean VST3 binding/runtime layer, while the higher-level DSP and synth-oriented functionality lives in the companion `synthkit` repository.

It contains:

- VST3 C API headers and C bridge glue
- Go bindings for VST3 types and process structures
- The plugin wrapper/runtime used to expose Go processors as VST3 plugins
- Generic plugin authoring primitives such as buses, parameters, process context, state, and MIDI/event transport

The higher-level audio and synth layer now lives in the companion `synthkit` repository.

For the web-rendered plugin editor flow, see [`docs/web-editor-bridge.md`](docs/web-editor-bridge.md).
For Windows build and packaging notes, see [`docs/windows-build.md`](docs/windows-build.md).

## Scope

`vst3go` keeps the VST3-facing pieces:

- `bridge/`
- `include/vst3/`
- `pkg/vst3/`
- `pkg/plugin/`
- `pkg/midi/`
- `pkg/framework/bus`
- `pkg/framework/param`
- `pkg/framework/plugin`
- `pkg/framework/process`
- `pkg/framework/state`

The companion `synthkit` repo owns:

- `pkg/dsp/`
- `pkg/framework/dsp`
- `pkg/framework/debug`
- `pkg/framework/voice`
- `examples/`

## Companion Repo

Recommended name: `synthkit`

Purpose:

- DSP building blocks
- instrument/effect showcase plugins
- higher-level audio helpers built on top of `vst3go`

## Development

Core checks in this repo:

```bash
just fmt
just fmt-check
just test
```

## Validation

- This repo validates the shared runtime layer with `just test` and the Windows build path with `just windows-smoke` when the host shell is available.
- Platform-specific integration checks, host app behavior, and downstream DSP/plugin wiring belong in the companion `synthkit` repository or the consumer project.
- Windows packaging and bundle layout checks in this repo stay focused on the generated DLL, header sidecar, and layout contract.

Windows-specific editor and packaging notes:

- `bridge/windows_dll.c` provides the DLL entry point.
- `bridge/bridge.c` exports `GetPluginFactory` for the Windows build.
- `pkg/plugin/cbridge/windows_dll_windows.go` makes the DLL entry source part of the Windows cgo build.
- `pkg/plugin/editor_view_windows.c` expects WebView2 headers and loader support.
- Windows support currently targets `amd64` and the `x86_64-win` bundle layout.
- `just windows-preflight` auto-selects a compatible Windows compiler when needed and checks for `windows.h`, `WebView2.h`, and `WebView2Loader`. If no compiler works, it reports the candidates it tried and why each one failed.
- `just windows-smoke` runs a local script-only smoke test that exercises selection, preflight, packaging, and bundle validation without a real Windows host.
- `just windows-build-dll` builds the Windows shared library when the Windows toolchain is available.
- `just windows-build <dll>` assembles a VST3 bundle directory from a built Windows DLL and requires the generated header sidecar.
- `just windows-check-bundle` validates the resulting bundle layout.
- `just windows-release` runs preflight once, then build and validation in one go.
- `docs/windows-build.md` describes the expected bundle shape.

## License

This project is licensed under the MIT License. The bundled VST3 SDK headers remain under their respective licenses.
