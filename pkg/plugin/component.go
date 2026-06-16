package plugin

// #cgo CFLAGS: -I../../include
// #include "../../include/vst3/vst3_c_api.h"
// #include "../../bridge/bridge.h"
import "C"

import (
	"bytes"
	"fmt"
	"sync"
	"unsafe"

	"github.com/cwbudde/vst3go/pkg/framework/bus"
	"github.com/cwbudde/vst3go/pkg/framework/param"
	"github.com/cwbudde/vst3go/pkg/framework/plugin"
	"github.com/cwbudde/vst3go/pkg/framework/process"
	"github.com/cwbudde/vst3go/pkg/framework/state"
	"github.com/cwbudde/vst3go/pkg/midi"
	"github.com/cwbudde/vst3go/pkg/vst3"
)

// componentImpl wraps a Processor to implement VST3 interfaces
type componentImpl struct {
	processor    Processor
	pluginInfo   plugin.Info
	processCtx   *process.Context
	sampleRate   float64
	maxBlockSize int32
	active       bool
	processing   bool
	mu           sync.RWMutex
	wrapper      *componentWrapper // Reference to wrapper for notifications
}

func (c *componentImpl) parameterRegistry() (*param.Registry, error) {
	if c.processor == nil {
		return nil, fmt.Errorf("component has no processor")
	}

	params := c.processor.GetParameters()
	if params == nil {
		return nil, fmt.Errorf("processor has no parameter registry")
	}

	return params, nil
}

func (c *componentImpl) GetParameters() *param.Registry {
	if c.processor == nil {
		return nil
	}

	return c.processor.GetParameters()
}

func (c *componentImpl) stateManager() (*state.Manager, error) {
	params, err := c.parameterRegistry()
	if err != nil {
		return nil, err
	}

	manager := state.NewManager(params)
	if stateful, ok := c.processor.(StatefulProcessor); ok {
		manager.SetCustomLoadFunc(stateful.LoadCustomState)
		manager.SetCustomSaveFunc(stateful.SaveCustomState)
	}

	return manager, nil
}

// newComponent creates a new component implementation
func newComponent(processor Processor, info plugin.Info) *componentImpl {
	params := processor.GetParameters()
	return &componentImpl{
		processor:    processor,
		pluginInfo:   info,
		processCtx:   process.NewContext(8192, params), // Default max block size
		maxBlockSize: 8192,
	}
}

// IComponent implementation
func (c *componentImpl) Initialize(_ interface{}) error {
	return c.processor.Initialize(48000, c.maxBlockSize) // Default sample rate
}

func (c *componentImpl) Terminate() error {
	return nil
}

func (c *componentImpl) GetControllerClassID() [16]byte {
	// Return same ID - we're a single component
	return [16]byte{}
}

func (c *componentImpl) SetIOMode(_ int32) error {
	return nil
}

func (c *componentImpl) GetBusCount(mediaType, direction int32) int32 {
	buses := c.processor.GetBuses()
	return buses.GetBusCount(bus.MediaType(mediaType), bus.Direction(direction))
}

func (c *componentImpl) GetBusInfo(mediaType, direction, index int32) (*vst3.BusInfo, error) {
	buses := c.processor.GetBuses()
	info := buses.GetBusInfo(bus.MediaType(mediaType), bus.Direction(direction), index)
	if info == nil {
		return nil, vst3.ErrNotImplemented
	}

	flags := uint32(1) // Default active
	if !info.IsActive {
		flags = 0
	}

	return &vst3.BusInfo{
		MediaType:    int32(info.MediaType),
		Direction:    int32(info.Direction),
		ChannelCount: info.ChannelCount,
		Name:         info.Name,
		BusType:      int32(info.BusType),
		Flags:        flags,
	}, nil
}

func (c *componentImpl) GetRoutingInfo(inInfo, outInfo interface{}) error {
	return vst3.ErrNotImplemented
}

func (c *componentImpl) ActivateBus(mediaType, direction, index int32, state bool) error {
	return nil
}

func (c *componentImpl) SetActive(active bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.active = active
	if !active && c.processCtx != nil {
		c.processCtx.ResetParameterChanges()
		c.processCtx.ClearAllEvents()
	}
	return c.processor.SetActive(active)
}

func (c *componentImpl) SetState(stateData []byte) error {
	stateManager, err := c.stateManager()
	if err != nil {
		return fmt.Errorf("set component state: %w", err)
	}

	buf := bytes.NewReader(stateData)
	if err := stateManager.Load(buf); err != nil {
		return fmt.Errorf("set component state: %w", err)
	}
	return nil
}

