// Package plugin provides the runtime-facing VST3 plugin interfaces.
//
// This package is the contract consumed by the VST3 wrapper/runtime layer:
// callers implement Plugin, which produces a Processor, and may optionally
// implement StatefulProcessor when parameter state alone is not sufficient.
package plugin

import (
	"io"

	"github.com/th-release/vst3go/pkg/framework/bus"
	"github.com/th-release/vst3go/pkg/framework/param"
	"github.com/th-release/vst3go/pkg/framework/plugin"
	"github.com/th-release/vst3go/pkg/framework/process"
)

// Plugin is the top-level runtime interface implemented by a plugin entrypoint.
//
// GetInfo returns the metadata shared with the host, while CreateProcessor
// constructs a fresh Processor instance for actual audio/event handling.
type Plugin interface {
	// GetInfo returns plugin metadata
	GetInfo() plugin.Info

	// CreateProcessor creates a new instance of the audio processor
	CreateProcessor() Processor
}

// Processor is the core runtime contract consumed by the VST3 wrapper.
//
// A Processor owns:
//   - initialization and activation lifecycle
//   - audio/event processing via process.Context
//   - parameter and bus configuration
//   - latency/tail reporting
//
// The wrapper/runtime creates a Processor from Plugin.CreateProcessor and
// delegates all host-facing lifecycle and processing calls through it.
type Processor interface {
	// Initialize is called when the plugin is created
	Initialize(sampleRate float64, maxBlockSize int32) error

	// ProcessAudio processes audio - ZERO ALLOCATIONS!
	ProcessAudio(ctx *process.Context)

	// GetParameters returns the parameter registry
	GetParameters() *param.Registry

	// GetBuses returns the bus configuration
	GetBuses() *bus.Configuration

	// SetActive is called when processing starts/stops
	SetActive(active bool) error

	// GetLatencySamples returns the plugin's latency in samples
	GetLatencySamples() int32

	// GetTailSamples returns the tail length in samples
	GetTailSamples() int32
}

// StatefulProcessor extends Processor with custom state save/load hooks.
//
// Implement this only when parameter state is insufficient to restore runtime
// behavior, for example when additional processor-owned state must be persisted.
type StatefulProcessor interface {
	Processor

	// SaveCustomState saves additional state beyond parameters
	// This is called after all parameters have been saved
	SaveCustomState(w io.Writer) error

	// LoadCustomState loads additional state beyond parameters
	// This is called after all parameters have been loaded
	LoadCustomState(r io.Reader) error
}
