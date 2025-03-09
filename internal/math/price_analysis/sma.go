// Package price_analysis include RSI, SMA, Fractal Dimension
package price_analysis

func (p *PriceAnalysis) CalculateMovingAverage(period int, prices []float64) float64 {
	var sum float64

	for _, price := range prices {
		sum += price
	}
	movingAvg := sum / float64(period)
	return movingAvg
}
