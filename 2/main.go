package main

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"math"
	"math/rand"
	"slices"
)

// Задание 5:

func GetFreqDistr(RParamsArr []float64, A, B float64, IntervalsCount int) []float64 {
	var dY float64 = (B - A) / float64(IntervalsCount)
	Freq := make([]float64, IntervalsCount)

	var fN int
	for i := 0; i < len(RParamsArr); i++ {
		fN = int(math.Floor(RParamsArr[i] / dY))
		Freq[fN]++
	}

	for i := 0; i < IntervalsCount; i++ {
		Freq[i] /= float64(len(RParamsArr)) * dY
	}

	return Freq
}

// Задание 4: функция RANDPeriod

func RANDPeriod(X []float64) []int64 {
	var n int64 = int64(len(X))

	var element float64
	for i := int64(0); i < n; i++ {
		element = X[i]
		for j := i + 1; j < n; j++ {
			if element == X[j] {
				return []int64{j - i, i, j}
			}
		}
	}
	return []int64{-1, -1, -1}
}

// Задание 1: функция RAND

func RAND(a, b, m, x_i int64) int64 {
	return (a*x_i + b) % m
}

func main() {

	// Задание 2: расчет последовательностей случайных чисел

	var a, b, m int64 = 22695477, 1, 1 << 32
	var x0 int64 = 1

	var A, B float64 = 0, 10
	var N int = 100

	RNumsArr := make([]float64, N)
	RNumsArr[0] = float64(RAND(a, b, m, x0))
	for i := 1; i < N; i++ {
		RNumsArr[i] = float64(RAND(a, b, m, int64(RNumsArr[i-1])))
	}
	for i := 0; i < N; i++ {
		RNumsArr[i] = RNumsArr[i] / float64(m)
	}

	RParamsArr := make([]float64, N)
	for i := 0; i < N; i++ {
		RParamsArr[i] = A + (B-A)*RNumsArr[i]
	}

	fmt.Println("Максимальное значение при N=100:", slices.Max(RParamsArr))
	fmt.Println("Минимальное значение при N=100:", slices.Min(RParamsArr), "\n")

	N = 1000

	RNumsArr = make([]float64, N)
	RNumsArr[0] = float64(RAND(a, b, m, x0))
	for i := 1; i < N; i++ {
		RNumsArr[i] = float64(RAND(a, b, m, int64(RNumsArr[i-1])))
	}
	for i := 0; i < N; i++ {
		RNumsArr[i] = RNumsArr[i] / float64(m)
	}

	RParamsArr_e3 := make([]float64, N)
	for i := 0; i < N; i++ {
		RParamsArr_e3[i] = A + (B-A)*RNumsArr[i]
	}

	fmt.Println("Максимальное значение при N=1000:", slices.Max(RParamsArr_e3))
	fmt.Println("Минимальное значение при N=1000:", slices.Min(RParamsArr_e3), "\n")

	N = 10000

	RNumsArr = make([]float64, N)
	RNumsArr[0] = float64(RAND(a, b, m, x0))
	for i := 1; i < N; i++ {
		RNumsArr[i] = float64(RAND(a, b, m, int64(RNumsArr[i-1])))
	}
	for i := 0; i < N; i++ {
		RNumsArr[i] = RNumsArr[i] / float64(m)
	}

	RParamsArr_e4 := make([]float64, N)
	for i := 0; i < N; i++ {
		RParamsArr_e4[i] = A + (B-A)*RNumsArr[i]
	}

	fmt.Println("Максимальное значение при N=10000:", slices.Max(RParamsArr_e4))
	fmt.Println("Минимальное значение при N=10000:", slices.Min(RParamsArr_e4), "\n")

	N = 100000

	RNumsArr = make([]float64, N)
	RNumsArr[0] = float64(RAND(a, b, m, x0))
	for i := 1; i < N; i++ {
		RNumsArr[i] = float64(RAND(a, b, m, int64(RNumsArr[i-1])))
	}
	for i := 0; i < N; i++ {
		RNumsArr[i] = RNumsArr[i] / float64(m)
	}

	RParamsArr_e5 := make([]float64, N)
	for i := 0; i < N; i++ {
		RParamsArr_e5[i] = A + (B-A)*RNumsArr[i]
	}

	fmt.Println("Максимальное значение при N=100000:", slices.Max(RParamsArr_e5))
	fmt.Println("Минимальное значение при N=100000:", slices.Min(RParamsArr_e5), "\n")

	// Задание 3

	var M float64 = (A + B) / 2
	fmt.Println("Теоретическое мат. ожидание:", M)
	var D float64 = ((B - A) * (B - A)) / 12
	fmt.Println("Теоретическая дисперсия:", D, "\n")

	N = 100

	var M_e2 float64
	for i := 0; i < N; i++ {
		M_e2 += RParamsArr[i]
	}
	M_e2 /= float64(N)
	fmt.Println("Мат. ожидание при N=100:", M_e2)

	var EpsM1 float64 = math.Abs((M-M_e2)/M) * 100
	fmt.Println("Погрешность мат. ожидания при N=100:", EpsM1)

	var D_e2 float64
	for i := 0; i < N; i++ {
		D_e2 += RParamsArr[i] * RParamsArr[i]
	}
	D_e2 /= float64(N)
	D_e2 -= M_e2 * M_e2
	D_e2 *= float64(N) / (float64(N) - 1)
	fmt.Println("Дисперсия при N=100:", D_e2)

	var EpsD1 float64 = math.Abs((D-D_e2)/D) * 100
	fmt.Println("Погрешность дисперсии при N=100:", EpsD1, "\n")

	N = 1000

	var M_e3 float64
	for i := 0; i < N; i++ {
		M_e3 += RParamsArr_e3[i]
	}
	M_e3 /= float64(N)
	fmt.Println("Мат. ожидание при N=1000:", M_e3)

	var EpsM2 float64 = math.Abs((M-M_e3)/M) * 100
	fmt.Println("Погрешность мат. ожидания при N=1000:", EpsM2)

	var D_e3 float64
	for i := 0; i < N; i++ {
		D_e3 += RParamsArr_e3[i] * RParamsArr_e3[i]
	}
	D_e3 /= float64(N)
	D_e3 -= M_e3 * M_e3
	D_e3 *= float64(N) / (float64(N) - 1)
	fmt.Println("Дисперсия при N=1000:", D_e3)

	var EpsD2 float64 = math.Abs((D-D_e3)/D) * 100
	fmt.Println("Погрешность дисперсии при N=1000:", EpsD2, "\n")

	N = 10000

	var M_e4 float64
	for i := 0; i < N; i++ {
		M_e4 += RParamsArr_e4[i]
	}
	M_e4 /= float64(N)
	fmt.Println("Мат. ожидание при N=10000:", M_e4)

	var EpsM3 float64 = math.Abs((M-M_e4)/M) * 100
	fmt.Println("Погрешность мат. ожидания при N=10000:", EpsM3)

	var D_e4 float64
	for i := 0; i < N; i++ {
		D_e4 += RParamsArr_e4[i] * RParamsArr_e4[i]
	}
	D_e4 /= float64(N)
	D_e4 -= M_e4 * M_e4
	D_e4 *= float64(N) / (float64(N) - 1)
	fmt.Println("Дисперсия при N=10000:", D_e4)

	var EpsD3 float64 = math.Abs((D-D_e4)/D) * 100
	fmt.Println("Погрешность дисперсии при N=10000:", EpsD3, "\n")

	N = 100000

	var M_e5 float64
	for i := 0; i < N; i++ {
		M_e5 += RParamsArr_e5[i]
	}
	M_e5 /= float64(N)
	fmt.Println("Мат. ожидание при N=100000:", M_e5)

	var EpsM4 float64 = math.Abs((M-M_e5)/M) * 100
	fmt.Println("Погрешность мат. ожидания при N=100000:", EpsM4)

	var D_e5 float64
	for i := 0; i < N; i++ {
		D_e5 += RParamsArr_e5[i] * RParamsArr_e5[i]
	}
	D_e5 /= float64(N)
	D_e5 -= M_e5 * M_e5
	D_e5 *= float64(N) / (float64(N) - 1)
	fmt.Println("Дисперсия при N=100000:", D_e5)

	var EpsD4 float64 = math.Abs((D-D_e5)/D) * 100
	fmt.Println("Погрешность дисперсии при N=100000:", EpsD4, "\n")

	// Задание 4

	var TEST_1 = RANDPeriod(RParamsArr)
	fmt.Println("Результаты теста на периодичность последовательности при N=100:", TEST_1)
	var TEST_2 = RANDPeriod(RParamsArr_e3)
	fmt.Println("Результаты теста на периодичность последовательности при N=1000:", TEST_2)
	var TEST_3 = RANDPeriod(RParamsArr_e4)
	fmt.Println("Результаты теста на периодичность последовательности при N=10000:", TEST_3)
	var TEST_4 = RANDPeriod(RParamsArr_e5)
	fmt.Println("Результаты теста на периодичность последовательности при N=100000:", TEST_4, "\n")

	// Задание 6

	var K int = 10
	resX := make([]float64, K)
	for k := 0; k < K; k++ {
		resX[k] = ((B - A) / float64(K)) * (0.5 + float64(k))
	}
	fmt.Println("Проверка функции GerFreqDistr:", resX, "\n")

	resY := GetFreqDistr(RParamsArr, A, B, K)
	fmt.Println("Значение функции GerFreqDistr для последовательности при N=100:", resY)

	var values plotter.Values
	values = append(values, resY...)
	hist, _ := plotter.NewBarChart(values, 10)

	pl := plot.New()
	pl.Add(hist)
	pl.Save(5*vg.Inch, 5*vg.Inch, "hist_e2.png")
	fmt.Println("Гистограмма относительных частот для последовательности случайных чисел длинной N=100 была сохранена в файл hist_e2.png")

	var pearsonCriterion_e2 float64
	for i := 0; i < K; i++ {
		pearsonCriterion_e2 += ((1/float64(K) - resY[i]) * (1/float64(K) - resY[i])) / resY[i]
	}
	fmt.Println("Критерий Пирсона для последовательности случайных чисел длинной N=100 равен:", pearsonCriterion_e2, "\n")

	resY = GetFreqDistr(RParamsArr_e3, A, B, K)
	fmt.Println("Значение функции GerFreqDistr для последовательности при N=1000:", resY)

	values = values[:0]
	values = append(values, resY...)
	hist, _ = plotter.NewBarChart(values, 10)

	pl = plot.New()
	pl.Add(hist)
	pl.Save(5*vg.Inch, 5*vg.Inch, "hist_e3.png")
	fmt.Println("Гистограмма относительных частот для последовательности случайных чисел длинной N=1000 была сохранена в файл hist_e3.png")

	var pearsonCriterion_e3 float64
	for i := 0; i < K; i++ {
		pearsonCriterion_e3 += ((1/float64(K) - resY[i]) * (1/float64(K) - resY[i])) / resY[i]
	}
	fmt.Println("Критерий Пирсона для последовательности случайных чисел длинной N=1000 равен:", pearsonCriterion_e3, "\n")

	resY = GetFreqDistr(RParamsArr_e4, A, B, K)
	fmt.Println("Значение функции GerFreqDistr для последовательности при N=10000:", resY)

	values = values[:0]
	values = append(values, resY...)
	hist, _ = plotter.NewBarChart(values, 10)

	pl = plot.New()
	pl.Add(hist)
	pl.Save(5*vg.Inch, 5*vg.Inch, "hist_e4.png")
	fmt.Println("Гистограмма относительных частот для последовательности случайных чисел длинной N=10000 была сохранена в файл hist_e4.png")

	var pearsonCriterion_e4 float64
	for i := 0; i < K; i++ {
		pearsonCriterion_e4 += ((1/float64(K) - resY[i]) * (1/float64(K) - resY[i])) / resY[i]
	}
	fmt.Println("Критерий Пирсона для последовательности случайных чисел длинной N=10000 равен:", pearsonCriterion_e4, "\n")

	resY = GetFreqDistr(RParamsArr_e5, A, B, K)
	fmt.Println("Значение функции GerFreqDistr для последовательности при N=100000:", resY)

	values = values[:0]
	values = append(values, resY...)
	hist, _ = plotter.NewBarChart(values, 10)

	pl = plot.New()
	pl.Add(hist)
	pl.Save(5*vg.Inch, 5*vg.Inch, "hist_e5.png")
	fmt.Println("Гистограмма относительных частот для последовательности случайных чисел длинной N=100000 была сохранена в файл hist_e5.png")

	var pearsonCriterion_e5 float64
	for i := 0; i < K; i++ {
		pearsonCriterion_e5 += ((1/float64(K) - resY[i]) * (1/float64(K) - resY[i])) / resY[i]
	}
	fmt.Println("Критерий Пирсона для последовательности случайных чисел длинной N=100000 равен:", pearsonCriterion_e5, "\n")

	// Задание 7

	N = 100

	rand_RParamsArr_e2 := make([]float64, N)
	for i := 0; i < N; i++ {
		rand_RParamsArr_e2[i] = A + (B-A)*rand.Float64()
	}

	fmt.Println("Максимальное значение при N=100 для встроенного генератора случайных чисел Go:", slices.Max(rand_RParamsArr_e2))
	fmt.Println("Минимальное значение при N=100 для встроенного генератора случайных чисел Go:", slices.Min(rand_RParamsArr_e2), "\n")

	N = 1000

	rand_RParamsArr_e3 := make([]float64, N)
	for i := 0; i < N; i++ {
		rand_RParamsArr_e3[i] = A + (B-A)*rand.Float64()
	}

	fmt.Println("Максимальное значение при N=1000 для встроенного генератора случайных чисел Go:", slices.Max(rand_RParamsArr_e3))
	fmt.Println("Минимальное значение при N=1000 для встроенного генератора случайных чисел Go:", slices.Min(rand_RParamsArr_e3), "\n")

	N = 10000

	rand_RParamsArr_e4 := make([]float64, N)
	for i := 0; i < N; i++ {
		rand_RParamsArr_e4[i] = A + (B-A)*rand.Float64()
	}

	fmt.Println("Максимальное значение при N=10000 для встроенного генератора случайных чисел Go:", slices.Max(rand_RParamsArr_e4))
	fmt.Println("Минимальное значение при N=10000 для встроенного генератора случайных чисел Go:", slices.Min(rand_RParamsArr_e4), "\n")

	N = 100000

	rand_RParamsArr_e5 := make([]float64, N)
	for i := 0; i < N; i++ {
		rand_RParamsArr_e5[i] = A + (B-A)*rand.Float64()
	}

	fmt.Println("Максимальное значение при N=100000 для встроенного генератора случайных чисел Go:", slices.Max(rand_RParamsArr_e5))
	fmt.Println("Минимальное значение при N=100000 для встроенного генератора случайных чисел Go:", slices.Min(rand_RParamsArr_e5), "\n")

	fmt.Println("Теоретическое мат. ожидание:", M)
	fmt.Println("Теоретическая дисперсия:", D, "\n")

	N = 100

	var rand_M_e2 float64
	for i := 0; i < N; i++ {
		rand_M_e2 += rand_RParamsArr_e2[i]
	}
	rand_M_e2 /= float64(N)
	fmt.Println("Мат. ожидание при N=100 для встроенного генератора случайных чисел Go:", rand_M_e2)

	var rand_EpsM1 float64 = math.Abs((M-rand_M_e2)/M) * 100
	fmt.Println("Погрешность мат. ожидания при N=100 для встроенного генератора случайных чисел Go:", rand_EpsM1)

	var rand_D_e2 float64
	for i := 0; i < N; i++ {
		rand_D_e2 += rand_RParamsArr_e2[i] * rand_RParamsArr_e2[i]
	}
	rand_D_e2 /= float64(N)
	rand_D_e2 -= rand_M_e2 * rand_M_e2
	rand_D_e2 *= float64(N) / (float64(N) - 1)
	fmt.Println("Дисперсия при N=100 для встроенного генератора случайных чисел Go:", rand_D_e2)

	var rand_EpsD1 float64 = math.Abs((D-rand_D_e2)/D) * 100
	fmt.Println("Погрешность дисперсии при N=100 для встроенного генератора случайных чисел Go:", rand_EpsD1, "\n")

	N = 1000

	var rand_M_e3 float64
	for i := 0; i < N; i++ {
		rand_M_e3 += rand_RParamsArr_e3[i]
	}
	rand_M_e3 /= float64(N)
	fmt.Println("Мат. ожидание при N=1000 для встроенного генератора случайных чисел Go:", rand_M_e3)

	var rand_EpsM2 float64 = math.Abs((M-rand_M_e3)/M) * 100
	fmt.Println("Погрешность мат. ожидания при N=1000 для встроенного генератора случайных чисел Go:", rand_EpsM2)

	var rand_D_e3 float64
	for i := 0; i < N; i++ {
		rand_D_e3 += rand_RParamsArr_e3[i] * rand_RParamsArr_e3[i]
	}
	rand_D_e3 /= float64(N)
	rand_D_e3 -= rand_M_e3 * rand_M_e3
	rand_D_e3 *= float64(N) / (float64(N) - 1)
	fmt.Println("Дисперсия при N=1000 для встроенного генератора случайных чисел Go:", rand_D_e3)

	var rand_EpsD2 float64 = math.Abs((D-rand_D_e3)/D) * 100
	fmt.Println("Погрешность дисперсии при N=1000 для встроенного генератора случайных чисел Go:", rand_EpsD2, "\n")

	N = 10000

	var rand_M_e4 float64
	for i := 0; i < N; i++ {
		rand_M_e4 += rand_RParamsArr_e4[i]
	}
	rand_M_e4 /= float64(N)
	fmt.Println("Мат. ожидание при N=10000 для встроенного генератора случайных чисел Go:", rand_M_e4)

	var rand_EpsM3 float64 = math.Abs((M-rand_M_e4)/M) * 100
	fmt.Println("Погрешность мат. ожидания при N=10000 для встроенного генератора случайных чисел Go:", rand_EpsM3)

	var rand_D_e4 float64
	for i := 0; i < N; i++ {
		rand_D_e4 += rand_RParamsArr_e4[i] * rand_RParamsArr_e4[i]
	}
	rand_D_e4 /= float64(N)
	rand_D_e4 -= rand_M_e4 * rand_M_e4
	rand_D_e4 *= float64(N) / (float64(N) - 1)
	fmt.Println("Дисперсия при N=10000 для встроенного генератора случайных чисел Go:", rand_D_e4)

	var rand_EpsD3 float64 = math.Abs((D-rand_D_e4)/D) * 100
	fmt.Println("Погрешность дисперсии при N=10000 для встроенного генератора случайных чисел Go:", rand_EpsD3, "\n")

	N = 100000

	var rand_M_e5 float64
	for i := 0; i < N; i++ {
		rand_M_e5 += rand_RParamsArr_e5[i]
	}
	rand_M_e5 /= float64(N)
	fmt.Println("Мат. ожидание при N=100000 для встроенного генератора случайных чисел Go:", rand_M_e5)

	var rand_EpsM4 float64 = math.Abs((M-rand_M_e5)/M) * 100
	fmt.Println("Погрешность мат. ожидания при N=100000 для встроенного генератора случайных чисел Go:", rand_EpsM4)

	var rand_D_e5 float64
	for i := 0; i < N; i++ {
		rand_D_e5 += rand_RParamsArr_e5[i] * rand_RParamsArr_e5[i]
	}
	rand_D_e5 /= float64(N)
	rand_D_e5 -= rand_M_e5 * rand_M_e5
	rand_D_e5 *= float64(N) / (float64(N) - 1)
	fmt.Println("Дисперсия при N=100000 для встроенного генератора случайных чисел Go:", rand_D_e5)

	var rand_EpsD4 float64 = math.Abs((D-rand_D_e5)/D) * 100
	fmt.Println("Погрешность дисперсии при N=100000 для встроенного генератора случайных чисел Go:", rand_EpsD4, "\n")

	var rand_TEST_1 = RANDPeriod(rand_RParamsArr_e2)
	fmt.Println("Результаты теста на периодичность последовательности при N=100 для встроенного генератора случайных чисел Go:", rand_TEST_1)
	var rand_TEST_2 = RANDPeriod(rand_RParamsArr_e3)
	fmt.Println("Результаты теста на периодичность последовательности при N=1000 для встроенного генератора случайных чисел Go:", rand_TEST_2)
	var rand_TEST_3 = RANDPeriod(rand_RParamsArr_e4)
	fmt.Println("Результаты теста на периодичность последовательности при N=10000 для встроенного генератора случайных чисел Go:", rand_TEST_3)
	var rand_TEST_4 = RANDPeriod(rand_RParamsArr_e5)
	fmt.Println("Результаты теста на периодичность последовательности при N=100000 для встроенного генератора случайных чисел Go:", rand_TEST_4, "\n")

	resY = GetFreqDistr(rand_RParamsArr_e2, A, B, K)
	fmt.Println("Значение функции GerFreqDistr для последовательности при N=100 для встроенного генератора случайных чисел Go:", resY)

	values = values[:0]
	values = append(values, resY...)
	hist, _ = plotter.NewBarChart(values, 10)

	pl = plot.New()
	pl.Add(hist)
	pl.Save(5*vg.Inch, 5*vg.Inch, "rand_hist_e2.png")
	fmt.Println("Гистограмма относительных частот для последовательности случайных чисел длинной N=100 для встроенного генератора случайных чисел Go была сохранена в файл rand_hist_e2.png")

	var rand_pearsonCriterion_e2 float64
	for i := 0; i < K; i++ {
		rand_pearsonCriterion_e2 += ((1/float64(K) - resY[i]) * (1/float64(K) - resY[i])) / resY[i]
	}
	fmt.Println("Критерий Пирсона для последовательности случайных чисел длинной N=100 для встроенного генератора случайных чисел Go равен:", rand_pearsonCriterion_e2, "\n")

	resY = GetFreqDistr(rand_RParamsArr_e3, A, B, K)
	fmt.Println("Значение функции GerFreqDistr для последовательности при N=1000 для встроенного генератора случайных чисел Go:", resY)

	values = values[:0]
	values = append(values, resY...)
	hist, _ = plotter.NewBarChart(values, 10)

	pl = plot.New()
	pl.Add(hist)
	pl.Save(5*vg.Inch, 5*vg.Inch, "rand_hist_e3.png")
	fmt.Println("Гистограмма относительных частот для последовательности случайных чисел длинной N=1000 для встроенного генератора случайных чисел Go была сохранена в файл rand_hist_e3.png")

	var rand_pearsonCriterion_e3 float64
	for i := 0; i < K; i++ {
		rand_pearsonCriterion_e3 += ((1/float64(K) - resY[i]) * (1/float64(K) - resY[i])) / resY[i]
	}
	fmt.Println("Критерий Пирсона для последовательности случайных чисел длинной N=1000 для встроенного генератора случайных чисел Go равен:", rand_pearsonCriterion_e3, "\n")

	resY = GetFreqDistr(rand_RParamsArr_e4, A, B, K)
	fmt.Println("Значение функции GerFreqDistr для последовательности при N=10000 для встроенного генератора случайных чисел Go:", resY)

	values = values[:0]
	values = append(values, resY...)
	hist, _ = plotter.NewBarChart(values, 10)

	pl = plot.New()
	pl.Add(hist)
	pl.Save(5*vg.Inch, 5*vg.Inch, "rand_hist_e4.png")
	fmt.Println("Гистограмма относительных частот для последовательности случайных чисел длинной N=10000 для встроенного генератора случайных чисел Go была сохранена в файл rand_hist_e4.png")

	var rand_pearsonCriterion_e4 float64
	for i := 0; i < K; i++ {
		rand_pearsonCriterion_e4 += ((1/float64(K) - resY[i]) * (1/float64(K) - resY[i])) / resY[i]
	}
	fmt.Println("Критерий Пирсона для последовательности случайных чисел длинной N=10000 для встроенного генератора случайных чисел Go равен:", rand_pearsonCriterion_e4, "\n")

	resY = GetFreqDistr(rand_RParamsArr_e5, A, B, K)
	fmt.Println("Значение функции GerFreqDistr для последовательности при N=100000 для встроенного генератора случайных чисел Go:", resY)

	values = values[:0]
	values = append(values, resY...)
	hist, _ = plotter.NewBarChart(values, 10)

	pl = plot.New()
	pl.Add(hist)
	pl.Save(5*vg.Inch, 5*vg.Inch, "rand_hist_e5.png")
	fmt.Println("Гистограмма относительных частот для последовательности случайных чисел длинной N=100000 для встроенного генератора случайных чисел Go была сохранена в файл rand_hist_e5.png")

	var rand_pearsonCriterion_e5 float64
	for i := 0; i < K; i++ {
		rand_pearsonCriterion_e5 += ((1/float64(K) - resY[i]) * (1/float64(K) - resY[i])) / resY[i]
	}
	fmt.Println("Критерий Пирсона для последовательности случайных чисел длинной N=100000 для встроенного генератора случайных чисел Go равен:", rand_pearsonCriterion_e5, "\n")

}