func (c *componentImpl) GetState() ([]byte, error) {
	stateManager, err := c.stateManager()
	if err != nil {
		return nil, fmt.Errorf("get component state: %w", err)
	}

	var buf bytes.Buffer
	if err := stateManager.Save(&buf); err != nil {
		return nil, fmt.Errorf("get component state: %w", err)
	}

	return buf.Bytes(), nil
}

// IAudioProcessor implementation
func (c *componentImpl) SetBusArrangements(inputs, outputs []int64) error {
	return nil
}

func (c *componentImpl) GetBusArrangement(direction, index int32) (int64, error) {
	// Return stereo by default
	return int64(3), nil // Left + Right
}

func (c *componentImpl) CanProcessSampleSize(symbolicSampleSize int32) error {
	// We only support 32-bit float
	if symbolicSampleSize == 0 { // kSample32
		return nil
	}
	return vst3.ErrNotImplemented
}

func (c *componentImpl) GetLatencySamples() uint32 {
	return uint32(c.processor.GetLatencySamples())
}

func (c *componentImpl) SetupProcessing(setup *vst3.ProcessSetup) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if setup == nil {
		return fmt.Errorf("setup processing: nil process setup")
	}

	c.sampleRate = setup.SampleRate
	if setup.MaxSamplesPerBlock > 0 {
		c.maxBlockSize = setup.MaxSamplesPerBlock
		// Recreate process context with new max block size
		params := c.processor.GetParameters()
		c.processCtx = process.NewContext(int(c.maxBlockSize), params)
	}

	if err := c.processor.Initialize(c.sampleRate, c.maxBlockSize); err != nil {
		return fmt.Errorf("setup processing: initialize processor: %w", err)
	}
	return nil
}

func (c *componentImpl) SetProcessing(state bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.processing = state
	return nil
}

func (c *componentImpl) Process(data unsafe.Pointer) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if data == nil {
		return fmt.Errorf("process: nil process data")
	}

	if !c.processing {
		return nil
	}

	// Get raw process data struct
	processData := (*C.struct_Steinberg_Vst_ProcessData)(data)

	c.prepareProcessContext()
	c.updateTransport(processData.processContext)
	c.mapAudioBuffers(processData)
	c.processInputEvents(processData.inputEvents)
	c.collectParameterChanges(unsafe.Pointer(processData.inputParameterChanges))

	// Process audio with sample-accurate parameter automation
	if c.processCtx.HasParameterChanges() {
		// Sort parameter changes by sample offset
		c.processCtx.SortParameterChanges()

		// Process audio in chunks between parameter changes
		c.processSampleAccurate()
	} else {
		// No parameter changes - process entire block
		c.processor.ProcessAudio(c.processCtx)
	}

	return nil
}

func (c *componentImpl) prepareProcessContext() {
	c.processCtx.SampleRate = c.sampleRate
	c.processCtx.ResetParameterChanges()
	c.processCtx.ClearAllEvents()
	c.processCtx.Input = c.processCtx.Input[:0]
	c.processCtx.Output = c.processCtx.Output[:0]
}

