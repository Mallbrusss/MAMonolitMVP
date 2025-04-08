package price_analysis

import (
	"errors"
	"fmt"
	"mamonolitmvp/internal/math/fractal_analysis"
)

type Signal struct {
	Hurst float64
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
	fa *fractal_analysis.FractalDimension
}

func NewPriceAnalysis() *PriceAnalysis {
	return &PriceAnalysis{
		fa: fractal_analysis.NewFractalDimension(),
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

	signal := Signal{
		Hurst: hurst,
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
