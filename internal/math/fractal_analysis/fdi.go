package fractal_analysis

import (
	"math"
)

func calcWidth(alpha []float64) float64 {
	min, max := alpha[0], alpha[0]
	for _, a := range alpha {
		if a < min {
			min = a
		}
		if a > max {
			max = a
		}
	}
	return max - min
}

func calcAsymmetry(fAlpha []float64) float64 {
	n := len(fAlpha)
	if n < 3 {
		return 0 // Недостаточно данных для вычисления асимметрии
	}

	// Шаг 1: Вычисляем среднее значение
	sum := 0.0
	for _, value := range fAlpha {
		sum += value
	}
	mean := sum / float64(n)

	// Шаг 2: Вычисляем стандартное отклонение
	varianceSum := 0.0
	for _, value := range fAlpha {
		deviation := value - mean
		varianceSum += deviation * deviation
	}
	variance := varianceSum / float64(n-1) // Используем несмещенную оценку
	stdDev := math.Sqrt(variance)

	if stdDev == 0 {
		return 0 // Если стандартное отклонение равно нулю, асимметрия не определена
	}

	// Шаг 3: Вычисляем сумму кубов стандартизированных значений
	cubedSum := 0.0
	for _, value := range fAlpha {
		z := (value - mean) / stdDev
		cubedSum += math.Pow(z, 3)
	}

	// Шаг 4: Применяем формулу для асимметрии
	skewness := (float64(n) / (float64(n-1) * float64(n-2))) * cubedSum
	return skewness
}

// Расчет кривизны через приближение второй производной
func calcCurvature(tau []float64) float64 {
	n := len(tau)
	if n < 3 {
		return 0 // Недостаточно точек для вычисления кривизны
	}

	totalCurvature := 0.0
	count := 0

	for i := 1; i < n-1; i++ {
		// Вычисляем первую производную
		firstDerivative := (tau[i+1] - tau[i-1]) / 2.0

		// Вычисляем вторую производную
		secondDerivative := tau[i+1] - 2*tau[i] + tau[i-1]

		// Вычисляем кривизну
		denominator := math.Pow(1+firstDerivative*firstDerivative, 1.5)
		if denominator == 0 {
			continue // Избегаем деления на ноль
		}
		curvature := math.Abs(secondDerivative) / denominator

		// Суммируем кривизну
		totalCurvature += curvature
		count++
	}

	// Возвращаем среднюю кривизну
	if count == 0 {
		return 0
	}
	return totalCurvature / float64(count)
}

func (fa *FractalDimension) CalcFdi(alpha, fAlpha, tau []float64) (float64, float64, float64, float64) {
	width := calcWidth(alpha)
	asym := calcAsymmetry(fAlpha)
	curvature := calcCurvature(tau)

	w1, w2, w3 := 0.4, 0.3, 0.3

	fdi := w1*width + w2*math.Abs(asym) + w3*curvature

	return width, asym, curvature, fdi
}

func (fa *FractalDimension) CalcNormalizedFdi(alpha, fAlpha, tau []float64) (float64, float64, float64, float64, error) {
	// Шаг 1: Вычисляем raw значения width, asymmetry и curvature
	width := calcWidth(alpha)
	asym := calcAsymmetry(fAlpha)
	curvature := calcCurvature(tau)

	// Шаг 2: Создаем массивы для нормализации
	values := []struct {
		value float64
	}{
		{value: width},
		{value: asym},
		{value: curvature},
	}

	// Шаг 3: Нормализуем каждое значение
	normValues := make([]float64, len(values))
	minD := values[0].value
	maxD := values[0].value

	for _, val := range values {
		if val.value < minD {
			minD = val.value
		}
		if val.value > maxD {
			maxD = val.value
		}
	}

	// Проверяем на случай одинаковых значений
	if minD == maxD {
		for i := range normValues {
			normValues[i] = 0
		}
	} else {
		for i, val := range values {
			normValues[i] = (val.value - minD) / (maxD - minD)
		}
	}

	// Шаг 4: Извлекаем нормализованные значения
	normWidth := normValues[0]
	normAsym := normValues[1]
	normCurvature := normValues[2]

	// Шаг 5: Вычисляем FDI с использованием нормализованных значений
	w1, w2, w3 := 0.4, 0.3, 0.3
	normFdi := w1*normWidth + w2*math.Abs(normAsym) + w3*normCurvature

	return normWidth, normAsym, normCurvature, normFdi, nil
}
