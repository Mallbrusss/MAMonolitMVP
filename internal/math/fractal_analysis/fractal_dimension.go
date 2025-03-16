package fractal_analysis

import (
	"fmt"
	"math"
)

type FractalDimension struct {
}

func NewFractalDimension() *FractalDimension {
	return &FractalDimension{}
}

// CalculateFractalDimension вычисляет фрактальную размерность методом FractalDimension
func (p *FractalDimension) CalculateFractalDimension(prices []float64) float64 {
	boxSize := p.determineNumBoxes(prices)
	if len(prices) == 0 || boxSize <= 0 {
		return 0
	}

	// Разбиваем данные на "ящики" заданного размера
	numBoxes := len(prices) / boxSize
	var totalRange float64 = 0

	for i := 0; i < numBoxes; i++ {
		minPrice := math.MaxFloat64
		maxPrice := -math.MaxFloat64

		for j := i * boxSize; j < (i+1)*boxSize && j < len(prices); j++ {
			if prices[j] < minPrice {
				minPrice = prices[j]
			}
			if prices[j] > maxPrice {
				maxPrice = prices[j]
			}
		}

		totalRange += maxPrice - minPrice
	}

	// Фрактальная размерность приближенно равна log(количество ящиков) / log(размер ящика)
	if totalRange == 0 {
		return 1 // Если все значения одинаковые, размерность равна 1
	}

	return 1 + math.Log(float64(numBoxes))/math.Log(float64(boxSize))
}

// FindFractals находит фракталы (пик и долина) в массиве цен
func (p *FractalDimension) FindFractals(prices []float64) ([]float64, []float64) {
	var peaks []float64
	var valleys []float64

	for i := 2; i < len(prices)-2; i++ {
		// Проверяем условие для пика
		if prices[i] > prices[i-1] && prices[i] > prices[i+1] &&
			prices[i] > prices[i-2] && prices[i] > prices[i+2] {
			peaks = append(peaks, prices[i])
		}

		// Проверяем условие для долины
		if prices[i] < prices[i-1] && prices[i] < prices[i+1] &&
			prices[i] < prices[i-2] && prices[i] < prices[i+2] {
			valleys = append(valleys, prices[i])
		}
	}

	return peaks, valleys
}

func (p *FractalDimension) AnalyzeFractals(prices []float64) string {
	peaks, valleys := p.FindFractals(prices)

	if len(peaks) > len(valleys) {
		return "Strong Uptrend"
	} else if len(peaks) < len(valleys) {
		return "Strong Downtrend"
	} else {
		return "Flat"
	}
}

func (p *FractalDimension) determineNumBoxes(prices []float64) int {
	if len(prices) < 5 {
		return 1
	}

	var percentage float64
	if len(prices) < 50 {
		percentage = 0.1
	} else {
		percentage = 0.4
	}

	numBoxes := int(float64(len(prices)) * percentage)

	if numBoxes < 5 {
		numBoxes = 5
	}

	fmt.Println("```````````Len:   ", len(prices))
	fmt.Println("!!!!!!!!!NUMBOXES:   ", numBoxes)
	fmt.Println("=========Last Price:   ", prices[len(prices)-1])

	return numBoxes
}
