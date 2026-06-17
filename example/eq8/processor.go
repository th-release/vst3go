package eq8

import (
	"fmt"
	"math"

	"github.com/th-release/vst3go/pkg/framework/bus"
	frameworkparam "github.com/th-release/vst3go/pkg/framework/param"
	"github.com/th-release/vst3go/pkg/framework/process"
)

const (
	inputGainID  uint32 = 1
	outputGainID uint32 = 2
	bypassID     uint32 = 3
	analyzerID   uint32 = 4

	bandBaseID uint32 = 100
	bandStride uint32 = 10
	bandCount         = 8
)

var _ interface {
	Initialize(float64, int32) error
	ProcessAudio(*process.Context)
	GetParameters() *frameworkparam.Registry
	GetBuses() *bus.Configuration
	SetActive(bool) error
	GetLatencySamples() int32
	GetTailSamples() int32
} = (*Processor)(nil)

type bandParamIDs struct {
	enable    uint32
	frequency uint32
	gain      uint32
	q         uint32
	typ       uint32
}

type eqBand struct {
	enabled bool
	typ     float64
	freq    float64
	gain    float64
	q       float64
	left    biquad
	right   biquad
}

// Processor implements a complete downstream EQ8-style processor.
type Processor struct {
	params       *frameworkparam.AutoRegistry
	buses        *bus.Configuration
	sampleRate   float64
	maxBlockSize int32
	active       bool
	bands        [bandCount]eqBand
}

// NewProcessor creates the EQ processor and registers all parameters.
func NewProcessor() *Processor {
	params := frameworkparam.NewAutoRegistry()

	mustRegister(params, inputGainID, frameworkparam.GainParameter(inputGainID, "Input Gain").Build())
	mustRegister(params, outputGainID, frameworkparam.GainParameter(outputGainID, "Output Gain").Build())
	mustRegister(params, bypassID, frameworkparam.New(bypassID, "Bypass").Toggle().Bypass().Build())
	mustRegister(params, analyzerID, frameworkparam.New(analyzerID, "Analyzer").Toggle().Build())

	for bandIndex := 0; bandIndex < bandCount; bandIndex++ {
		registerBand(params, bandIndex+1)
	}

	return &Processor{
		params:       params,
		buses:        bus.Stereo(),
		sampleRate:   48000,
		maxBlockSize: 8192,
	}
}

func mustRegister(registry *frameworkparam.AutoRegistry, id uint32, param *frameworkparam.Parameter) {
	if err := registry.RegisterWithID(id, param); err != nil {
		panic(err)
	}
}

func registerBand(registry *frameworkparam.AutoRegistry, bandNumber int) {
	baseID := bandBaseID + uint32(bandNumber-1)*bandStride
	prefix := fmt.Sprintf("Band %d", bandNumber)

	options := []frameworkparam.ChoiceOption{
		{Value: 0, Name: "Bell"},
		{Value: 1, Name: "Low Shelf"},
		{Value: 2, Name: "High Shelf"},
		{Value: 3, Name: "Low Cut"},
		{Value: 4, Name: "High Cut"},
		{Value: 5, Name: "Notch"},
	}

	defaultFrequency := []float64{60, 120, 250, 500, 1000, 2000, 4000, 8000}
	if bandNumber < 1 || bandNumber > len(defaultFrequency) {
		defaultFrequency = append(defaultFrequency, 1000)
	}

	mustRegister(registry, baseID+0, frameworkparam.New(baseID+0, prefix+" Enable").Toggle().Build())
	mustRegister(registry, baseID+1, frameworkparam.FrequencyParameter(baseID+1, prefix+" Frequency", 20, 20000, defaultFrequency[bandNumber-1]).Build())
	mustRegister(registry, baseID+2, frameworkparam.GainParameter(baseID+2, prefix+" Gain").Default(0).Build())
	mustRegister(registry, baseID+3, frameworkparam.QParameter(baseID+3, prefix+" Q", 0.1, 12.0, 0.707).Build())
	mustRegister(registry, baseID+4, frameworkparam.Choice(baseID+4, prefix+" Type", options).Build())
}

