# EQ8 Example

This directory contains a complete downstream plugin example built on `vst3go`.

What it includes:

- a real `Plugin` entrypoint
- a stereo `Processor`
- eight EQ bands with per-band enable, type, frequency, gain, and Q controls
- bypass and analyzer controls
- a browser-rendered editor shell under `web/editor/`
- a real Vite + TypeScript + React source tree under `web/src/`
- plain `go test` coverage for the public shape of the plugin

Parameter layout:

- `1` `Input Gain`
- `2` `Output Gain`
- `3` `Bypass`
- `4` `Analyzer`
- `100..104` `Band 1`
- `110..114` `Band 2`
- `120..124` `Band 3`
- `130..134` `Band 4`
- `140..144` `Band 5`
- `150..154` `Band 6`
- `160..164` `Band 7`
- `170..174` `Band 8`

The processor is intentionally simple:

- it uses stereo buses
- it stores no custom state beyond parameters
- it bypasses cleanly when the host sets the bypass control
- it rebuilds band coefficients from the live registry values on each block

The web editor is intentionally simple too:

- it reads the same `EditorSnapshot` shape the runtime uses
- it renders the eight band cards, global controls, and response graph
- it posts parameter changes back in the same normalized/plain message shape used by the host bridge

The `web/` directory is the Vite project root:

- `npm install` installs the editor dependencies
- `npm run build` writes the generated shell back into `web/editor/`
- `example/eq8/web.go` embeds the generated shell for the Go example runtime

Bundle commands for this example:

- macOS: `just eq8-mac-release`
- Windows: `just eq8-win-release`

For a shorter step-by-step version, see [`docs/eq8-build-quickstart.md`](../../docs/eq8-build-quickstart.md).

You can import this package from a downstream repo and wire it into the host/runtime layer without changing the example plugin shape.
