# VST3Go TODO

## Repo Boundary

- [x] Keep `vst3go` scoped as the VST3 binding/runtime layer.
- [x] Move DSP packages, showcase plugins, and synth-product planning out to `synthkit`.
- [x] Rename the module to `github.com/cwbudde/vst3go`.
- [x] Update the README to state that this fork intentionally differs from the original project.
- [x] Review the retained package boundary once more.
  - [x] `pkg/vst3`
  - [x] `pkg/plugin`
  - [x] `pkg/midi`
  - [x] `pkg/framework/{bus,param,plugin,process,state}`
  - [x] Confirm no higher-level audio/product concerns have drifted back in.

## Rules

- [x] Keep the C bridge minimal and direct.
- [x] Keep business logic in Go rather than in C.
- [x] Keep DSP and product-specific logic out of this repo.
- [x] Keep example/showcase plugin ownership out of this repo.
- [x] Keep these rules reflected in docs and future planning updates.

## Current Baseline

- [x] Plugin implementations can import `pkg/plugin/cbridge` instead of including bridge code directly.
- [x] Core runtime packages exist for buses, parameters, process context, state, MIDI/events, plugin runtime, and VST3 bindings.
- [x] Sample-accurate parameter automation is implemented.
- [x] MIDI/event handling is implemented.
- [x] Advanced bus support is implemented.
- [x] Parameter IDs are standardized around typed `uint32` constants with `ParamXxx` naming.
- [x] Most runtime reset paths are in place on deactivation.
- [x] `go test ./pkg/...` passes after the split and module rename.
- [x] Clean up the remaining bridge warning.
  - [x] Address the `const` qualifier mismatch around `GoCreateInstance` in `bridge.c` / `bridge.h`.

## Phase 1: Runtime API Cleanup

- [ ] Resolve naming inconsistencies.
  - [x] Eliminate duplicate constant spellings and casing variants by making `ResultOK` canonical and keeping aliases compatibility-only.
  - [x] Standardize normalized/plain parameter naming semantics by adding explicit `GetNormalized`/`SetNormalized` and `GetPlain`/`SetPlain` accessors.
  - [x] Standardize canonical bus configuration entry points by adding primary constructors such as `Stereo`, `Mono`, `EffectStereo`, `Generator`, and `MIDIEffect`.
- [ ] Clarify interface hierarchy.
  - [x] Document the relationships between plugin, processor, and optional stateful processor behavior in code comments.
  - [x] Keep optional interfaces explicit rather than implicit.
- [ ] Tighten parameter and registry APIs.
  - [x] Audit and clean up the obvious unsafe parameter access patterns in retained runtime code.
  - [x] Prefer checked access where nil is possible by adding safe registry helpers such as `GetOK`, `Has`, `GetNormalized`, and `GetPlain`.
  - [x] Clarify duplicate-registration behavior and failure modes by making `Registry.Add` fail explicitly on duplicate IDs.
- [ ] Tighten error handling.
  - [x] Replace silent failures where normal error returns are appropriate in the state/registry layer.
  - [x] Wrap errors with useful context in state save/load and core runtime setup paths.

## Phase 2: Runtime Architecture Hardening

- [x] Reassess package boundaries.
  - [x] Keep VST3-facing code, bridge glue, and generic runtime primitives clearly separated.
  - [x] Remove mixed-responsibility code where low-level and higher-level concerns are entangled.
- [x] Review lifecycle behavior.
  - [x] Audit `SetActive(false)` and related reset behavior in retained runtime code.
  - [x] Ensure state, parameter, and event handling reset deterministically in retained runtime paths.
- [x] Review process-layer structure.
  - [x] Refactor oversized runtime processing functions where readability is poor.
  - [x] Split transport updates, buffer mapping, and automation collection into smaller runtime helpers.
- [x] Keep helper abstractions honest.
  - [x] Only add base helpers if they reduce boilerplate without hiding core VST3 mechanics.

## Phase 3: Testing, Validation, And Performance

