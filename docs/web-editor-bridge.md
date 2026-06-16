# Web Editor Bridge

This repo now treats the VST editor as a browser-rendered surface instead of a standalone website.

## Rendering Flow

1. `pkg/plugin/component.go` builds an `EditorModel` from the live parameter registry.
2. `pkg/plugin/wrapper_controller.go` serializes that model into the HTML bootstrap returned by `GoEditControllerGetEditorHTML`.
3. `pkg/plugin/editor_view_darwin.c` loads the HTML into `WKWebView` on macOS.
4. Parameter changes flow back through the web message handler and end up in the Go processor state.

## HTML, CSS, And JS

The current bridge inlines the shell HTML, CSS, and JS so the editor can work without a bundler.

When you switch to a React build, keep the same data contract:

- `EditorModel` or `EditorSnapshot` stays the source of truth in Go.
- The build output should replace the inline UI shell, not the data contract.
- CSS can remain external or be inlined by the build pipeline.
- JS should still receive the initial snapshot from Go, then hydrate the controls from that data.

Recommended build layout for a React editor:

- `web/editor/index.html` — shell page
- `web/editor/assets/*.css` — extracted styles
- `web/editor/assets/*.js` — bundled React runtime and editor app

The shell only needs a bootstrap hook that injects the encoded snapshot before the app mounts.

## State Save And Restore

There are two state layers:

- **Plugin state**: managed by `pkg/framework/state.Manager` and the VST3 `GetState` / `SetState` callbacks.
- **Editor snapshot**: the web editor's live parameter snapshot, which can be saved locally and restored into the web UI.

Current behavior:

- The editor snapshot mirrors the live parameter registry.
- `EditorSnapshot.Apply` restores snapshot values back into the registry.
- The web editor can save a snapshot in local storage and restore it back into the controls.

This keeps the host-controlled plugin state and the editor-controlled view state aligned.

## Windows Bridge Design

Windows should reuse the same snapshot contract and HTML bootstrap.

Recommended shape:

- Add a Windows editor bridge beside the macOS one.
- Back the editor window with WebView2.
- Keep the Go-side snapshot and parameter callback flow unchanged.
- Treat window host integration as platform glue only.
- Avoid duplicating editor data structures for Windows.

The Windows implementation should focus on:

- `HWND` / controller attachment
- WebView2 initialization and message routing
- resize and focus handling
- `GetPluginFactory` export plumbing
- build tags and CGO flags for the Windows toolchain

## Why This Shape

This approach keeps the editor surface portable while preserving the repo's boundary:

- Go owns the plugin data model.
- Platform code only hosts the browser renderer.
- The browser renderer only speaks snapshot JSON and parameter updates.
