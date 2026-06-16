# Windows Build Notes

This repo now has a Windows editor-view scaffold and a `GetPluginFactory` export path, but Windows packaging still needs the host toolchain and bundle layout to be wired by the consumer build.

## Required Pieces

- A Windows CGO toolchain
- The Microsoft WebView2 SDK headers and loader library
- The WebView2 runtime installed on the target machine

## Current Contract

- `bridge/windows_dll.c` provides a minimal `DllMain`.
- `bridge/bridge.c` exports `GetPluginFactory` with a Windows-friendly symbol annotation.
- `pkg/plugin/editor_view_windows.c` hosts the editor through WebView2.
- `pkg/plugin/editor_assets/` provides the HTML/CSS/JS snapshot shell.

## Expected Bundle Shape

- The VST3 consumer should place the built library in the normal plugin bundle layout for the host.
- The Windows host-side package should ship the DLL together with the WebView2 loader/runtime dependencies it expects.
- The editor bridge should stay inside the plugin DLL; the web assets remain embedded in the binary.
- The repo ships `scripts/package_windows_vst3.sh` and `just windows-package` to assemble the bundle from an already-built Windows DLL.

## Build Expectations

- Cross-compiling from Linux is possible only when the Windows C toolchain and WebView2 headers are available.
- The repo-level Go tests continue to validate the shared runtime code on the current platform.
- Windows-specific packaging validation should happen on an actual Windows toolchain.
