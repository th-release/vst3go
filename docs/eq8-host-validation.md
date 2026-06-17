# EQ8 Host Validation

Use this after you have built and installed the EQ8 bundle into a real VST3 host.

## Before You Start

- build the shared runtime and editor shell
- package the bundle for your OS
- copy the `.vst3` bundle into the host's VST3 folder
- scan or rescan the plugin in the host

For the build and packaging commands, see `docs/eq8-operator-guide.md`.

## Validation Steps

### 1. Plugin Discovery

- the host sees the EQ8 plugin
- the plugin name and vendor text look correct
- the host does not crash while scanning the plugin

### 2. Editor Open

- the editor window opens
- the browser-rendered UI appears
- the band controls and graph render without blank regions

### 3. Audio and Parameters

- parameter changes in the editor reach the processor
- host automation updates are reflected in the UI
- bypass behaves as expected
- input and output gain controls respond correctly

### 4. State Restore

- save the project or preset
- close and reopen the host
- the EQ8 state returns correctly

### 5. Reopen and Rescan

- close the editor and reopen it
- scan or rescan the plugin again
- the host remains stable across repeated open/close cycles

## macOS Host Checklist

Use this when testing the `.vst3` bundle built by `just eq8-mac-release`:

- copy `dist/macos/eq8.vst3` into `~/Library/Audio/Plug-Ins/VST3/`
- rescan the plugin inside the host
- confirm the host sees `EQ8 Example` from `example/eq8/plugin.go`
- confirm the editor opens inside `WKWebView`
- confirm band controls, bypass, and analyzer all update correctly
- reopen the editor after closing it once
- save and reopen the project to confirm state restore

If the host exposes a plugin manager or rescan list, use it to confirm the bundle is discovered by name, not just by file presence.

## Windows Host Checklist

Use this when testing the `.vst3` bundle built by `just eq8-win-release`:

- copy `dist/windows/eq8.vst3` into `%APPDATA%\\VST3\\` or the host's configured VST3 folder
- rescan the plugin inside the host
- confirm the host sees `EQ8 Example` from `example/eq8/plugin.go`
- confirm the editor opens inside `WebView2`
- confirm `WebView2Loader.dll` is present in the bundle when packaging expects it
- confirm parameter updates reach the processor and redraw in the UI
- reopen the editor after closing it once
- save and reopen the project to confirm state restore

If a host keeps a plugin cache, make sure you invalidate or rescan it after each new bundle build.

## macOS Notes

- install the bundle under `~/Library/Audio/Plug-Ins/VST3/`
- verify `WKWebView` opens the editor and keeps focus reasonably

## Windows Notes

- install the bundle under `%APPDATA%\\VST3\\` or the host's configured plugin directory
- verify `WebView2` opens the editor and the loader DLL is present in the bundle

## Pass Criteria

EQ8 validation is successful when:

- the plugin loads in a real host
- the editor opens
- parameter changes work
- state restore works
- reopen and rescan remain stable
- the macOS or Windows host checklist above passes for the target OS
