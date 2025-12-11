package main

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"math"
	"math/rand"
	"time"
)

// ========== ЗАДАНИЕ 1 ==========

// WeibullPDF - функция плотности вероятности распределения Вейбулла
func WeibullPDF(x, lambda, k float64) float64 {
	if x < 0 || lambda <= 0 || k <= 0 {
		return 0
	}

	// f(x;λ,k) = (k/λ) * (x/λ)^(k-1) * exp(-(x/λ)^k)
	if x == 0 && k > 1 {
		return 0
	}
	if x == 0 && k == 1 {
		return 1 / lambda
	}
	if x == 0 && k < 1 {
		return math.Inf(1)
	}

	term1 := k / lambda
	term2 := math.Pow(x/lambda, k-1)
	term3 := math.Exp(-math.Pow(x/lambda, k))

	return term1 * term2 * term3
}

// ========== ЗАДАНИЕ 2 ==========

// InverseWeibull - функция обратного преобразования для распределения Вейбулла
func InverseWeibull(p, lambda, k float64) float64 {
	if p <= 0 || p >= 1 || lambda <= 0 || k <= 0 {
		return 0
	}
	// x = λ * (-ln(1-p))^(1/k)
	return lambda * math.Pow(-math.Log(1-p), 1/k)
}

// ========== ЗАДАНИЕ 3 ==========

// GenerateWeibullDistribution - генерация распределения Вейбулла методом обратной функции
func GenerateWeibullDistribution(lambda, k float64, n int) []float64 {
	data := make([]float64, n)
	for i := 0; i < n; i++ {
		// Генерируем равномерно распределенное число в (0,1)
		u := rand.Float64()
		// Применяем обратную функцию
		data[i] = InverseWeibull(u, lambda, k)
	}
	return data
}

// ========== ЗАДАНИЕ 4 ==========

// CalculateWeibullHistogram - вычисление гистограммы для распределения Вейбулла
func CalculateWeibullHistogram(data []float64, bins int, maxVal float64) ([]float64, []float64) {
	// Инициализируем массив для подсчета
	counts := make([]int, bins)
	binWidth := maxVal / float64(bins)

	// Подсчитываем попадания в бины
	for _, value := range data {
		if value >= 0 && value < maxVal {
			binIndex := int(value / binWidth)
			if binIndex >= bins {
				binIndex = bins - 1
			}
			counts[binIndex]++
		}
	}

	// Преобразуем в относительные частоты (плотность вероятности)
	frequencies := make([]float64, bins)
	binCenters := make([]float64, bins)
	total := float64(len(data))

	for i := 0; i < bins; i++ {
		binCenters[i] = (float64(i) + 0.5) * binWidth
		frequencies[i] = float64(counts[i]) / (total * binWidth)
	}

	return binCenters, frequencies
}

// CalculateWeibullRMSE - расчет RMSE между экспериментальным и теоретическим распределением Вейбулла
func CalculateWeibullRMSE(experimental []float64, binCenters []float64, lambda, k float64) float64 {
	n := len(experimental)
	if n == 0 {
		return 0
	}

	var sum float64
	for i := 0; i < n; i++ {
		// Теоретическое значение плотности вероятности
		theoryVal := WeibullPDF(binCenters[i], lambda, k)
		// Квадрат разности
		diff := experimental[i] - theoryVal
		sum += diff * diff
	}

	return math.Sqrt(sum / float64(n))
}

// Вспомогательная функция для вычисления статистик
func calculateStatistics(data []float64) (mean, variance, stdDev float64) {
	n := float64(len(data))
	if n == 0 {
		return 0, 0, 0
	}

	var sum, sumSq float64
	for _, val := range data {
		sum += val
		sumSq += val * val
	}

	mean = sum / n
	variance = sumSq/n - mean*mean
	if variance < 0 {
		variance = 0
	}
	stdDev = math.Sqrt(variance)

	return mean, variance, stdDev
}

// Вспомогательная структура для цвета
type color struct {
	R, G, B uint8
}

// Метод для преобразования цвета в RGBA
func (c color) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	r |= r << 8
	g = uint32(c.G)
	g |= g << 8
	b = uint32(c.B)
	b |= b << 8
	a = uint32(255)
	a |= a << 8
	return
}

