package price_analysis

import (
	"fmt"
	"mamonolitmvp/internal/math/fractal_analysis"
)

type PriceAnalysis struct {
	ShortSmaPeriod int
	LongSmaPeriod  int
	RSIPeriod      int
	BoxSize        int
	fa             *fractal_analysis.FractalDimension
}

func NewPriceAnalysis(shortSmaPeriod, longSmaPeriod, rsiPeriod int) *PriceAnalysis {
	return &PriceAnalysis{
		ShortSmaPeriod: shortSmaPeriod,
		LongSmaPeriod:  longSmaPeriod,
		RSIPeriod:      rsiPeriod,
		fa:             fractal_analysis.NewFractalDimension(),
	}
}

func (p *PriceAnalysis) CalculateRSI(prices []float64) ([]float64, error) {
	if len(prices) == 0 {
		return nil, fmt.Errorf("price data is empty")
	}
	if len(prices) < p.RSIPeriod+1 {
		return nil, fmt.Errorf("not enough data to calculate RSI with period %d", p.RSIPeriod)
	}

	gains := make([]float64, len(prices)-1)
	losses := make([]float64, len(prices)-1)

	for i := 1; i < len(prices); i++ {
		diff := prices[i] - prices[i-1]
		if diff > 0 {
			gains[i-1] = diff
		} else {
			losses[i-1] = -diff
		}
	}

	rsiValues := make([]float64, len(prices)-p.RSIPeriod)
	avgGain := sum(gains[:p.RSIPeriod]) / float64(p.RSIPeriod)
	avgLoss := sum(losses[:p.RSIPeriod]) / float64(p.RSIPeriod)

	if avgLoss == 0 {
		rsiValues[0] = 100.0
	} else {
		rs := avgGain / avgLoss
		rsiValues[0] = 100 - (100 / (1 + rs))
	}

	for i := p.RSIPeriod; i < len(prices)-1; i++ {
		avgGain = (avgGain*(float64(p.RSIPeriod-1)) + gains[i-1]) / float64(p.RSIPeriod)
		avgLoss = (avgLoss*(float64(p.RSIPeriod-1)) + losses[i-1]) / float64(p.RSIPeriod)

		if avgLoss == 0 {
			rsiValues[i-p.RSIPeriod+1] = 100.0
		} else {
			rs := avgGain / avgLoss
			rsiValues[i-p.RSIPeriod+1] = 100 - (100 / (1 + rs))
		}
	}

	return rsiValues, nil
}

func sum(values []float64) float64 {
	total := 0.0
	for _, v := range values {
		total += v
	}
	return total
}

func (p *PriceAnalysis) Average(numbers []float64) float64 {
	sum := 0.0
	for _, number := range numbers {
		sum += number
	}
	return sum / float64(len(numbers))
}

func (p *PriceAnalysis) AnalyzeTrendWithRSI(prices []float64) string {
	rsiValues, err := p.CalculateRSI(prices)
	if err != nil {
		return fmt.Sprintf("Error calculating RSI: %v", err)
	}

	latestRSI := rsiValues[len(rsiValues)-1]

	if latestRSI > 70 {
		return "Strong Uptrend"
	} else if latestRSI < 30 {
		return "Strong Downtrend"
	} else if latestRSI > 40 && latestRSI < 60 {
		return "Flat Market"
	}

	return "Neutral"
}
