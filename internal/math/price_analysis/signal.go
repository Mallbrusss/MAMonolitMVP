package price_analysis

import (
	"errors"
	"fmt"
	"mamonolitmvp/internal/math/fractal_analysis"
)

type Signal struct {
	Total       int
	ShortSMA    []float64
	LongSMA     []float64
	TrendFactor float64
	Hurst       float64
	Mfdfa
	MfSpectrum
	Fdi
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

type SlidingWindow struct {
	FdiSeries   []Fdi
	HurstSeries []float64
}

func NewPriceAnalysis() *PriceAnalysis {
	return &PriceAnalysis{
		fa:             fractal_analysis.NewFractalDimension(),
		ShortSmaPeriod: 50,
		LongSmaPeriod:  100,
	}
}

type Fdi struct {
	Width     float64
	Asym      float64
	Curvature float64
	Fdi       float64
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

	width, asym, curvature, fdi := p.fa.CalcFdi(alpha, fAlpha, tau)

	signal := Signal{
		Total:       len(prices),
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
		Fdi: Fdi{
			Width:     width,
			Asym:      asym,
			Curvature: curvature,
			Fdi:       fdi,
		},
	}

	return signal, nil
}

func (p *PriceAnalysis) SlidingWindowAnalysis(prices []float64, windowSize int) ([]Fdi, []float64, error) {
	if len(prices) < windowSize {
		return nil, nil, errors.New("длина данных меньше размера окна")
	}

	var fdiSeries []Fdi
	var hurstSeries []float64

	for i := 0; i <= len(prices)-windowSize; i++ {
		window := prices[i : i+windowSize]

		signal, err := p.TotalSignal(window)
		if err != nil {
			return nil, nil, fmt.Errorf("ошибка в окне %d: %w", i, err)
		}

		fdiSeries = append(fdiSeries, signal.Fdi)
		hurstSeries = append(hurstSeries, signal.Hurst)

	}

	return fdiSeries, hurstSeries, nil
}
