// Package price_analysis include RSI, SMA, Fractal Dimension
package price_analysis

import "fmt"

func (p *PriceAnalysis) CalculateShortMovingAverage(prices []float64) ([]float64, error) {
	if len(prices) < p.ShortSmaPeriod {
		return nil, fmt.Errorf("not enough data points for short SMA calculation: %d < %d", len(prices), p.ShortSmaPeriod)
	}

	smaValues := make([]float64, len(prices)-p.ShortSmaPeriod+1) // +1 здесь важно
	for i := p.ShortSmaPeriod - 1; i < len(prices); i++ {
		smaValues[i-p.ShortSmaPeriod+1] = sum(prices[i-p.ShortSmaPeriod+1:i+1]) / float64(p.ShortSmaPeriod)
	}
	return smaValues, nil
}

func (p *PriceAnalysis) CalculateLongMovingAverage(prices []float64) ([]float64, error) {
	if len(prices) < p.LongSmaPeriod {
		return nil, fmt.Errorf("not enough data points for long SMA calculation: %d < %d", len(prices), p.LongSmaPeriod)
	}

	smaValues := make([]float64, len(prices)-p.LongSmaPeriod+1) // +1 здесь важно
	for i := p.LongSmaPeriod - 1; i < len(prices); i++ {
		smaValues[i-p.LongSmaPeriod+1] = sum(prices[i-p.LongSmaPeriod+1:i+1]) / float64(p.LongSmaPeriod)
	}
	return smaValues, nil
}

func (p *PriceAnalysis) CalculateAmplitude(prices []float64) float64 {
	minPrice, maxPrice := prices[0], prices[0]

	for _, price := range prices {
		if price < minPrice {
			minPrice = price
		}
		if price > maxPrice {
			maxPrice = price
		}
	}

	return maxPrice - minPrice
}

func (p *PriceAnalysis) AnalyzeTrendWithSMA(prices []float64) string {
	shortSMA, err := p.CalculateShortMovingAverage(prices)
	if err != nil {
		return ""
	}
	longSMA, err := p.CalculateLongMovingAverage(prices)
	if err != nil {
		return ""
	}

	if shortSMA[len(shortSMA)-1] > longSMA[len(longSMA)-1] {
		return "Uptrend"
	} else if shortSMA[len(shortSMA)-1] < longSMA[len(longSMA)-1] {
		return "Downtrend"
	} else {
		return "Flat"
	}
}

func (p *PriceAnalysis) AnalyzeTrendWithAmplitude(prices []float64) string {
	amplitude := p.CalculateAmplitude(prices)

	if amplitude < 2.0 {
		return "Flat"
	} else {
		return "Trend"
	}
}
