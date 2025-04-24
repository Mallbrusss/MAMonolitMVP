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

func calcAsymmetry(alpha, fAlpha []float64) float64 {
	if len(alpha) == 0 || len(fAlpha) == 0 || len(alpha) != len(fAlpha) {
		return 0 // Если массивы пустые или имеют разную длину, возвращаем 0
	}

	// Находим индекс максимума в fAlpha
	maxF := fAlpha[0]
	modeIndex := 0
	for i, f := range fAlpha {
		if f > maxF {
			maxF = f
			modeIndex = i
		}
	}

	// Находим минимальное и максимальное значения в alpha
	alphaMin, alphaMax := alpha[0], alpha[0]
	for _, a := range alpha {
		if a < alphaMin {
			alphaMin = a
		}
		if a > alphaMax {
			alphaMax = a
		}
	}

	// Вычисляем центральный угол и ширину
	alphaCenter := (alphaMax + alphaMin) / 2
	alphaMode := alpha[modeIndex]
	width := alphaMax - alphaMin

	// Проверяем, что ширина не равна 0
	if width == 0 {
		return 0
	}

	// Возвращаем нормированную асимметрию
	return (alphaMode - alphaCenter) / width
}

// Расчет кривизны через приближение второй производной
func calcCurvature(alpha, fAlpha []float64) float64 {
	// Находим индекс максимального значения в массиве fAlpha
	modeIndex := 0
	maxF := fAlpha[0]
	for i, f := range fAlpha {
		if f > maxF {
			maxF = f
			modeIndex = i
		}
	}

	// Проверяем, что modeIndex не находится на границах массива
	if modeIndex <= 0 || modeIndex >= len(fAlpha)-1 {
		return 0 // Недостаточно данных для расчета второй производной
	}

	// Вычисляем шаг дискретизации alpha
	deltaAlpha := alpha[1] - alpha[0]

	// Приближенное значение второй производной
	secondDerivative := (fAlpha[modeIndex-1] - 2*fAlpha[modeIndex] + fAlpha[modeIndex+1]) / (deltaAlpha * deltaAlpha)

	return secondDerivative
}

func (fa *FractalDimension) CalcFdi(alpha, fAlpha []float64) (float64, float64, float64, float64) {
	width := calcWidth(alpha)
	asym := calcAsymmetry(alpha, fAlpha)
	curvature := calcCurvature(alpha, fAlpha)

	w1, w2, w3 := 0.4, 0.3, 0.3

	fdi := w1*width + w2*math.Abs(asym) + w3*curvature

	return width, asym, curvature, fdi
}
