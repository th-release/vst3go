package param

import (
	"sync"
	"testing"
)

func TestNormalizedAccessors(t *testing.T) {
	p := &Parameter{Min: -12, Max: 12}

	p.SetNormalized(0.75)
	if got := p.GetNormalized(); got != 0.75 {
		t.Fatalf("GetNormalized() = %f, want 0.75", got)
	}

	if got := p.GetValue(); got != 0.75 {
		t.Fatalf("GetValue() = %f, want 0.75", got)
	}
}

func TestPlainAccessors(t *testing.T) {
	p := &Parameter{Min: -12, Max: 12}

	p.SetPlain(6)
	if got := p.GetPlain(); got != 6 {
		t.Fatalf("GetPlain() = %f, want 6", got)
	}

	if got := p.GetPlainValue(); got != 6 {
		t.Fatalf("GetPlainValue() = %f, want 6", got)
	}
}

func TestParameterConcurrentAccess(t *testing.T) {
	p := &Parameter{Min: -12, Max: 12}

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(seed float64) {
			defer wg.Done()

			for j := 0; j < 1_000; j++ {
				value := seed + float64(j%10)/10.0
				p.SetNormalized(value / 10.0)
				_ = p.GetNormalized()
				p.SetPlain(value)
				_ = p.GetPlain()
			}
		}(float64(i))
	}

	wg.Wait()
}
