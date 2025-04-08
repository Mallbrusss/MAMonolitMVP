package fractal_analysis

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
	"math"
	"slices"
)

type FractalDimension struct {
}

func NewFractalDimension() *FractalDimension {
	return &FractalDimension{}
}

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

// Полиномиальная регрессия
func polyFit(x, y []float64, degree int) []float64 {
	n := len(x)
	X := mat.NewDense(n, degree+1, nil)

	// Формируем матрицу признаков (дизайн-матрицу)
	for i := 0; i < n; i++ {
		for j := 0; j <= degree; j++ {
			X.Set(i, j, math.Pow(x[i], float64(j)))
		}
	}

	// Вектор Y как матрица-столбец
	Y := mat.NewDense(n, 1, y)

	// Решаем уравнение X * beta = Y
	var beta mat.Dense
	err := beta.Solve(X, Y)
	if err != nil {
		panic(fmt.Sprintf("least squares solution failed: %v", err))
	}

	// Переносим результат в срез []float64
	result := make([]float64, degree+1)
	for i := range result {
		result[i] = beta.At(i, 0)
	}
	return result
}

func polyEval(x float64, coeffs []float64) float64 {
	y := 0.0
	for i, c := range coeffs {
		y += c * math.Pow(x, float64(i))
	}
	return y
}

func (fa *FractalDimension) DFA(series []float64, scaleMin, scaleMax, scaleStep, degree int) ([]float64, []float64, float64, float64) {
	n := len(series)
	mean := stat.Mean(series, nil)

	// Интеграция
	y := make([]float64, n)
	y[0] = series[0] - mean
	for i := 1; i < n; i++ {
		y[i] = y[i-1] + (series[i] - mean)
	}

	var logS, logF []float64
	for s := scaleMin; s <= scaleMax; s += scaleStep {
		numSegments := n / s
		var rmsList []float64

		for i := 0; i < numSegments; i++ {
			start := i * s
			end := start + s
			segment := y[start:end]

			x := make([]float64, s)
			for j := range x {
				x[j] = float64(j)
			}

			coeffs := polyFit(x, segment, degree)

			var rms float64
			for j := 0; j < s; j++ {
				trend := polyEval(x[j], coeffs)
				rms += math.Pow(segment[j]-trend, 2)
			}
			rms = math.Sqrt(rms / float64(s))
			rmsList = append(rmsList, rms)
		}

		Fs := stat.Mean(rmsList, nil)
		logS = append(logS, math.Log(float64(s)))
		logF = append(logF, math.Log(Fs))
	}

	// Линейная регрессия logF vs logS
	alpha, intercept := stat.LinearRegression(logS, logF, nil, false)

	// Коэффициент детерминации R²
	var ssTot, ssRes float64
	meanF := stat.Mean(logF, nil)
	for i := range logF {
		pred := alpha*logS[i] + intercept
		ssTot += math.Pow(logF[i]-meanF, 2)
		ssRes += math.Pow(logF[i]-pred, 2)
	}
	r2 := 1 - ssRes/ssTot

	return logS, logF, alpha, r2
}

// MFDFA MF-DFA: qList — список q (например, [-5,-4,...,5]), degree — порядок тренда (1, 2, ...)
func (fa *FractalDimension) MFDFA(series []float64, qList []float64, scaleMin, scaleMax, scaleStep, degree int) (map[float64][]float64, map[float64]float64, map[string]float64) {
	n := len(series)
	mean := stat.Mean(series, nil)

	// Интеграция
	y := make([]float64, n)
	y[0] = series[0] - mean
	for i := 1; i < n; i++ {
		y[i] = y[i-1] + (series[i] - mean)
	}

	hq := make(map[float64]float64)      // Херст для каждого q
	logFq := make(map[float64][]float64) // log F_q(s)
	logS := make([]float64, 0)

	for s := scaleMin; s <= scaleMax; s += scaleStep {
		numSegments := n / s
		if numSegments < 2 {
			continue
		}
		logS = append(logS, math.Log(float64(s)))

		// Вычисляем дисперсии по окнам
		flucts := make([]float64, 2*numSegments)

		for part := 0; part < 2; part++ {
			for i := 0; i < numSegments; i++ {
				start := i * s
				if part == 1 {
					start = n - (i+1)*s
				}
				end := start + s
				if end > n || start < 0 {
					continue
				}

				segment := y[start:end]
				x := make([]float64, s)
				for j := 0; j < s; j++ {
					x[j] = float64(j)
				}
				coeffs := polyFit(x, segment, degree)

				var rms float64
				for j := 0; j < s; j++ {
					trend := polyEval(x[j], coeffs)
					rms += math.Pow(segment[j]-trend, 2)
				}
				rms /= float64(s)
				flucts[i+part*numSegments] = rms
			}
		}

		// Считаем F_q(s) для каждого q
		for _, q := range qList {
			var Fq float64
			if q == 0 {
				// логарифмическое усреднение
				var sumLog float64
				for _, f2 := range flucts {
					if f2 > 0 {
						sumLog += math.Log(f2)
					}
				}
				Fq = math.Exp(0.5 * sumLog / float64(len(flucts)))
			} else {
				var sum float64
				for _, f2 := range flucts {
					sum += math.Pow(f2, q/2)
				}
				Fq = math.Pow(sum/float64(len(flucts)), 1.0/q)
			}
			logFq[q] = append(logFq[q], math.Log(Fq))
		}
	}

	// Линейная регрессия logFq[q] ~ logS => h(q)
	for _, q := range qList {
		h, _ := stat.LinearRegression(logS, logFq[q], nil, false)
		hq[q] = h
	}

	return logFq, hq, map[string]float64{"scale0": logS[0]}

}

func (fa *FractalDimension) CalcMultifractalSpectrum(hq map[float64]float64) (qList []float64, tauList, alphaList, fAlphaList []float64) {
	// Сортируем q по возрастанию
	qSorted := make([]float64, 0, len(hq))
	for q := range hq {
		qSorted = append(qSorted, q)
	}
	slices.Sort(qSorted)

	tau := make([]float64, len(qSorted))
	for i, q := range qSorted {
		tau[i] = q*hq[q] - 1
	}

	alpha := make([]float64, len(qSorted))
	fAlpha := make([]float64, len(qSorted))

	// Производная по q для alpha
	for i := 1; i < len(qSorted)-1; i++ {
		dq := qSorted[i+1] - qSorted[i-1]
		dTau := tau[i+1] - tau[i-1]
		alpha[i] = dTau / dq
		fAlpha[i] = qSorted[i]*alpha[i] - tau[i]
	}

	// Края — через одностороннюю производную
	alpha[0] = (tau[1] - tau[0]) / (qSorted[1] - qSorted[0])
	alpha[len(alpha)-1] = (tau[len(tau)-1] - tau[len(tau)-2]) / (qSorted[len(qSorted)-1] - qSorted[len(qSorted)-2])
	fAlpha[0] = qSorted[0]*alpha[0] - tau[0]
	fAlpha[len(fAlpha)-1] = qSorted[len(qSorted)-1]*alpha[len(alpha)-1] - tau[len(tau)-1]

	return qSorted, tau, alpha, fAlpha
}
