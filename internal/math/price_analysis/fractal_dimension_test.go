package price_analysis

import (
	"fmt"
	"testing"
)

func TestCalculateFractalDimension(t *testing.T) {
	prices := []float64{100, 102, 103, 101, 99, 100, 102, 105, 104, 106, 103, 101, 100, 102, 103, 101, 99, 100, 102, 105, 104, 106, 103, 101}
	boxSize := 6

	fractalDim := CalculateFractalDimension(prices, boxSize)
	fmt.Printf("Фрактальная размерность: %.2f\n", fractalDim)
}

func TestFindFractals(t *testing.T) {
	prices := []float64{100, 102, 103, 101, 99, 100, 102, 105, 104, 106, 103, 101, 100, 102, 103, 101, 99, 100, 102, 105, 104, 106, 103, 101}

	peaks, valleys := FindFractals(prices)

	fmt.Println("Пики:", peaks)
	fmt.Println("Долины:", valleys)
}