func (c *componentImpl) updateTransport(ctx *C.struct_Steinberg_Vst_ProcessContext) {
	transport := c.processCtx.Transport
	*transport = process.TransportInfo{}

	if ctx == nil {
		return
	}

	transport.IsPlaying = (ctx.state & C.Steinberg_Vst_ProcessContext_StatesAndFlags_kPlaying) != 0
	transport.IsRecording = (ctx.state & C.Steinberg_Vst_ProcessContext_StatesAndFlags_kRecording) != 0
	transport.IsCycling = (ctx.state & C.Steinberg_Vst_ProcessContext_StatesAndFlags_kCycleActive) != 0

	transport.HasTempo = (ctx.state & C.Steinberg_Vst_ProcessContext_StatesAndFlags_kTempoValid) != 0
	if transport.HasTempo {
		transport.Tempo = float64(ctx.tempo)
	}

	transport.HasTimeSignature = (ctx.state & C.Steinberg_Vst_ProcessContext_StatesAndFlags_kTimeSigValid) != 0
	if transport.HasTimeSignature {
		transport.TimeSigNumerator = int32(ctx.timeSigNumerator)
		transport.TimeSigDenominator = int32(ctx.timeSigDenominator)
	}

	transport.HasMusicalTime = (ctx.state & C.Steinberg_Vst_ProcessContext_StatesAndFlags_kProjectTimeMusicValid) != 0
	if transport.HasMusicalTime {
		transport.ProjectTimeMusic = float64(ctx.projectTimeMusic)
	}

	transport.HasBarPosition = (ctx.state & C.Steinberg_Vst_ProcessContext_StatesAndFlags_kBarPositionValid) != 0
	if transport.HasBarPosition {
		transport.BarPositionMusic = float64(ctx.barPositionMusic)
	}

	transport.HasCycle = (ctx.state & C.Steinberg_Vst_ProcessContext_StatesAndFlags_kCycleValid) != 0
	if transport.HasCycle {
		transport.CycleStartMusic = float64(ctx.cycleStartMusic)
		transport.CycleEndMusic = float64(ctx.cycleEndMusic)
	}

	transport.ProjectTimeSamples = int64(ctx.projectTimeSamples)
	transport.ContinuousTimeSamples = int64(ctx.continousTimeSamples)

	if (ctx.state & C.Steinberg_Vst_ProcessContext_StatesAndFlags_kClockValid) != 0 {
		transport.SamplesToNextClock = int32(ctx.samplesToNextClock)
	}
}

func (c *componentImpl) mapAudioBuffers(processData *C.struct_Steinberg_Vst_ProcessData) {
	numSamples := int(processData.numSamples)

	if processData.numInputs > 0 && processData.inputs != nil {
		inputBuses := (*[1]C.struct_Steinberg_Vst_AudioBusBuffers)(unsafe.Pointer(processData.inputs))[:processData.numInputs:processData.numInputs]
		c.appendAudioBuffers(&c.processCtx.Input, inputBuses, numSamples)
	}

	if processData.numOutputs > 0 && processData.outputs != nil {
		outputBuses := (*[1]C.struct_Steinberg_Vst_AudioBusBuffers)(unsafe.Pointer(processData.outputs))[:processData.numOutputs:processData.numOutputs]
		c.appendAudioBuffers(&c.processCtx.Output, outputBuses, numSamples)
	}
}

func (c *componentImpl) appendAudioBuffers(dst *[][]float32, buses []C.struct_Steinberg_Vst_AudioBusBuffers, numSamples int) {
	for _, bus := range buses {
		channelBuffers32 := getChannelBuffers32(&bus)
		if bus.numChannels <= 0 || channelBuffers32 == nil {
			continue
		}

		channels := (*[16]*float32)(unsafe.Pointer(channelBuffers32))[:bus.numChannels:bus.numChannels]
		for _, channel := range channels {
			if channel == nil {
				continue
			}

			samples := (*[vst3.MaxArraySize]float32)(unsafe.Pointer(channel))[:numSamples:numSamples]
			*dst = append(*dst, samples)
		}
	}
}

func (c *componentImpl) collectParameterChanges(inputParameterChanges unsafe.Pointer) {
	if inputParameterChanges == nil {
		return
	}

	paramCount := C.getParameterChangeCount(inputParameterChanges)
	for i := C.int32_t(0); i < paramCount; i++ {
		paramQueue := C.getParameterData(inputParameterChanges, i)
		if paramQueue == nil {
			continue
		}

		paramID := C.getParameterId(paramQueue)
		pointCount := C.getPointCount(paramQueue)
		for j := C.int32_t(0); j < pointCount; j++ {
			var sampleOffset C.int32_t
			var value C.double

			result := C.getPoint(paramQueue, j, &sampleOffset, &value)
			if result == 0 {
				c.processCtx.AddParameterChange(uint32(paramID), float64(value), int(sampleOffset))
			}
		}
	}
}

func (c *componentImpl) GetTailSamples() uint32 {
	return uint32(c.processor.GetTailSamples())
}

// processInputEvents processes MIDI events from the host
func (c *componentImpl) processInputEvents(eventList *C.struct_Steinberg_Vst_IEventList) {
	if eventList == nil {
		return
	}

	// Get event count using helper
	eventCount := C.getEventCount(unsafe.Pointer(eventList))
	if eventCount <= 0 {
		return
	}

	// Process each event
	for i := C.int32_t(0); i < eventCount; i++ {
		var event C.struct_Steinberg_Vst_Event
		if C.getEvent(unsafe.Pointer(eventList), i, &event) == 0 { // kResultOk
			c.processSingleEvent(&event)
		}
	}
}

