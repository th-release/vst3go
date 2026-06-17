# Developing With `vst3go`

This guide shows how to build a plugin on top of `vst3go`, from the first `Plugin`/`Processor` pair all the way to editor metadata, state restore, and release validation.

The short version is:

- `vst3go` is the runtime and VST3 shell.
- Your downstream repo owns the actual effect/instrument behavior.
- Keep DSP logic in Go where possible, but keep it outside this repository if it is product-specific.

If you want a mental model, think of `vst3go` as the layer that answers:

- What is the plugin called?
- What parameters does it expose?
- What audio and event buses does it have?
- How does it process one block of audio?
- How does it save and restore state?
- How does the native editor reflect and update parameters?

It does not try to be a full DAW framework.

## 1. The Three Layers You Need To Understand

When you build a plugin with this repo, there are usually three separate layers:

1. **Metadata layer**
   - `pkg/framework/plugin.Info` contains the name, ID, vendor, version, and category.
   - `pkg/plugin.Plugin.GetInfo()` returns that metadata to the host.

2. **Runtime layer**
   - `pkg/plugin.Processor` is the main contract.
   - It owns parameters, buses, lifecycle, and block processing.
   - The wrapper in `pkg/plugin` turns this into a VST3 component/controller pair.

3. **Editor/state layer**
   - `pkg/plugin.EditorModel` and `pkg/plugin.EditorSnapshot` describe what the browser-rendered editor sees.
   - `pkg/framework/state.Manager` and `StatefulProcessor` handle persistence beyond parameter values.
   - The web editor shell is only a view of the same parameter/state model.

The important rule is that the audio processor is the source of truth. The editor mirrors it, not the other way around.

## 2. What A Plugin Looks Like In Practice

At minimum, a plugin needs:

- a `pkg/plugin.Plugin` implementation
- a `pkg/plugin.Processor` implementation
- a `pkg/framework/param.Registry`
- a `pkg/framework/bus.Configuration`
- optional custom state hooks if parameter values are not enough

The runtime wrapper uses the following flow:

1. host asks for plugin metadata
2. host creates a processor instance
3. wrapper calls `Initialize(sampleRate, maxBlockSize)`
4. host starts processing
5. `ProcessAudio(ctx)` gets called repeatedly
6. editor/view asks for current parameter state when needed
7. host saves or restores state through `GetState` and `SetState`

## 3. A Minimal Plugin Skeleton

This is the shape you usually want in a downstream repo:

```go
package myplugin

import (
	"github.com/cwbudde/vst3go/pkg/framework/bus"
	frameworkplugin "github.com/cwbudde/vst3go/pkg/framework/plugin"
	"github.com/cwbudde/vst3go/pkg/framework/param"
	"github.com/cwbudde/vst3go/pkg/framework/process"
	"github.com/cwbudde/vst3go/pkg/plugin"
)

type Plugin struct{}

func (p *Plugin) GetInfo() frameworkplugin.Info {
	return frameworkplugin.Info{
		ID:       "com.example.myplugin",
		Name:     "My Plugin",
		Version:  "0.1.0",
		Vendor:   "Example Co",
		Category: "Fx",
	}
}

func (p *Plugin) CreateProcessor() plugin.Processor {
	return NewProcessor()
}

type Processor struct {
	params *param.Registry
	buses  *bus.Configuration
}
```

That part is boring on purpose. The important work is in the processor itself.

## 4. Build The Processor Around Real Parameters

Parameters should be designed first, not added as an afterthought.

For each control, decide:

- is this a continuous or discrete parameter?
- should it be automatable?
- should it be hidden?
- should it be read-only?
- what is the plain range?
- how should it be formatted in the editor?

The builder helpers in `pkg/framework/param` are meant to keep this clean:

- `GainParameter`
- `FrequencyParameter`
- `MixParameter`
- `TimeParameter`
- `RatioParameter`
- `QParameter`
- `PanParameter`
- `PhaseParameter`
- `Choice`

### Example Parameter Block

For a 3-band EQ, a practical parameter set might be:

- `Input Gain`
- `Low Gain`
- `Low Frequency`
- `Low Q`
- `Mid Gain`
- `Mid Frequency`
- `Mid Q`
- `High Gain`
- `High Frequency`
- `Output Gain`
- `Bypass`

Example setup:

```go
func NewProcessor() *Processor {
	params := param.NewRegistry()

	_ = params.Add(
		param.GainParameter(1, "Input Gain").Default(0).Build(),
		param.GainParameter(2, "Low Gain").Default(0).Build(),
		param.FrequencyParameter(3, "Low Frequency", 20, 2000, 120).Build(),
		param.QParameter(4, "Low Q", 0.1, 8.0, 0.707).Build(),
		param.GainParameter(5, "Mid Gain").Default(0).Build(),
		param.FrequencyParameter(6, "Mid Frequency", 200, 5000, 1000).Build(),
		param.QParameter(7, "Mid Q", 0.1, 12.0, 1.0).Build(),
		param.GainParameter(8, "High Gain").Default(0).Build(),
		param.FrequencyParameter(9, "High Frequency", 1000, 20000, 8000).Build(),
		param.GainParameter(10, "Output Gain").Default(0).Build(),
		param.New(11, "Bypass").Toggle().Bypass().Build(),
	)

	return &Processor{
		params: params,
		buses:  bus.Stereo(),
	}
}
```

For a `Utility`-style plugin, the parameter set is different:

- `Input Gain`
- `Output Gain`
- `Pan`
- `Width`
- `Mono`
- `Phase Invert L`
- `Phase Invert R`
- `Mute`

Those are all stock controls that users already understand, which is exactly why the plugin feels like a utility tool instead of a special effect.

## 5. Choose Your Bus Layout First

The bus layout determines how the host sees your plugin.

Use the helpers in `pkg/framework/bus`:

- `bus.Stereo()` for the common stereo-in/stereo-out effect
- `bus.Mono()` for a simple single-channel utility or instrument shell
- `bus.NewBuilder()` for sidechain or multibus layouts

### Stereo EQ

Most EQs should start with a stereo main input and output:

```go
buses := bus.Stereo()
```

### Utility Plugin

If you are building a Utility-like tool, stereo I/O is still the normal choice:

```go
buses := bus.Stereo()
```

Even if the processing is simple, keep the same I/O shape as the host expects. A utility plugin usually manipulates the stereo image, gain, and phase, so stereo in/out is the right default.

### Sidechain Example

If you later want ducking or external analysis, use a builder:

```go
buses := bus.NewBuilder().
	WithStereoInput("Main In").
	WithStereoOutput("Main Out").
	WithSidechain("Sidechain In").
	MustBuild()
```

The guide rule here is simple: define the buses from the user-facing behavior backward, not from the implementation forward.

## 6. Implement `Initialize`, `SetActive`, And `ProcessAudio`

The lifecycle matters because `vst3go` is designed around realtime audio constraints.

### `Initialize`

Use `Initialize(sampleRate, maxBlockSize)` to allocate or configure anything that depends on host sample rate or block size.

Typical responsibilities:

- initialize filter coefficients
- reset envelopes or delay lines
- cache sample-rate-dependent constants
- prepare lookup tables

Do not do host-independent plugin setup here if it can happen earlier.

### `SetActive`

Use `SetActive(true)` when processing starts, and `SetActive(false)` when the host stops or deactivates the plugin.

Typical responsibilities:

- reset smoothing
- clear delay lines if needed
- clear pending MIDI/events
- reset transport followers
- prepare the processor for a clean next start

### `ProcessAudio`

This is the realtime callback. It should be:

- allocation-free
- lock-free where possible
- branch-light
- deterministic

The `process.Context` gives you convenient access to:

- `Input`
- `Output`
- `NumSamples()`
- `NumInputChannels()`
- `NumOutputChannels()`
- `Param(id)`
- `ParamPlain(id)`
- `ProcessChannels`
- `ProcessStereo`
- `ProcessSamples`
- `PassThrough`
- `Clear`

### A Good Default Pattern

For a stereo EQ, a clean starting point is:

```go
func (p *Processor) ProcessAudio(ctx *process.Context) {
	if ctx == nil {
		return
	}

	if ctx.Param(11) >= 0.5 {
		ctx.PassThrough()
		return
	}

	ctx.ProcessStereo(func(ch int, input, output []float32) {
		for i := range input {
			sample := input[i]
			sample = p.processLow(sample, ctx.ParamPlain(2), ctx.ParamPlain(3), ctx.ParamPlain(4))
			sample = p.processMid(sample, ctx.ParamPlain(5), ctx.ParamPlain(6), ctx.ParamPlain(7))
			sample = p.processHigh(sample, ctx.ParamPlain(8), ctx.ParamPlain(9))
			sample = applyGain(sample, ctx.ParamPlain(1), ctx.ParamPlain(10))
			output[i] = sample
		}
	})
}
```

The filter functions above are your downstream DSP code. `vst3go` gives you the shell and the block plumbing, not a built-in EQ engine.

## 7. How To Build A 3-Band EQ

This is the most natural example because it exercises almost everything:

- continuous parameters
- tonal controls
- sample-rate-dependent math
- realtime audio processing
- editor controls
- state persistence

