// Package price_analysis include RSI, SMA, Fractal Dimension
package price_analysis

import "math"

// CalculateFractalDimension вычисляет фрактальную размерность методом Бокса
func CalculateFractalDimension(prices []float64, boxSize int) float64 {
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
func FindFractals(prices []float64) ([]int, []int) {
	var peaks []int
	var valleys []int

	for i := 2; i < len(prices)-2; i++ {
		// Проверяем условие для пика
		if prices[i] > prices[i-1] && prices[i] > prices[i+1] &&
			prices[i] > prices[i-2] && prices[i] > prices[i+2] {
			peaks = append(peaks, i)
		}

		// Проверяем условие для долины
		if prices[i] < prices[i-1] && prices[i] < prices[i+1] &&
			prices[i] < prices[i-2] && prices[i] < prices[i+2] {
			valleys = append(valleys, i)
		}
	}

	return peaks, valleys
}
