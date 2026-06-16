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

## License

This project is licensed under the MIT License. The bundled VST3 SDK headers remain under their respective licenses.
