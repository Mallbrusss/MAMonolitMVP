package fractal_analysis

import "math"

func (fa *FractalDimension) HurstExponent(timeSeries []float64) float64 {
	n := len(timeSeries)
	if n < 2 {
		return 0.5 // Стандартное значение для случайного процесса
	}

	// Шаг 1: Вычисляем среднее значение ряда
	mean := 0.0
	for _, value := range timeSeries {
		mean += value
	}
	mean /= float64(n)

	// Шаг 2: Строим отклонения от среднего (девиации)
	deviations := make([]float64, n)
	for i := 0; i < n; i++ {
		deviations[i] = timeSeries[i] - mean
	}

	// Шаг 3: Вычисляем накопленные суммы отклонений
	accumulated := make([]float64, n)
	accumulated[0] = deviations[0]
	for i := 1; i < n; i++ {
		accumulated[i] = accumulated[i-1] + deviations[i]
	}

	// Шаг 4: Находим диапазон (Range) и стандартное отклонение (Standard Deviation)
	rangeValue := math.Max(math.Abs(accumulated[0]), math.Abs(accumulated[n-1]))
	for i := 1; i < n; i++ {
		rangeValue = math.Max(rangeValue, math.Abs(accumulated[i]))
	}

	// Вычисляем стандартное отклонение временного ряда
	var sumSquares float64
	for _, value := range deviations {
		sumSquares += value * value
	}
	stdDev := math.Sqrt(sumSquares / float64(n))

	if stdDev == 0 {
		return 0.5 // Возвращаем значение 0.5 в случае отсутствия вариации (например, если все значения одинаковы)
	}
	// Шаг 5: Результат для Hurst Exponent
	rs := rangeValue / stdDev
	logN := math.Log(float64(n))
	logRS := math.Log(rs)

	// Коэффициент Hurst
	hurst := logRS / logN

	return hurst
}