func main() {
	fmt.Println("Практическая работа №6")
	fmt.Println("Моделирование закона распределения Вейбулла\n")

	// Инициализация генератора случайных чисел
	rand.Seed(time.Now().UnixNano())

	// ========== ЗАДАНИЕ 1 ==========
	fmt.Println("=== ЗАДАНИЕ 1 ===")
	fmt.Println("Функция плотности вероятности распределения Вейбулла")

	// Параметры из задания
	params := []struct {
		lambda float64
		k      float64
		label  string
		color  color
	}{
		{2.0, 2.0, "λ=2, k=2", color{255, 0, 0}}, // Красный
		{2.0, 1.0, "λ=2, k=1", color{0, 0, 255}}, // Синий
		{1.0, 1.0, "λ=1, k=1", color{0, 128, 0}}, // Зеленый
	}

	// Создаем график для PDF
	p1 := plot.New()
	p1.Title.Text = "Плотность вероятности распределения Вейбулла"
	p1.X.Label.Text = "x"
	p1.Y.Label.Text = "f(x)"
	p1.Legend.Top = true
	p1.Legend.Left = true

	// Диапазон для построения графиков
	xMax := 5.0
	points := 200

	for _, param := range params {
		pts := make(plotter.XYs, points)
		count := 0

		for i := 0; i < points; i++ {
			x := float64(i) * xMax / float64(points-1)
			y := WeibullPDF(x, param.lambda, param.k)

			// Проверяем на особые случаи
			if math.IsInf(y, 0) || math.IsNaN(y) {
				continue
			}

			pts[count].X = x
			pts[count].Y = y
			count++
		}

		pts = pts[:count]

		line, err := plotter.NewLine(pts)
		if err != nil {
			fmt.Printf("Ошибка создания линии: %v\n", err)
			continue
		}

		line.Color = param.color
		line.Width = vg.Points(1.5)
		p1.Add(line)
		p1.Legend.Add(param.label, line)
	}

	// Сохраняем график
	if err := p1.Save(10*vg.Inch, 6*vg.Inch, "task1_weibull_pdf.png"); err != nil {
		fmt.Printf("Ошибка сохранения графика: %v\n", err)
	} else {
		fmt.Println("График плотности вероятности сохранен: task1_weibull_pdf.png")
	}

	// ========== ЗАДАНИЕ 2 ==========
	fmt.Println("\n=== ЗАДАНИЕ 2 ===")
	fmt.Println("Функция обратного преобразования для λ=2, k=2")

	lambda1 := 2.0
	k1 := 2.0

	// Создаем график для обратной функции
	p2 := plot.New()
	p2.Title.Text = fmt.Sprintf("Обратная функция Вейбулла (λ=%.1f, k=%.1f)", lambda1, k1)
	p2.X.Label.Text = "p ∈ (0,1)"
	p2.Y.Label.Text = "x = F⁻¹(p)"

	// Строим обратную функцию
	points2 := 100
	pts2 := make(plotter.XYs, points2-2) // Исключаем p=0 и p=1
	idx := 0

	for i := 1; i < points2-1; i++ {
		p := float64(i) / float64(points2)
		x := InverseWeibull(p, lambda1, k1)

		pts2[idx].X = p
		pts2[idx].Y = x
		idx++
	}

	line2, err := plotter.NewLine(pts2)
	if err != nil {
		fmt.Printf("Ошибка создания линии: %v\n", err)
	} else {
		line2.Color = color{255, 0, 0}
		line2.Width = vg.Points(2)
		p2.Add(line2)
	}

	// Сохраняем график
	if err := p2.Save(10*vg.Inch, 6*vg.Inch, "task2_inverse_weibull.png"); err != nil {
		fmt.Printf("Ошибка сохранения графика: %v\n", err)
	} else {
		fmt.Println("График обратной функции сохранен: task2_inverse_weibull.png")
	}

	// Тестовые значения для проверки
	fmt.Println("\nТестовые значения обратной функции:")
	testPs := []float64{0.1, 0.5, 0.9}
	for _, p := range testPs {
		x := InverseWeibull(p, lambda1, k1)
		fmt.Printf("  F⁻¹(%.1f) = %.4f\n", p, x)
	}

	// Проверка: при k=1 распределение Вейбулла должно совпадать с экспоненциальным
	fmt.Println("\nПроверка частного случая (k=1 -> экспоненциальное распределение):")
	p_check := 0.5
	x_weibull := InverseWeibull(p_check, 2.0, 1.0)
	// Для экспоненциального: x = -ln(1-p)/λ
	x_exp := -math.Log(1-p_check) / 2.0
	fmt.Printf("  Вейбулл(λ=2, k=1): F⁻¹(%.1f) = %.4f\n", p_check, x_weibull)
	fmt.Printf("  Экспоненциальное(λ=2): F⁻¹(%.1f) = %.4f\n", p_check, x_exp)
	fmt.Printf("  Совпадают: %v\n", math.Abs(x_weibull-x_exp) < 1e-10)

	// ========== ЗАДАНИЕ 3 ==========
	fmt.Println("\n=== ЗАДАНИЕ 3 ===")
	fmt.Println("Моделирование распределения Вейбулла методом обратной функции")

	// Различные значения N из задания
	Ns := []int{1000, 10000, 100000, 1000000}

	// Массивы для хранения сгенерированных данных
	generatedData := make([][]float64, len(Ns))

	fmt.Println("\nГенерация распределения Вейбулла (λ=2, k=2):")
	for idx, N := range Ns {
		fmt.Printf("  N=%d... ", N)
		start := time.Now()

		data := GenerateWeibullDistribution(lambda1, k1, N)
		generatedData[idx] = data

		// Вычисляем статистики
		mean, variance, stdDev := calculateStatistics(data)

		// Теоретические моменты распределения Вейбулла
		// μ = λ * Γ(1 + 1/k)
		// σ² = λ² * [Γ(1 + 2/k) - Γ²(1 + 1/k)]
		gamma1 := math.Gamma(1 + 1/k1)
		gamma2 := math.Gamma(1 + 2/k1)
		theoryMean := lambda1 * gamma1
		theoryVar := lambda1 * lambda1 * (gamma2 - gamma1*gamma1)
		theoryStdDev := math.Sqrt(theoryVar)

		elapsed := time.Since(start)

		fmt.Printf("завершено за %v\n", elapsed)
		fmt.Printf("    Выборочное среднее: %.4f (теоретическое: %.4f)\n", mean, theoryMean)
		fmt.Printf("    Выборочная дисперсия: %.4f (теоретическая: %.4f)\n", variance, theoryVar)
		fmt.Printf("    Выборочное СКО: %.4f (теоретическое: %.4f)\n", stdDev, theoryStdDev)
		fmt.Printf("    Относительная ошибка среднего: %.2f%%\n",
			math.Abs(mean-theoryMean)/theoryMean*100)
	}

	// ========== ЗАДАНИЕ 4 ==========
	fmt.Println("\n=== ЗАДАНИЕ 4 ===")
	fmt.Println("Гистограммы относительных частот и расчет RMSE")

	// Параметры для гистограммы (20 интервалов по заданию)
	bins := 20
	histMax := 5.0 // Максимальное значение для гистограммы

	// Массив для хранения значений RMSE
	rmseValues := make([]float64, len(Ns))

	// Создаем график для гистограмм
	p3 := plot.New()
	p3.Title.Text = "Гистограммы распределения Вейбулла (λ=2, k=2)"
	p3.X.Label.Text = "x"
	p3.Y.Label.Text = "Плотность вероятности"
	p3.Legend.Top = true
	p3.Legend.Left = true

	// Также создаем отдельный график для зависимости RMSE от N
	p4 := plot.New()
	p4.Title.Text = "Зависимость RMSE от числа экспериментов"
	p4.X.Label.Text = "Число экспериментов (N)"
	p4.Y.Label.Text = "RMSE"
	p4.X.Scale = plot.LogScale{}
	p4.Y.Scale = plot.LogScale{}

	// Цвета для разных N
	histColors := []color{
		{255, 0, 0},   // N=10^3 - красный
		{0, 0, 255},   // N=10^4 - синий
		{0, 128, 0},   // N=10^5 - зеленый
		{255, 0, 255}, // N=10^6 - фиолетовый
	}

	fmt.Println("\nРасчет гистограмм и RMSE:")
	for idx, N := range Ns {
		data := generatedData[idx]

		// Вычисляем гистограмму
		binCenters, frequencies := CalculateWeibullHistogram(data, bins, histMax)

		// Вычисляем RMSE
		rmse := CalculateWeibullRMSE(frequencies, binCenters, lambda1, k1)
		rmseValues[idx] = rmse

		fmt.Printf("  N=%d: RMSE = %.6f\n", N, rmse)

		// Добавляем гистограмму на график (только для N=10^5 для наглядности)
		if N == 100000 {
			// Создаем точки для гистограммы
			histPts := make(plotter.XYs, bins)
			validCount := 0

			for i := 0; i < bins; i++ {
				if !math.IsNaN(frequencies[i]) && !math.IsInf(frequencies[i], 0) {
					histPts[validCount].X = binCenters[i]
					histPts[validCount].Y = frequencies[i]
					validCount++
				}
			}

			histPts = histPts[:validCount]

			if validCount > 0 {
				histLine, err := plotter.NewLine(histPts)
				if err != nil {
					fmt.Printf("Ошибка создания линии гистограммы: %v\n", err)
				} else {
					histLine.Color = histColors[idx]
					histLine.Width = vg.Points(1)
					histLine.Dashes = []vg.Length{vg.Points(2), vg.Points(2)}
					p3.Add(histLine)
					p3.Legend.Add(fmt.Sprintf("Экспериментальное (N=%d)", N), histLine)
				}
			}
		}

		// Сохраняем отдельную гистограмму для каждого N
		saveIndividualWeibullHistogram(data, bins, histMax, lambda1, k1, N, idx)
	}

	// Добавляем теоретическую кривую на график гистограмм
	theoryPoints := 200
	theoryPts := make(plotter.XYs, theoryPoints)
	validTheory := 0

	for i := 0; i < theoryPoints; i++ {
		x := float64(i) * histMax / float64(theoryPoints-1)
		y := WeibullPDF(x, lambda1, k1)

		if !math.IsNaN(y) && !math.IsInf(y, 0) {
			theoryPts[validTheory].X = x
			theoryPts[validTheory].Y = y
			validTheory++
		}
	}

	theoryPts = theoryPts[:validTheory]

	if validTheory > 0 {
		theoryLine, err := plotter.NewLine(theoryPts)
		if err != nil {
			fmt.Printf("Ошибка создания теоретической линии: %v\n", err)
		} else {
			theoryLine.Color = color{0, 0, 0}
			theoryLine.Width = vg.Points(2)
			p3.Add(theoryLine)
			p3.Legend.Add("Теоретическое", theoryLine)
		}
	}

	// Сохраняем график гистограмм
	if err := p3.Save(10*vg.Inch, 6*vg.Inch, "task4_weibull_histograms.png"); err != nil {
		fmt.Printf("Ошибка сохранения графика гистограмм: %v\n", err)
	} else {
		fmt.Println("\nГрафик сравнения гистограмм сохранен: task4_weibull_histograms.png")
	}

	// Добавляем точки RMSE на график зависимости
	rmsePts := make(plotter.XYs, len(Ns))
	for i, N := range Ns {
		rmsePts[i].X = float64(N)
		rmsePts[i].Y = rmseValues[i]
	}

	scatter, err := plotter.NewScatter(rmsePts)
	if err != nil {
		fmt.Printf("Ошибка создания scatter plot: %v\n", err)
	} else {
		scatter.GlyphStyle.Color = color{255, 0, 0}
		scatter.GlyphStyle.Radius = vg.Points(4)
		scatter.GlyphStyle.Shape = draw.CircleGlyph{}
		p4.Add(scatter)

		// Добавляем линию тренда
		trendLine, err := plotter.NewLine(rmsePts)
		if err != nil {
			fmt.Printf("Ошибка создания линии тренда: %v\n", err)
		} else {
			trendLine.Color = color{0, 0, 255}
			trendLine.Width = vg.Points(1)
			trendLine.Dashes = []vg.Length{vg.Points(2), vg.Points(2)}
			p4.Add(trendLine)
		}
	}

	// Сохраняем график зависимости RMSE от N
	if err := p4.Save(10*vg.Inch, 6*vg.Inch, "task4_weibull_rmse_vs_n.png"); err != nil {
		fmt.Printf("Ошибка сохранения графика RMSE: %v\n", err)
	} else {
		fmt.Println("График зависимости RMSE от N сохранен: task4_weibull_rmse_vs_n.png")
	}

	// Дополнительный анализ: проверка для разных параметров
	fmt.Println("\n=== ДОПОЛНИТЕЛЬНЫЙ АНАЛИЗ ===")
	fmt.Println("Анализ для различных параметров распределения Вейбулла")

	for _, param := range params {
		if param.lambda == lambda1 && param.k == k1 {
			continue // Пропускаем основной случай
		}

		fmt.Printf("\nλ=%.1f, k=%.1f:\n", param.lambda, param.k)

		// Генерируем данные
		N_test := 100000
		data := GenerateWeibullDistribution(param.lambda, param.k, N_test)

		// Вычисляем статистики
		mean, variance, _ := calculateStatistics(data)

		// Теоретические моменты
		gamma1 := math.Gamma(1 + 1/param.k)
		gamma2 := math.Gamma(1 + 2/param.k)
		theoryMean := param.lambda * gamma1
		theoryVar := param.lambda * param.lambda * (gamma2 - gamma1*gamma1)

		fmt.Printf("  Теоретическое среднее: %.4f\n", theoryMean)
		fmt.Printf("  Выборочное среднее: %.4f\n", mean)
		fmt.Printf("  Теоретическая дисперсия: %.4f\n", theoryVar)
		fmt.Printf("  Выборочная дисперсия: %.4f\n", variance)
		fmt.Printf("  Относительная ошибка среднего: %.2f%%\n",
			math.Abs(mean-theoryMean)/theoryMean*100)

		// Специальный случай: k=1 (экспоненциальное распределение)
		if param.k == 1 {
			fmt.Printf("  При k=1 распределение должно быть экспоненциальным с λ=%.1f\n", param.lambda)
			expMean := 1 / param.lambda
			expVar := 1 / (param.lambda * param.lambda)
			fmt.Printf("  Теоретическое экспоненциальное среднее: %.4f\n", expMean)
			fmt.Printf("  Теоретическая экспоненциальная дисперсия: %.4f\n", expVar)
		}
	}

	// Анализ формы распределения при разных k
	fmt.Println("\n=== АНАЛИЗ ВЛИЯНИЯ ПАРАМЕТРА ФОРМЫ k ===")

	const testLambda = 2.0
	testKs := []float64{0.5, 1.0, 1.5, 2.0, 3.0, 5.0}

	for _, k := range testKs {
		// Мода распределения Вейбулла
		var mode float64
		if k > 1 {
			mode = testLambda * math.Pow((k-1)/k, 1/k)
		} else if k == 1 {
			mode = 0
		} else { // k < 1
			mode = 0 // функция не определена или бесконечна в 0
		}

		// Медиана
		median := testLambda * math.Pow(math.Log(2), 1/k)

		fmt.Printf("λ=%.1f, k=%.1f: мода=%.4f, медиана=%.4f\n",
			testLambda, k, mode, median)
	}

	fmt.Println("\n=== ПРАКТИЧЕСКАЯ РАБОТА ЗАВЕРШЕНА ===")
	fmt.Println("Созданы файлы:")
	fmt.Println("1. task1_weibull_pdf.png - плотности вероятности для разных параметров")
	fmt.Println("2. task2_inverse_weibull.png - обратная функция для λ=2, k=2")
	fmt.Println("3. task4_weibull_histogram_N*.png - гистограммы для разных N")
	fmt.Println("4. task4_weibull_histograms.png - сравнение гистограмм")
	fmt.Println("5. task4_weibull_rmse_vs_n.png - зависимость RMSE от N")
}

