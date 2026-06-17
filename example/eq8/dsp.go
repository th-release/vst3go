package eq8

import "math"

type biquadCoefficients struct {
	a0 float64
	b0 float64
	b1 float64
	b2 float64
	a1 float64
	a2 float64
}

type biquad struct {
	b0 float64
	b1 float64
	b2 float64
	a1 float64
	a2 float64
	z1 float64
	z2 float64
}

func (b *biquad) Reset() {
	b.z1 = 0
	b.z2 = 0
}

func (b *biquad) SetCoefficients(coefficients biquadCoefficients) {
	if coefficients.a0 == 0 {
		coefficients = identityCoefficients()
	}
	invA0 := 1 / coefficients.a0
	b.b0 = coefficients.b0 * invA0
	b.b1 = coefficients.b1 * invA0
	b.b2 = coefficients.b2 * invA0
	b.a1 = coefficients.a1 * invA0
	b.a2 = coefficients.a2 * invA0
}

func (b *biquad) Process(input float64) float64 {
	output := b.b0*input + b.z1
	b.z1 = b.b1*input - b.a1*output + b.z2
	b.z2 = b.b2*input - b.a2*output
	return output
}

func identityCoefficients() biquadCoefficients {
	return biquadCoefficients{
		a0: 1,
		b0: 1,
		b1: 0,
		b2: 0,
		a1: 0,
		a2: 0,
	}
}

func makeBandCoefficients(sampleRate float64, bandType float64, frequency, gainDB, q float64) biquadCoefficients {
	if sampleRate <= 0 {
		sampleRate = 48000
	}
	if frequency < 1 {
		frequency = 1
	}
	if q < 0.1 {
		q = 0.1
	}
	if frequency > sampleRate*0.49 {
		frequency = sampleRate * 0.49
	}

	w0 := 2 * math.Pi * frequency / sampleRate
	sinW0 := math.Sin(w0)
	cosW0 := math.Cos(w0)
	alpha := sinW0 / (2 * q)
	A := math.Pow(10, gainDB/40.0)

	switch int(math.Round(bandType)) {
	case 1: // Low Shelf
		sqrtA := math.Sqrt(A)
		beta := 2 * sqrtA * alpha
		return biquadCoefficients{
			a0: (A + 1) + (A-1)*cosW0 + beta,
			b0: A * ((A + 1) - (A-1)*cosW0 + beta),
			b1: 2 * A * ((A - 1) - (A+1)*cosW0),
			b2: A * ((A + 1) - (A-1)*cosW0 - beta),
			a1: -2 * ((A - 1) + (A+1)*cosW0),
			a2: (A + 1) + (A-1)*cosW0 - beta,
		}
	case 2: // High Shelf
		sqrtA := math.Sqrt(A)
		beta := 2 * sqrtA * alpha
		return biquadCoefficients{
			a0: (A + 1) - (A-1)*cosW0 + beta,
			b0: A * ((A + 1) + (A-1)*cosW0 + beta),
			b1: -2 * A * ((A - 1) + (A+1)*cosW0),
			b2: A * ((A + 1) + (A-1)*cosW0 - beta),
			a1: 2 * ((A - 1) - (A+1)*cosW0),
			a2: (A + 1) - (A-1)*cosW0 - beta,
		}
	case 3: // Low Cut
		return biquadCoefficients{
			a0: 1 + alpha,
			b0: (1 + cosW0) / 2,
			b1: -(1 + cosW0),
			b2: (1 + cosW0) / 2,
			a1: -2 * cosW0,
			a2: 1 - alpha,
		}
	case 4: // High Cut
		return biquadCoefficients{
			a0: 1 + alpha,
			b0: (1 - cosW0) / 2,
			b1: 1 - cosW0,
			b2: (1 - cosW0) / 2,
			a1: -2 * cosW0,
			a2: 1 - alpha,
		}
	case 5: // Notch
		return biquadCoefficients{
			a0: 1 + alpha,
			b0: 1,
			b1: -2 * cosW0,
			b2: 1,
			a1: -2 * cosW0,
			a2: 1 - alpha,
		}
	default: // Bell
		return biquadCoefficients{
			a0: 1 + alpha/A,
			b0: 1 + alpha*A,
			b1: -2 * cosW0,
			b2: 1 - alpha*A,
			a1: -2 * cosW0,
			a2: 1 - alpha/A,
		}
	}
}
