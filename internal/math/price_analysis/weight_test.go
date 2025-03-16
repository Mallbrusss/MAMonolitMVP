package price_analysis

import (
	"fmt"
	"testing"
)

func TestPriceAnalysis_TotalSignal(t *testing.T) {
	pricesStrongUpTrend := []float64{100, 102, 105, 110, 115, 120, 125, 130, 135, 140}
	pricesLittleUpTrend := []float64{100, 101, 102, 103, 104, 105, 106, 107, 108, 109}
	pricesFlatTrend := []float64{100, 101, 99, 100, 101, 100, 99, 100, 101, 100}
	pricesLittleDownTred := []float64{100, 99, 98, 97, 96, 95, 94, 93, 92, 91}
	pricesStrongDownTrend := []float64{100, 95, 90, 85, 60, 55, 15, 10, 8, 40}

	pa := NewPriceAnalysis(3, 7, 3)
	totalSigStrongUpTrend, err := pa.TotalSignal(pricesStrongUpTrend)
	if err != nil {
		t.Fatal(err)
	}
	totalLittleUpTrend, err := pa.TotalSignal(pricesLittleUpTrend)
	if err != nil {
		t.Fatal(err)
	}
	totalFlatTrend, err := pa.TotalSignal(pricesFlatTrend)
	if err != nil {
		t.Fatal(err)
	}
	totalLittleDownTrend, err := pa.TotalSignal(pricesLittleDownTred)
	if err != nil {
		t.Fatal(err)
	}
	totalStrongDownTrend, err := pa.TotalSignal(pricesStrongDownTrend)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Print("Strong Up Trend: ", totalSigStrongUpTrend,
		"\n Little Up Trend: ", totalLittleUpTrend,
		"\n Flat Trend: ", totalFlatTrend,
		"\nLittle Down Trend: ", totalLittleDownTrend,
		"\n Strong Down Trend: ", totalStrongDownTrend)
}
