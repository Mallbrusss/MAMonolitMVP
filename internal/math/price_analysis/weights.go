package price_analysis

import (
	"errors"
	"fmt"
	"mamonolitmvp/internal/math/fractal_analysis"
)

type Signal struct {
	ShortSMA    []float64
	LongSMA     []float64
	TrendFactor float64
	Hurst       float64
	Mfdfa
	MfSpectrum
}

type Mfdfa struct {
	LogFq map[string][]float64
	Hq    map[string]float64
	LogS  map[string]float64
}

type MfSpectrum struct {
	Qsorted []float64
	Tau     []float64
	Alpha   []float64
	FAlpha  []float64
}

type PriceAnalysis struct {
	fa             *fractal_analysis.FractalDimension
	ShortSmaPeriod int
	LongSmaPeriod  int
}

func NewPriceAnalysis() *PriceAnalysis {
	return &PriceAnalysis{
		fa:             fractal_analysis.NewFractalDimension(),
		ShortSmaPeriod: 50,
		LongSmaPeriod:  100,
	}
}

func (p *PriceAnalysis) TotalSignal(prices []float64) (Signal, error) {
	if len(prices) == 0 {
		return Signal{}, errors.New("no prices provided")
	}

	hurst := p.fa.HurstExponent(prices)

	qList := []float64{-5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5} // Значения q
	scaleMin, scaleMax, scaleStep := 10, 100, 10             // Масштабы для анализа
	degree := 2
	logFq, hq, logS := p.fa.MFDFA(prices, qList, scaleMin, scaleMax, scaleStep, degree)

	qSorted, tau, alpha, fAlpha := p.fa.CalcMultifractalSpectrum(hq)

	stringLogFq := make(map[string][]float64)
	for k, v := range logFq {
		key := fmt.Sprintf("%v", k)
		stringLogFq[key] = v
	}
	hqString := make(map[string]float64)
	for k, v := range hq {
		key := fmt.Sprintf("%v", k)
		hqString[key] = v
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

	trendFactor := (smaShort[len(smaShort)-1] - smaLong[len(smaLong)-1]) / smaLong[len(smaLong)-1]
	fmt.Printf("Short SMA: %.4f, Long SMA: %.4f, Trend Factor: %.4f",
		smaShort[len(smaShort)-1], smaLong[len(smaLong)-1], trendFactor)

	signal := Signal{
		ShortSMA:    smaShort,
		LongSMA:     smaLong,
		TrendFactor: trendFactor,
		Hurst:       hurst,
		Mfdfa: Mfdfa{
			LogFq: stringLogFq,
			Hq:    hqString,
			LogS:  logS,
		},
		MfSpectrum: MfSpectrum{
			Qsorted: qSorted,
			Tau:     tau,
			Alpha:   alpha,
			FAlpha:  fAlpha,
		},
	}

	return signal, nil
}
