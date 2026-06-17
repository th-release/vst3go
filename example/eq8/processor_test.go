package eq8

import (
	"math"
	"testing"

	"github.com/th-release/vst3go/pkg/framework/bus"
	"github.com/th-release/vst3go/pkg/framework/process"
)

func TestNewProcessorRegistersExpectedShape(t *testing.T) {
	processor := NewProcessor()

	if processor.GetParameters() == nil {
		t.Fatal("GetParameters() returned nil")
	}

	if got, want := processor.GetParameters().Count(), int32(44); got != want {
		t.Fatalf("parameter count = %d, want %d", got, want)
	}

	if got := processor.GetBuses().GetBusCount(bus.MediaTypeAudio, bus.DirectionInput); got != 1 {
		t.Fatalf("input bus count = %d, want 1", got)
	}
	if got := processor.GetBuses().GetBusCount(bus.MediaTypeAudio, bus.DirectionOutput); got != 1 {
		t.Fatalf("output bus count = %d, want 1", got)
	}
}

func TestProcessAudioPassThroughWhenBypassed(t *testing.T) {
	processor := NewProcessor()
	if err := processor.Initialize(48000, 64); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	registry := processor.GetParameters()
	registry.Get(bypassID).SetNormalized(1)

	ctx := process.NewContext(4, registry)
	ctx.Input = [][]float32{
		{0.25, -0.5, 0.75, -1.0},
		{0.10, 0.20, 0.30, 0.40},
	}
	ctx.Output = [][]float32{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}

	processor.ProcessAudio(ctx)

	for channelIndex := range ctx.Input {
		for sampleIndex := range ctx.Input[channelIndex] {
			if got, want := ctx.Output[channelIndex][sampleIndex], ctx.Input[channelIndex][sampleIndex]; got != want {
				t.Fatalf("output[%d][%d] = %v, want %v", channelIndex, sampleIndex, got, want)
			}
		}
	}
}

func TestProcessAudioAppliesInputGain(t *testing.T) {
	processor := NewProcessor()
	if err := processor.Initialize(48000, 64); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	registry := processor.GetParameters()
	registry.Get(bypassID).SetNormalized(0)
	registry.Get(inputGainID).SetPlain(6)
	registry.Get(outputGainID).SetPlain(0)

	for bandNumber := 1; bandNumber <= bandCount; bandNumber++ {
		baseID := bandBaseID + uint32(bandNumber-1)*bandStride
		registry.Get(baseID + 0).SetNormalized(0)
	}

	ctx := process.NewContext(4, registry)
	ctx.Input = [][]float32{{1, 1, 1, 1}, {1, 1, 1, 1}}
	ctx.Output = [][]float32{{0, 0, 0, 0}, {0, 0, 0, 0}}

	processor.ProcessAudio(ctx)

	want := float32(math.Pow(10, 6.0/20.0))
	for channelIndex := range ctx.Output {
		for sampleIndex, got := range ctx.Output[channelIndex] {
			if math.Abs(float64(got-want)) > 0.0005 {
				t.Fatalf("output[%d][%d] = %v, want %v", channelIndex, sampleIndex, got, want)
			}
		}
	}
}