### 7.1 Decide The Control Model

A practical 3-band EQ usually has:

- input gain
- low shelf gain and frequency
- mid peaking gain, frequency, and Q
- high shelf gain and frequency
- output gain
- bypass

You may also want:

- output meter
- phase invert
- oversampling toggle
- analyzer enable

Keep the first version small. You can add advanced extras later.

### 7.2 Decide The Filter Topology

In a downstream repo, you will likely implement:

- a low shelf biquad
- a peaking mid band biquad
- a high shelf biquad

For each channel, you usually keep coefficient/state structs like:

```go
type biquad struct {
	b0, b1, b2 float64
	a1, a2     float64
	z1, z2     float64
}
```

And then:

```go
func (b *biquad) Process(x float64) float64 {
	y := b.b0*x + b.z1
	b.z1 = b.b1*x - b.a1*y + b.z2
	b.z2 = b.b2*x - b.a2*y
	return y
}
```

That math is downstream-specific, but the important part for `vst3go` is where you store it:

- processor fields
- reconfigured in `Initialize`
- reset in `SetActive(false)`
- driven from current parameter values during `ProcessAudio`

### 7.3 Keep Parameter Reads Cheap

In `ProcessAudio`, read values once per block unless you truly need per-sample automation.

Example:

```go
inputGain := ctx.ParamPlain(1)
lowGain := ctx.ParamPlain(2)
lowFreq := ctx.ParamPlain(3)
lowQ := ctx.ParamPlain(4)
```

If you want sample-accurate automation, use the process context’s parameter-change flow and update coefficients or targets at the relevant offsets.

### 7.4 Handle Bypass Cleanly

Bypass should usually be a true bypass path:

```go
if ctx.Param(11) >= 0.5 {
	ctx.PassThrough()
	return
}
```

For EQs, that is usually the most expected behavior.

### 7.5 Add State Later

If your EQ only uses parameters and no hidden state, `StatefulProcessor` may be unnecessary.

If you add features like:

- analyzer freeze
- preset name
- custom band modes
- user-defined routing

then add `SaveCustomState` and `LoadCustomState`.

## 8. How To Build A Utility Plugin

An Ableton Utility-style plugin is a better example for “small but useful” host behavior.

Common Utility controls:

- gain
- pan
- width
- mono on/off
- phase invert L/R
- mute

### 8.1 Parameters

Use plain, obvious values:

```go
_ = params.Add(
	param.GainParameter(1, "Input Gain").Default(0).Build(),
	param.GainParameter(2, "Output Gain").Default(0).Build(),
	param.PanParameter(3, "Pan").Build(),
	param.MixParameter(4, "Width").Default(100).Build(),
	param.New(5, "Mono").Toggle().Build(),
	param.New(6, "Phase Invert L").Toggle().Build(),
	param.New(7, "Phase Invert R").Toggle().Build(),
	param.New(8, "Mute").Toggle().Build(),
)
```

You can make `Width` a `MixParameter` or a custom `0..200%` control depending on how you want to define stereo expansion.

### 8.2 Processing Strategy

The processing order usually looks like:

1. apply input gain
2. optionally collapse to mono
3. apply pan
4. apply width/stereo matrix
5. apply phase inversion
6. apply output gain
7. zero the output if muted

That gives a predictable mental model for users.

### 8.3 Example Processing Shape

```go
func (p *Processor) ProcessAudio(ctx *process.Context) {
	if ctx.Param(8) >= 0.5 {
		ctx.Clear()
		return
	}

	inputGain := ctx.ParamPlain(1)
	outputGain := ctx.ParamPlain(2)
	pan := ctx.ParamPlain(3)
	width := ctx.ParamPlain(4)
	mono := ctx.Param(5) >= 0.5
	invertLeft := ctx.Param(6) >= 0.5
	invertRight := ctx.Param(7) >= 0.5

	ctx.ProcessStereo(func(ch int, input, output []float32) {
		for i := range input {
			sample := input[i]

			if ch == 0 && invertLeft {
				sample = -sample
			}
			if ch == 1 && invertRight {
				sample = -sample
			}

			sample = applyGain(sample, inputGain, 0)

			if mono {
				// Downstream implementation would usually average both channels.
				// This branch is just the structural idea.
			}

			if ch == 0 {
				sample = applyPanLeft(sample, pan, width)
			} else {
				sample = applyPanRight(sample, pan, width)
			}

			sample = applyGain(sample, outputGain, 0)
			output[i] = sample
		}
	})
}
```

The exact gain math is up to you. The point is to keep the plugin behavior user-facing and simple.

## 9. Editor Integration

The editor path in this repo is built on `EditorModel` and `EditorSnapshot`.

