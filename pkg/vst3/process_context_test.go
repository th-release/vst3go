package vst3

import "testing"

func TestProcessContextFlags(t *testing.T) {
	ctx := &ProcessContext{
		State:              ProcessContextFlagPlaying | ProcessContextFlagTempoValid | ProcessContextFlagTimeSigValid | ProcessContextFlagBarPositionValid | ProcessContextFlagProjectTimeValid,
		SampleRate:         48000,
		ProjectTimeSamples:  1024,
		SystemTime:         2048,
		ProjectTimeMusic:   5.5,
		BarPositionMusic:   4.0,
		CycleStartMusic:    1.0,
		CycleEndMusic:      8.0,
		Tempo:              120,
		TimeSigNumerator:   4,
		TimeSigDenominator: 4,
		SamplesToNextClock: 240,
	}

	if !ctx.IsPlaying() {
		t.Fatal("IsPlaying() should report true")
	}
	if !ctx.HasTempo() || !ctx.HasTimeSignature() || !ctx.HasBarPosition() || !ctx.HasProjectTimeMusic() {
		t.Fatal("validity flags were not preserved")
	}
	if got := ctx.SampleRate; got != 48000 {
		t.Fatalf("SampleRate = %f, want 48000", got)
	}
	if got := ctx.ProjectTimeMusic; got != 5.5 {
		t.Fatalf("ProjectTimeMusic = %f, want 5.5", got)
	}
	if got := ctx.SamplesToNextClock; got != 240 {
		t.Fatalf("SamplesToNextClock = %d, want 240", got)
	}
}