- [ ] Expand automated coverage for retained packages.
  - [ ] `pkg/framework/param`
  - [ ] `pkg/framework/process`
  - [ ] `pkg/framework/state`
  - [ ] `pkg/midi`
  - [ ] `pkg/plugin` where practical
  - [ ] `pkg/vst3`
- [ ] Add race-detector coverage where thread-safety assumptions are non-trivial.
- [ ] Keep allocation-sensitive paths measurable.
  - [ ] Benchmark process-context internals.
  - [ ] Benchmark runtime hot paths where useful.
  - [ ] Revisit fixed-size internal buffers where limits may be too rigid.
- [ ] Keep validation expectations explicit.
  - [ ] Define minimal validation expectations for the runtime layer itself.
  - [ ] Document what downstream repos such as `synthkit` must validate on their own.

## Phase 4: Documentation And Consumer Experience

- [ ] Document retained scope clearly.
  - [x] State the split boundary in README and PLAN.
  - [ ] Add a concise “what belongs here vs. in `synthkit`” section if the README still needs tightening.
- [ ] Improve API documentation.
  - [ ] package-level docs
  - [ ] lifecycle expectations
  - [ ] thread-safety guarantees
  - [ ] state/persistence expectations
- [ ] Add migration guidance for the fork.
  - [x] Clarify intentional differences from the original project.
  - [x] Clarify the new module path.
  - [ ] Document how downstream repos should consume this fork once versions are tagged.
- [ ] Keep build and validation instructions accurate for the retained runtime repo.

## Phase 5: Web-Rendered Plugin UI

- [ ] Define the editor surface and render scope.
  - [ ] Identify the current plugin UI entry points and what `createView` should expose.
  - [ ] Decide which controls, layout blocks, and parameter bindings must be rendered first.
- [ ] Make the editor renderable through a web-based surface.
  - [ ] Add a web-rendered view path for the VST editor instead of a standalone site.
  - [ ] Keep the rendered UI aligned with the current plugin design and controls.
  - [ ] Make local editing and parameter changes visible immediately in the view.
  - [ ] Keep an explicit editor snapshot for save/restore and future React hydration.
  - [ ] Document how HTML, CSS, and JS build outputs plug into the editor shell.
- [ ] Treat the web-rendered editor as the prerequisite for platform expansion work.

## Phase 6: Cross-Platform Windows Support

- [ ] Update the build system for Windows.
  - [ ] Detect Windows cleanly.
  - [ ] Produce the correct library format and bundle layout.
- [ ] Extend the bridge for Windows entry points.
  - [x] Add required DLL entry handling.
  - [x] Export `GetPluginFactory`.
  - [ ] Host the web-rendered editor with WebView2 on Windows.
  - [x] Lay down the Windows editor-view scaffold.
- [ ] Review platform-specific CGO directives.
  - [ ] Add Windows-specific flags only where necessary.
  - [ ] Minimize platform divergence.
- [ ] Validate on real toolchains and hosts.
  - [ ] Start with MinGW-w64.
  - [ ] Keep MSVC as a later enhancement if justified.
- [ ] Capture Windows-specific risks.
  - [ ] Shared-library behavior with Go.
  - [ ] Scheduler differences.
  - [ ] File locking and reload friction.
  - [ ] Path-length edge cases.

## Phase 7: Release Readiness

- [ ] Keep architectural guardrails enforced.
- [ ] Keep retained runtime packages tested and documented.
- [ ] Confirm the public API boundary is stable enough for downstream repos.
- [ ] Decide supported release scope explicitly.
  - [ ] Linux/macOS only, or include Windows.
  - [ ] Runtime-only, with higher-level features delegated to companion repos.

## Ongoing Rules

- [ ] If work belongs in `synthkit`, remove it from this file instead of duplicating it here.
- [ ] Avoid standalone roadmap documents; merge planning back here.
- [ ] Before adding a new abstraction, check whether it belongs in `vst3go` or in a higher-level companion repo.