// GetParameters returns the live parameter registry.
func (p *Processor) GetParameters() *frameworkparam.Registry {
	if p == nil || p.params == nil {
		return nil
	}
	return p.params.Registry
}

// GetBuses returns the stereo bus layout.
func (p *Processor) GetBuses() *bus.Configuration {
	if p == nil {
		return nil
	}
	return p.buses
}

// Initialize prepares the processor for host processing.
func (p *Processor) Initialize(sampleRate float64, maxBlockSize int32) error {
	if p == nil {
		return fmt.Errorf("eq8 processor is nil")
	}
	if sampleRate <= 0 {
		sampleRate = 48000
	}
	p.sampleRate = sampleRate
	p.maxBlockSize = maxBlockSize
	p.resetFilters()
	return nil
}

// SetActive resets the filter state when the host deactivates the plugin.
func (p *Processor) SetActive(active bool) error {
	p.active = active
	if !active {
		p.resetFilters()
	}
	return nil
}

// GetLatencySamples returns zero because this EQ is fully causal.
func (p *Processor) GetLatencySamples() int32 {
	return 0
}

// GetTailSamples returns zero because the EQ has no audible tail.
func (p *Processor) GetTailSamples() int32 {
	return 0
}

// ProcessAudio applies the current EQ state to the incoming block.
func (p *Processor) ProcessAudio(ctx *process.Context) {
	if p == nil || ctx == nil {
		return
	}

	if ctx.Param(bypassID) >= 0.5 {
		ctx.PassThrough()
		return
	}

	sampleRate := p.sampleRate
	if sampleRate <= 0 {
		sampleRate = ctx.SampleRate
	}
	if sampleRate <= 0 {
		sampleRate = 48000
	}

	p.updateBands(ctx, sampleRate)

	inputGain := dbToLinear(ctx.ParamPlain(inputGainID))
	outputGain := dbToLinear(ctx.ParamPlain(outputGainID))

	ctx.ProcessStereo(func(ch int, input, output []float32) {
		for sampleIndex := range input {
			value := float64(input[sampleIndex]) * inputGain
			for bandIndex := range p.bands {
				band := &p.bands[bandIndex]
				if !band.enabled {
					continue
				}
				if ch == 0 {
					value = band.left.Process(value)
				} else {
					value = band.right.Process(value)
				}
			}
			output[sampleIndex] = float32(value * outputGain)
		}
	})
}

func (p *Processor) resetFilters() {
	for bandIndex := range p.bands {
		p.bands[bandIndex].left.Reset()
		p.bands[bandIndex].right.Reset()
	}
}

func (p *Processor) updateBands(ctx *process.Context, sampleRate float64) {
	for bandIndex := range p.bands {
		params := bandParamsForBand(bandIndex)

		enabled := ctx.Param(params.enable) >= 0.5
		frequency := clamp(ctx.ParamPlain(params.frequency), 20, sampleRate*0.45)
		gain := ctx.ParamPlain(params.gain)
		q := clamp(ctx.ParamPlain(params.q), 0.1, 20)
		typ := ctx.ParamPlain(params.typ)

		band := &p.bands[bandIndex]
		band.enabled = enabled
		band.freq = frequency
		band.gain = gain
		band.q = q
		band.typ = typ

		coefficients := makeBandCoefficients(sampleRate, typ, frequency, gain, q)
		band.left.SetCoefficients(coefficients)
		band.right.SetCoefficients(coefficients)
	}
}

func bandParamsForBand(index int) bandParamIDs {
	baseID := bandBaseID + uint32(index)*bandStride
	return bandParamIDs{
		enable:    baseID + 0,
		frequency: baseID + 1,
		gain:      baseID + 2,
		q:         baseID + 3,
		typ:       baseID + 4,
	}
}

func clamp(value, minValue, maxValue float64) float64 {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

func dbToLinear(db float64) float64 {
	return math.Pow(10, db/20.0)
}