// processSingleEvent converts a VST3 event to our MIDI event format
func (c *componentImpl) processSingleEvent(event *C.struct_Steinberg_Vst_Event) {
	// Use helper function to get event type
	eventType := C.getEventType(event)

	switch eventType {
	case C.Steinberg_Vst_Event_EventTypes_kNoteOnEvent:
		// Note On event - use helper to get the event data
		noteOn := C.getNoteOnEvent(event)
		c.processCtx.AddInputEvent(midi.NoteOnEvent{
			BaseEvent: midi.BaseEvent{
				EventChannel: uint8(noteOn.channel),
				Offset:       int32(event.sampleOffset),
			},
			NoteNumber: uint8(noteOn.pitch),
			Velocity:   uint8(noteOn.velocity * 127), // VST3 uses 0-1, MIDI uses 0-127
		})

	case C.Steinberg_Vst_Event_EventTypes_kNoteOffEvent:
		// Note Off event - use helper to get the event data
		noteOff := C.getNoteOffEvent(event)
		c.processCtx.AddInputEvent(midi.NoteOffEvent{
			BaseEvent: midi.BaseEvent{
				EventChannel: uint8(noteOff.channel),
				Offset:       int32(event.sampleOffset),
			},
			NoteNumber: uint8(noteOff.pitch),
			Velocity:   uint8(noteOff.velocity * 127),
		})

		// Add more event types as needed
	}
}

// IEditController implementation
func (c *componentImpl) SetComponentState(state []byte) error {
	return nil
}

func (c *componentImpl) GetParameterCount() int32 {
	return c.processor.GetParameters().Count()
}

func (c *componentImpl) GetParameterInfo(index int32) (*vst3.ParameterInfo, error) {
	p := c.processor.GetParameters().GetByIndex(index)
	if p == nil {
		return nil, vst3.ErrInvalidArgument
	}

	return &vst3.ParameterInfo{
		ID:           p.ID,
		Title:        p.Name,
		ShortTitle:   p.ShortName,
		Units:        p.Unit,
		StepCount:    p.StepCount,
		DefaultValue: p.DefaultValue,
		UnitID:       p.UnitID,
		Flags:        int32(p.Flags),
	}, nil
}

func (c *componentImpl) GetParamStringByValue(id uint32, value float64) (string, error) {
	if p := c.processor.GetParameters().Get(id); p != nil {
		result := p.FormatValue(value)
		// fmt.Printf("Component.GetParamStringByValue: id=%d, value=%.3f -> '%s'\n", id, value, result)
		return result, nil
	}
	return "", vst3.ErrInvalidArgument
}

func (c *componentImpl) GetParamValueByString(id uint32, str string) (float64, error) {
	if p := c.processor.GetParameters().Get(id); p != nil {
		return p.ParseValue(str)
	}
	return 0, vst3.ErrInvalidArgument
}

func (c *componentImpl) NormalizedParamToPlain(id uint32, normalized float64) float64 {
	if p, ok := c.processor.GetParameters().GetOK(id); ok {
		return p.Denormalize(normalized)
	}
	return normalized
}

func (c *componentImpl) PlainParamToNormalized(id uint32, plain float64) float64 {
	if p, ok := c.processor.GetParameters().GetOK(id); ok {
		return p.Normalize(plain)
	}
	return plain
}

func (c *componentImpl) GetParamNormalized(id uint32) float64 {
	if value, ok := c.processor.GetParameters().GetNormalized(id); ok {
		return value
	}
	return 0
}

func (c *componentImpl) SetParamNormalized(id uint32, value float64) error {
	if p, ok := c.processor.GetParameters().GetOK(id); ok {
		// Debug parameter changes
		fmt.Printf("[PARAM_CHANGE] SetParamNormalized: id=%d, value=%.3f, plain=%.1f\n",
			id, value, p.Denormalize(value))
		p.SetNormalized(value)
		return nil
	}
	return vst3.ErrInvalidArgument
}

func (c *componentImpl) SetComponentHandler(handler interface{}) error {
	return nil
}

func (c *componentImpl) CreateView(name string) (interface{}, error) {
	if name != "" && name != "editor" && name != "web" {
		return nil, vst3.ErrNotImplemented
	}

	return c.EditorModel()
}

func (c *componentImpl) EditorModel() (*EditorModel, error) {
	params, err := c.parameterRegistry()
	if err != nil {
		return nil, err
	}

	return BuildEditorModel(c.pluginInfo, params)
}

