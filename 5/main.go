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

// ExponentialPDF - функция плотности вероятности экспоненциального распределения
func ExponentialPDF(x, lambda float64) float64 {
	if x < 0 {
		return 0
	}
	return lambda * math.Exp(-lambda*x)
}

// ========== ЗАДАНИЕ 2 ==========

// InverseExponential - функция обратного преобразования для экспоненциального распределения
func InverseExponential(p, lambda float64) float64 {
	if p <= 0 || p >= 1 {
		return 0
	}
	return -math.Log(1-p) / lambda
}

// ========== ЗАДАНИЕ 3 ==========

// GenerateExponentialDistribution - генерация экспоненциально распределенных чисел методом обратной функции
func GenerateExponentialDistribution(lambda float64, n int) []float64 {
	data := make([]float64, n)
	for i := 0; i < n; i++ {
		// Генерируем равномерно распределенное число в (0,1)
		u := rand.Float64()
		// Применяем обратную функцию
		data[i] = InverseExponential(u, lambda)
	}
	return data
}

// ========== ЗАДАНИЕ 4 ==========

// CalculateHistogram - вычисление гистограммы
func CalculateHistogram(data []float64, bins int, maxVal float64) ([]float64, []float64) {
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

// CalculateRMSE - расчет среднеквадратического отклонения между экспериментальным и теоретическим распределением
func CalculateRMSE(experimental, theoretical []float64, binCenters []float64, lambda float64) float64 {
	n := len(experimental)
	if n != len(theoretical) {
		return 0
	}

	var sum float64
	for i := 0; i < n; i++ {
		// Теоретическое значение плотности вероятности
		theoryVal := ExponentialPDF(binCenters[i], lambda)
		// Квадрат разности
		diff := experimental[i] - theoryVal
		sum += diff * diff
	}

	return math.Sqrt(sum / float64(n))
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
	fmt.Println("Практическая работа №5")
	fmt.Println("Моделирование экспоненциального закона распределения\n")

	// Инициализация генератора случайных чисел
	rand.Seed(time.Now().UnixNano())

	// ========== ЗАДАНИЕ 1 ==========
	fmt.Println("=== ЗАДАНИЕ 1 ===")
	fmt.Println("Функция плотности вероятности экспоненциального распределения")

	// Параметры лямбда из задания
	lambdas := []float64{1.5, 1.0, 0.5}
	colors := []color{
		{255, 0, 0}, // Красный
		{0, 0, 255}, // Синий
		{0, 128, 0}, // Зеленый
	}

	// Создаем график для PDF
	p1 := plot.New()
	p1.Title.Text = "Плотность вероятности экспоненциального распределения"
	p1.X.Label.Text = "x"
	p1.Y.Label.Text = "f(x)"
	p1.Legend.Top = true

	// Диапазон для построения графиков
	xMax := 6.0
	points := 200

	for idx, lambda := range lambdas {
		pts := make(plotter.XYs, points)
		for i := 0; i < points; i++ {
			x := float64(i) * xMax / float64(points-1)
			pts[i].X = x
			pts[i].Y = ExponentialPDF(x, lambda)
		}

		line, err := plotter.NewLine(pts)
		if err != nil {
			fmt.Printf("Ошибка создания линии: %v\n", err)
			continue
		}

		line.Color = colors[idx]
		line.Width = vg.Points(1.5)
		p1.Add(line)
		p1.Legend.Add(fmt.Sprintf("λ=%.1f", lambda), line)
	}

	// Сохраняем график
	if err := p1.Save(10*vg.Inch, 6*vg.Inch, "task1_exponential_pdf.png"); err != nil {
		fmt.Printf("Ошибка сохранения графика: %v\n", err)
	} else {
		fmt.Println("График плотности вероятности сохранен: task1_exponential_pdf.png")
	}

	// ========== ЗАДАНИЕ 2 ==========
	fmt.Println("\n=== ЗАДАНИЕ 2 ===")
	fmt.Println("Функция обратного преобразования для λ=1.5")

	lambda1 := 1.5

	// Создаем график для обратной функции
	p2 := plot.New()
	p2.Title.Text = fmt.Sprintf("Обратная функция для экспоненциального распределения (λ=%.1f)", lambda1)
	p2.X.Label.Text = "p ∈ (0,1)"
	p2.Y.Label.Text = "x = F⁻¹(p)"

	// Строим обратную функцию
	points2 := 100
	pts2 := make(plotter.XYs, points2)
	for i := 1; i < points2; i++ { // Начинаем с 1, чтобы избежать p=0
		p := float64(i) / float64(points2)
		pts2[i].X = p
		pts2[i].Y = InverseExponential(p, lambda1)
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
	if err := p2.Save(10*vg.Inch, 6*vg.Inch, "task2_inverse_function.png"); err != nil {
		fmt.Printf("Ошибка сохранения графика: %v\n", err)
	} else {
		fmt.Println("График обратной функции сохранен: task2_inverse_function.png")
	}

	// Тестовые значения для проверки
	fmt.Println("\nТестовые значения обратной функции:")
	testPs := []float64{0.1, 0.5, 0.9}
	for _, p := range testPs {
		x := InverseExponential(p, lambda1)
		fmt.Printf("  F⁻¹(%.1f) = %.4f\n", p, x)
	}

	// ========== ЗАДАНИЕ 3 ==========
	fmt.Println("\n=== ЗАДАНИЕ 3 ===")
	fmt.Println("Моделирование методом обратной функции")

	// Различные значения N из задания
	Ns := []int{1000, 10000, 100000, 1000000}

	// Массивы для хранения сгенерированных данных
	generatedData := make([][]float64, len(Ns))

	fmt.Println("\nГенерация экспоненциально распределенных чисел:")
	for idx, N := range Ns {
		fmt.Printf("  N=%d... ", N)
		start := time.Now()

		data := GenerateExponentialDistribution(lambda1, N)
		generatedData[idx] = data

		// Вычисляем статистики
		var sum, sumSq float64
		for _, val := range data {
			sum += val
			sumSq += val * val
		}

		mean := sum / float64(N)
		variance := sumSq/float64(N) - mean*mean
		stdDev := math.Sqrt(math.Abs(variance))

		elapsed := time.Since(start)

		fmt.Printf("завершено за %v\n", elapsed)
		fmt.Printf("    Выборочное среднее: %.4f (теоретическое: %.4f)\n", mean, 1/lambda1)
		fmt.Printf("    Выборочное СКО: %.4f (теоретическое: %.4f)\n", stdDev, 1/lambda1)
	}

	// ========== ЗАДАНИЕ 4 ==========
	fmt.Println("\n=== ЗАДАНИЕ 4 ===")
	fmt.Println("Гистограммы относительных частот и расчет RMSE")

	// Параметры для гистограммы
	bins := 100
	histMax := 6.0 // Максимальное значение для гистограммы

	// Массив для хранения значений RMSE
	rmseValues := make([]float64, len(Ns))

	// Создаем график для гистограмм
	p3 := plot.New()
	p3.Title.Text = "Гистограммы экспоненциального распределения"
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
		binCenters, frequencies := CalculateHistogram(data, bins, histMax)

		// Вычисляем теоретическую плотность в центрах бинов
		theoretical := make([]float64, bins)
		for i := 0; i < bins; i++ {
			theoretical[i] = ExponentialPDF(binCenters[i], lambda1)
		}

		// Вычисляем RMSE
		rmse := CalculateRMSE(frequencies, theoretical, binCenters, lambda1)
		rmseValues[idx] = rmse

		fmt.Printf("  N=%d: RMSE = %.6f\n", N, rmse)

		// Добавляем гистограмму на график (только для N=10^5 для наглядности)
		if N == 100000 {
			// Создаем точки для гистограммы
			histPts := make(plotter.XYs, bins)
			for i := 0; i < bins; i++ {
				histPts[i].X = binCenters[i]
				histPts[i].Y = frequencies[i]
			}

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

		// Сохраняем отдельную гистограмму для каждого N
		saveIndividualHistogram(data, bins, histMax, lambda1, N, idx)
	}

	// Добавляем теоретическую кривую на график гистограмм
	theoryPts := make(plotter.XYs, 200)
	for i := 0; i < 200; i++ {
		x := float64(i) * histMax / 199.0
		theoryPts[i].X = x
		theoryPts[i].Y = ExponentialPDF(x, lambda1)
	}

	theoryLine, err := plotter.NewLine(theoryPts)
	if err != nil {
		fmt.Printf("Ошибка создания теоретической линии: %v\n", err)
	} else {
		theoryLine.Color = color{0, 0, 0}
		theoryLine.Width = vg.Points(2)
		p3.Add(theoryLine)
		p3.Legend.Add("Теоретическое", theoryLine)
	}

	// Сохраняем график гистограмм
	if err := p3.Save(10*vg.Inch, 6*vg.Inch, "task4_histograms_comparison.png"); err != nil {
		fmt.Printf("Ошибка сохранения графика гистограмм: %v\n", err)
	} else {
		fmt.Println("\nГрафик сравнения гистограмм сохранен: task4_histograms_comparison.png")
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
	if err := p4.Save(10*vg.Inch, 6*vg.Inch, "task4_rmse_vs_n.png"); err != nil {
		fmt.Printf("Ошибка сохранения графика RMSE: %v\n", err)
	} else {
		fmt.Println("График зависимости RMSE от N сохранен: task4_rmse_vs_n.png")
	}

	// Дополнительный анализ: проверка свойства отсутствия памяти
	fmt.Println("\n=== ДОПОЛНИТЕЛЬНЫЙ АНАЛИЗ ===")
	fmt.Println("Проверка свойства отсутствия памяти экспоненциального распределения")

	// Берем данные для N=100000
	N_test := 100000
	testData := GenerateExponentialDistribution(lambda1, N_test)

	// Проверяем свойство: P(X > s + t | X > s) = P(X > t)
	s := 1.0
	t := 0.5

	var count_s_plus_t, count_s int
	for _, val := range testData {
		if val > s {
			count_s++
			if val > s+t {
				count_s_plus_t++
			}
		}
	}

	if count_s > 0 {
		conditionalProb := float64(count_s_plus_t) / float64(count_s)
		unconditionalProb := math.Exp(-lambda1 * t)

		fmt.Printf("P(X > %.1f + %.1f | X > %.1f) = %.6f\n", s, t, s, conditionalProb)
		fmt.Printf("P(X > %.1f) = %.6f\n", t, unconditionalProb)
		fmt.Printf("Разница: %.6f\n", math.Abs(conditionalProb-unconditionalProb))
	}

	// Анализ для разных значений lambda
	fmt.Println("\n=== АНАЛИЗ ДЛЯ РАЗНЫХ λ ===")

	for _, lambda := range lambdas {
		fmt.Printf("\nλ = %.1f:\n", lambda)

		// Генерируем данные
		data := GenerateExponentialDistribution(lambda, 10000)

		// Вычисляем статистики
		var sum float64
		for _, val := range data {
			sum += val
		}
		mean := sum / 10000

		// Теоретические значения
		theoryMean := 1 / lambda
		theoryVar := 1 / (lambda * lambda)

		fmt.Printf("  Теоретическое среднее: %.4f\n", theoryMean)
		fmt.Printf("  Выборочное среднее: %.4f\n", mean)
		fmt.Printf("  Теоретическая дисперсия: %.4f\n", theoryVar)
		fmt.Printf("  Относительная ошибка: %.2f%%\n", math.Abs(mean-theoryMean)/theoryMean*100)
	}

	fmt.Println("\n=== ПРАКТИЧЕСКАЯ РАБОТА ЗАВЕРШЕНА ===")
	fmt.Println("Созданы файлы:")
	fmt.Println("1. task1_exponential_pdf.png - плотности вероятности для разных λ")
	fmt.Println("2. task2_inverse_function.png - обратная функция для λ=1.5")
	fmt.Println("3. task4_histogram_N*.png - гистограммы для разных N")
	fmt.Println("4. task4_histograms_comparison.png - сравнение гистограмм")
	fmt.Println("5. task4_rmse_vs_n.png - зависимость RMSE от N")
}

// Функция для сохранения отдельных гистограмм
func saveIndividualHistogram(data []float64, bins int, maxVal, lambda float64, N, idx int) {
	// Вычисляем гистограмму
	binCenters, frequencies := CalculateHistogram(data, bins, maxVal)

	// Создаем график
	p := plot.New()
	p.Title.Text = fmt.Sprintf("Экспоненциальное распределение, λ=%.1f, N=%d", lambda, N)
	p.X.Label.Text = "x"
	p.Y.Label.Text = "Плотность вероятности"

	// Добавляем экспериментальную гистограмму
	expPts := make(plotter.XYs, bins)
	for i := 0; i < bins; i++ {
		expPts[i].X = binCenters[i]
		expPts[i].Y = frequencies[i]
	}

	expLine, err := plotter.NewLine(expPts)
	if err != nil {
		return
	}

	expLine.Color = color{255, 0, 0}
	expLine.Width = vg.Points(1)
	p.Add(expLine)
	p.Legend.Add("Экспериментальное", expLine)

	// Добавляем теоретическую кривую
	theoryPts := make(plotter.XYs, 200)
	for i := 0; i < 200; i++ {
		x := float64(i) * maxVal / 199.0
		theoryPts[i].X = x
		theoryPts[i].Y = ExponentialPDF(x, lambda)
	}

	theoryLine, err := plotter.NewLine(theoryPts)
	if err != nil {
		return
	}

	theoryLine.Color = color{0, 0, 255}
	theoryLine.Width = vg.Points(1.5)
	p.Add(theoryLine)
	p.Legend.Add("Теоретическое", theoryLine)
	p.Legend.Top = true

	// Сохраняем график
	filename := fmt.Sprintf("task4_histogram_N%d.png", N)
	p.Save(8*vg.Inch, 6*vg.Inch, filename)
}
