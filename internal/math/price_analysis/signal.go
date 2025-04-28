package price_analysis

import (
	"errors"
	"fmt"
	"mamonolitmvp/internal/math/fractal_analysis"
	"math"
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
	NormalizeFdi
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

type NormalizeFdi struct {
	NormWidth     float64
	NormAsym      float64
	NormCurvature float64
	NormFdi       float64
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

	normWidth, normAsym, normCurvature, normFdi, err := p.fa.CalcNormalizedFdi(alpha, fAlpha, tau)
	if err != nil {
		return Signal{}, err
	}

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
		NormalizeFdi: NormalizeFdi{
			NormWidth:     normWidth,
			NormAsym:      normAsym,
			NormCurvature: normCurvature,
			NormFdi:       normFdi,
		},
	}

	return signal, nil
}

func (p *PriceAnalysis) SlidingWindowAnalysis(prices []float64, windowSize int) ([]Fdi, []NormalizeFdi, []float64, error) {
	if len(prices) < windowSize {
		return nil, nil, nil, errors.New("длина данных меньше размера окна")
	}

	var fdiSeries []Fdi
	var hurstSeries []float64

	for i := 0; i <= len(prices)-windowSize; i++ {
		window := prices[i : i+windowSize]

		signal, err := p.TotalSignal(window)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("ошибка в окне %d: %w", i, err)
		}

		fdiSeries = append(fdiSeries, signal.Fdi)
		hurstSeries = append(hurstSeries, signal.Hurst)
	}
	normFdi, err := p.CalculateWindoNormalFdi(fdiSeries)
	if err != nil {
		return nil, nil, nil, err
	}

	return fdiSeries, normFdi, hurstSeries, nil
}

func (fa *PriceAnalysis) CalculateWindoNormalFdi(fdiList []Fdi) ([]NormalizeFdi, error) {
	// Extract fields into separate slices
	widths := make([]float64, len(fdiList))
	asyms := make([]float64, len(fdiList))
	curvs := make([]float64, len(fdiList))

	for i, fdi := range fdiList {
		widths[i] = fdi.Width
		asyms[i] = fdi.Asym
		curvs[i] = fdi.Curvature
	}

	// Define weights
	w1, w2, w3 := 0.4, 0.3, 0.3

	// Normalize width, asymmetry, and curvature
	normWidth, err := normalizeValue(widths)
	if err != nil {
		return nil, err
	}
	normAsym, err := normalizeValue(asyms)
	if err != nil {
		return nil, err
	}
	normCurv, err := normalizeValue(curvs)
	if err != nil {
		return nil, err
	}

	// Ensure all slices have the same length
	if len(normWidth) != len(normAsym) || len(normWidth) != len(normCurv) {
		return nil, fmt.Errorf("input slices must have the same length")
	}

	// Calculate normalized FDI for each index
	result := make([]NormalizeFdi, len(fdiList))
	for i := range normWidth {
		fdi := NormalizeFdi{
			NormWidth:     normWidth[i],
			NormAsym:      normAsym[i],
			NormCurvature: normCurv[i],
			NormFdi:       w1*normWidth[i] + w2*math.Abs(normAsym[i]) + w3*normCurv[i],
		}
		result[i] = fdi
	}

	return result, nil
}
func normalizeValue(data []float64) ([]float64, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data slice is empty")
	}

	minD := data[0]
	maxD := data[0]

	// Find min and max values
	for _, value := range data {
		if value < minD {
			minD = value
		}
		if value > maxD {
			maxD = value
		}
	}

	// If all values are the same, return a slice of zeros
	if minD == maxD {
		normalized := make([]float64, len(data))
		for i := range normalized {
			normalized[i] = 0
		}
		return normalized, nil
	}

	// Normalize each value
	normalized := make([]float64, len(data))
	for i, value := range data {
		normalized[i] = (value - minD) / (maxD - minD)
	}

	return normalized, nil
}
