package fractal_analysis

import (
	"fmt"
	"mamonolitmvp/internal/math/price_analysis"
	"testing"
)

func TestCalculateFractalDimension(t *testing.T) {
	pa := price_analysis.NewPriceAnalysis(10, 10, 10)
	prices := []float64{100, 102, 103, 50, 60, 40, 102, 68, 77, 106}
	fractalDim := pa.CalculateFractalDimension(prices)
	fmt.Printf("Фрактальная размерность: %.2f\n", fractalDim)
}

func TestFindFractals(t *testing.T) {
	pa := price_analysis.NewPriceAnalysis(10, 10, 10)

	prices := []float64{100, 102, 103, 101, 99, 100, 102, 105, 104, 106, 103, 101, 100, 102, 103, 101, 99, 100, 102, 105, 104, 106, 103, 101}

	peaks, valleys := pa.FindFractals(prices)

	fmt.Println("Пики:", peaks)
	fmt.Println("Долины:", valleys)
}