That means your plugin should think in terms of:

- parameter metadata
- current normalized values
- plain values for display
- hidden/read-only semantics

### 9.1 What The Editor Receives

The editor model is built from the parameter registry and plugin metadata.

Your job is to make sure the registry is clean:

- correct IDs
- correct names
- correct min/max/default values
- correct flags
- sensible short names and units

### 9.2 What The Editor Should Not Do

The editor should not invent its own state source.

Instead:

- the host edits parameters
- the runtime updates the registry
- the editor mirrors that registry
- the snapshot can be saved and restored

### 9.3 Hidden And Read-Only Controls

Use these flags deliberately:

- `IsHidden` for internal or unsupported controls
- `IsReadOnly` for meters or values that should display but not be edited
- `CanAutomate` for values the host may write to

If you add a meter, you usually want it read-only and not automatable.

## 10. Saving And Restoring State

If parameter values are enough, you may not need custom state at all.

If not, implement `StatefulProcessor`:

```go
type Processor struct {
	params *param.Registry
	buses  *bus.Configuration

	presetName string
	lastMode   int32
}

func (p *Processor) SaveCustomState(w io.Writer) error {
	// write custom fields after parameter state
	return nil
}

func (p *Processor) LoadCustomState(r io.Reader) error {
	// restore custom fields after parameter state
	return nil
}
```

Use custom state for things like:

- preset names
- mode selections that are not parameterized
- internal routing choices
- analysis display settings

Do not use custom state to duplicate regular parameter values unless you have a strong reason.

## 11. A Recommended File Layout For A Downstream Plugin Repo

If you are building a real plugin in another repo, a good layout is:

```text
my-plugin/
  cmd/
  internal/
    dsp/
    ui/
  plugin/
    processor.go
    plugin.go
    state.go
  web/
    editor/
      index.html
      assets/
  justfile
```

Suggested division:

- `plugin/processor.go` owns the runtime processor
- `plugin/plugin.go` owns the `Plugin` entrypoint and metadata
- `plugin/state.go` owns optional custom state
- `internal/dsp/` owns the filter/math code
- `web/editor/` owns the browser-rendered UI if you generate one

Keep `vst3go` itself as the runtime dependency, not the product workspace.

## 12. How To Think About A Real Plugin Project

When you start a new plugin, ask these questions in order:

1. What is the user trying to do?
2. What are the 5 to 12 controls that actually matter?
3. Is this mono, stereo, sidechain, or multibus?
4. What is the processing topology?
5. What state must survive reload?
6. What must the editor show immediately?
7. What needs to be automatable?
8. What should be hidden?

That sequence usually leads to a sane first version.

## 13. Suggested Development Order

For a new plugin, do the work in this order:

1. define plugin metadata
2. define buses
3. define parameters
4. implement processor lifecycle
5. implement the basic audio algorithm
6. add state save/restore if needed
7. expose the editor model
8. connect the editor shell
9. test with `just test`
10. run Windows validation if Windows is in scope

## 14. A Practical Example Workflow

Here is a realistic workflow for a downstream repo:

1. Start with a stereo utility plugin.
2. Verify gain, pan, width, and mute in the editor.
3. Add state save/restore for the last selected mode.
4. Turn the utility into a 3-band EQ by swapping the DSP layer.
5. Keep the parameter contract stable while the backend changes.
6. Reuse the same web editor shell and only change the data fed into it.

That sequence keeps the project from becoming a pile of disconnected ideas.

## 15. What To Avoid

- do not put product-specific DSP into `vst3go`
- do not treat the browser editor as a separate app
- do not bypass the parameter registry to keep extra state
- do not allocate in `ProcessAudio` unless you have measured and accepted the tradeoff
- do not make platform glue own business logic
- do not hide important user controls behind undocumented behavior

## 16. Validation Checklist For A New Plugin

Before you call a plugin “done”, make sure:

- it loads in the host
- it processes audio without crash or glitch
- its parameters reflect correctly in the editor
- its state survives save and reload
- its bus layout matches the intended use case
- its realtime path is allocation-conscious
- its build and packaging story are documented

If Windows is in scope, also verify:

- the DLL builds
- the bundle layout is correct
- the editor opens in WebView2
- repeated reopen cycles are stable
- the plugin survives host reloads and rescans

## 17. Closing Advice

The best `vst3go` plugin code is usually very boring in the right places:

- metadata is explicit
- parameters are named carefully
- buses are obvious
- processing is deterministic
- state is intentional
- the editor mirrors the runtime instead of replacing it

If you keep those rules, you can build anything from a straightforward EQ to a polished Utility-style plugin without fighting the framework.
