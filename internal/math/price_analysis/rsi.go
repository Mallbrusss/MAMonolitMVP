package price_analysis

import (
	"fmt"
	"math"
)

type PriceAnalysis struct{}

func NewRsi() *PriceAnalysis {
	return &PriceAnalysis{}
}

func (p *PriceAnalysis) CalculateRSI(prices []float64, period int) (float64, error) {
	if len(prices) == 0 {
		return math.NaN(), fmt.Errorf("price data is empty")
	}
	if len(prices) < period+1 {
		return math.NaN(), fmt.Errorf("not enough data to calculate RSI with period %d", period)
	}

	changes := make([]float64, len(prices)-1)
	for i := range changes {
		changes[i] = prices[i+1] - prices[i]
	}

	gains := make([]float64, len(changes))
	losses := make([]float64, len(changes))
	for i, change := range changes {
		if change >= 0 {
			gains[i] = change
		} else {
			losses[i] = math.Abs(change)
		}
	}

	avgGain := p.Average(gains[:period])
	avgLoss := p.Average(losses[:period])
	rs := avgGain / avgLoss
	var rsi float64
	if avgLoss == 0 {
		if avgGain == 0 {
			rsi = 50
		} else {
			rsi = 100
		}
	} else {
		rsi = 100 - (100 / (1 + rs))
	}
	return rsi, nil
}

func (p *PriceAnalysis) Average(numbers []float64) float64 {
	sum := 0.0
	for _, number := range numbers {
		sum += number
	}
	return sum / float64(len(numbers))
}