func (c *componentImpl) EditorSnapshot() (*EditorSnapshot, error) {
	params, err := c.parameterRegistry()
	if err != nil {
		return nil, err
	}

	return BuildEditorSnapshot(c.pluginInfo, params)
}

func (c *componentImpl) ApplyEditorSnapshot(snapshot *EditorSnapshot) error {
	params, err := c.parameterRegistry()
	if err != nil {
		return err
	}

	return snapshot.Apply(params)
}

func (c *componentImpl) SetEditorParameter(id uint32, value float64) error {
	return c.SetParamNormalizedWithNotification(id, value)
}

// SetParamNormalizedWithNotification sets a parameter value and notifies the host
// This should be used when the plugin changes a parameter value internally
func (c *componentImpl) SetParamNormalizedWithNotification(id uint32, value float64) error {
	if p, ok := c.processor.GetParameters().GetOK(id); ok {
		// Notify host of parameter change
		if c.wrapper != nil {
			c.wrapper.notifyParamBeginEdit(id)
			p.SetNormalized(value)
			c.wrapper.notifyParamPerformEdit(id, value)
			c.wrapper.notifyEditorParameterChanged(id, p.GetNormalized(), p.GetPlain())
			c.wrapper.notifyParamEndEdit(id)
		} else {
			// Fallback if no wrapper available
			p.SetNormalized(value)
		}
		return nil
	}
	return vst3.ErrInvalidArgument
}

// processSampleAccurate processes audio with sample-accurate parameter automation
func (c *componentImpl) processSampleAccurate() {
	changes := c.processCtx.GetParameterChanges()
	numSamples := c.processCtx.NumSamples()
	lastOffset := 0

	// Store original buffers
	origInput := c.processCtx.Input
	origOutput := c.processCtx.Output

	// Process each chunk between parameter changes
	for _, change := range changes {
		if change.SampleOffset > lastOffset {
			// Process chunk from lastOffset to change.SampleOffset
			chunkSize := change.SampleOffset - lastOffset

			// Temporarily update context buffers to point to sub-slices (no allocation)
			c.processCtx.Input = nil
			c.processCtx.Output = nil

			// Set up input sub-slices
			for ch := 0; ch < len(origInput); ch++ {
				if lastOffset < len(origInput[ch]) {
					endOffset := lastOffset + chunkSize
					if endOffset > len(origInput[ch]) {
						endOffset = len(origInput[ch])
					}
					c.processCtx.Input = append(c.processCtx.Input, origInput[ch][lastOffset:endOffset])
				}
			}

			// Set up output sub-slices
			for ch := 0; ch < len(origOutput); ch++ {
				if lastOffset < len(origOutput[ch]) {
					endOffset := lastOffset + chunkSize
					if endOffset > len(origOutput[ch]) {
						endOffset = len(origOutput[ch])
					}
					c.processCtx.Output = append(c.processCtx.Output, origOutput[ch][lastOffset:endOffset])
				}
			}

			// Process this chunk
			c.processor.ProcessAudio(c.processCtx)

			lastOffset = change.SampleOffset
		}

		// Apply the parameter change
		c.processCtx.ApplyParameterChange(change)

		// Debug output
		if p, ok := c.processor.GetParameters().GetOK(change.ParamID); ok {
			fmt.Printf("[SAMPLE_ACCURATE] Applied param %d change at sample %d: value=%.6f, plain=%.1f\n",
				change.ParamID, change.SampleOffset, change.Value, p.GetPlain())
		}
	}

	// Process final chunk if there are samples remaining
	if lastOffset < numSamples {
		// Temporarily update context buffers for final chunk
		c.processCtx.Input = nil
		c.processCtx.Output = nil

		// Set up input sub-slices for final chunk
		for ch := 0; ch < len(origInput); ch++ {
			if lastOffset < len(origInput[ch]) {
				c.processCtx.Input = append(c.processCtx.Input, origInput[ch][lastOffset:])
			}
		}

		// Set up output sub-slices for final chunk
		for ch := 0; ch < len(origOutput); ch++ {
			if lastOffset < len(origOutput[ch]) {
				c.processCtx.Output = append(c.processCtx.Output, origOutput[ch][lastOffset:])
			}
		}

		// Process final chunk
		c.processor.ProcessAudio(c.processCtx)
	}

	// Restore original buffers
	c.processCtx.Input = origInput
	c.processCtx.Output = origOutput
}
