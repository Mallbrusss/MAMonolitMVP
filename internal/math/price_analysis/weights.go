package price_analysis

import (
	"errors"
	"fmt"
	"math"
)

type Signal struct {
	NormRSI          float64
	FractalDimension float64
	Peaks            []float64
	Valleys          []float64
	ShortSMA         float64
	LongSMA          float64
	TrendFactor      float64
	TotalSignal      float64
}

const tolerance = 0.02

func (p *PriceAnalysis) TotalSignal(prices []float64) (Signal, error) {
	if len(prices) == 0 {
		return Signal{}, errors.New("no prices provided")
	}

	fractalDimension := p.CalculateFractalDimension(prices)
	rsi, err := p.CalculateRSI(prices)
	if err != nil {
		return Signal{}, err
	}
	if len(rsi) == 0 {
		return Signal{}, fmt.Errorf("RSI calculation failed or returned no values")
	}

	smaShort, err := p.CalculateShortMovingAverage(prices)
	if err != nil {
		return Signal{}, err
	}
	if len(smaShort) == 0 {
		return Signal{}, fmt.Errorf("failed to calculate short SMA")
	}

	smaLong, err := p.CalculateLongMovingAverage(prices)
	if err != nil {
		return Signal{}, err
	}
	if len(smaLong) == 0 {
		return Signal{}, fmt.Errorf("failed to calculate long SMA")
	}

	normRSI := rsi[len(rsi)-1] / 100

	if smaLong[len(smaLong)-1] == 0 {
		return Signal{}, fmt.Errorf("long SMA is zero, cannot calculate trend factor")
	}
	trendFactor := (smaShort[len(smaShort)-1] - smaLong[len(smaLong)-1]) / smaLong[len(smaLong)-1]

	peaks, valleys := p.FindFractals(prices)

	currentPrice := prices[len(prices)-1]
	var fractalAdjustment float64 = 1.0

	for _, peak := range peaks {
		if peak != 0 {
			distance := math.Abs((currentPrice - peak) / peak)
			if distance <= tolerance {
				fractalAdjustment = 0.4 + 0.4*(distance/tolerance)
				break
			}
		}
	}

	for _, valley := range valleys {
		if valley != 0 {
			distance := math.Abs((currentPrice - valley) / valley)
			if distance <= tolerance {
				fractalAdjustment = 0.8 + 0.4*(distance/tolerance)
				break
			}
		}
	}

	totalSignal := (2 - fractalDimension) * (normRSI + trendFactor) * fractalAdjustment

	signal := Signal{
		NormRSI:          normRSI,
		FractalDimension: fractalDimension,
		Peaks:            peaks,
		Valleys:          valleys,
		ShortSMA:         smaShort[len(smaShort)-1],
		LongSMA:          smaLong[len(smaLong)-1],
		TrendFactor:      trendFactor,
		TotalSignal:      totalSignal,
	}

	fmt.Printf("RSI: %.4f, Fractal Dimension: %.4f, Peaks: %.4f, Vallyes: %.4f, Short SMA: %.4f, Long SMA: %.4f, Trend Factor: %.4f, Signal: %.4f\n",
		normRSI, fractalDimension, peaks, valleys, smaShort[len(smaShort)-1], smaLong[len(smaLong)-1], trendFactor, totalSignal)

	return signal, nil
}