// Функция для сохранения отдельных гистограмм
func saveIndividualWeibullHistogram(data []float64, bins int, maxVal, lambda, k float64, N, idx int) {
	// Вычисляем гистограмму
	binCenters, frequencies := CalculateWeibullHistogram(data, bins, maxVal)

	// Создаем график
	p := plot.New()
	p.Title.Text = fmt.Sprintf("Распределение Вейбулла, λ=%.1f, k=%.1f, N=%d", lambda, k, N)
	p.X.Label.Text = "x"
	p.Y.Label.Text = "Плотность вероятности"
	p.Legend.Top = true

	// Добавляем экспериментальную гистограмму
	expPts := make(plotter.XYs, 0, bins)
	for i := 0; i < bins; i++ {
		if !math.IsNaN(frequencies[i]) && !math.IsInf(frequencies[i], 0) {
			expPts = append(expPts, plotter.XY{X: binCenters[i], Y: frequencies[i]})
		}
	}

	if len(expPts) > 0 {
		expLine, err := plotter.NewLine(expPts)
		if err != nil {
			return
		}

		expLine.Color = color{255, 0, 0}
		expLine.Width = vg.Points(1)
		p.Add(expLine)
		p.Legend.Add("Экспериментальное", expLine)
	}

	// Добавляем теоретическую кривую
	theoryPoints := 200
	theoryPts := make(plotter.XYs, 0, theoryPoints)

	for i := 0; i < theoryPoints; i++ {
		x := float64(i) * maxVal / float64(theoryPoints-1)
		y := WeibullPDF(x, lambda, k)

		if !math.IsNaN(y) && !math.IsInf(y, 0) {
			theoryPts = append(theoryPts, plotter.XY{X: x, Y: y})
		}
	}

	if len(theoryPts) > 0 {
		theoryLine, err := plotter.NewLine(theoryPts)
		if err != nil {
			return
		}

		theoryLine.Color = color{0, 0, 255}
		theoryLine.Width = vg.Points(1.5)
		p.Add(theoryLine)
		p.Legend.Add("Теоретическое", theoryLine)
	}

	// Сохраняем график
	filename := fmt.Sprintf("task4_weibull_histogram_N%d.png", N)
	p.Save(8*vg.Inch, 6*vg.Inch, filename)
}
