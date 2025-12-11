package main

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"math"
	"math/rand"
	"time"
)

// Задание 1: Функция плотности вероятности нормального распределения
func NormalPDF(x, mean, sigma float64) float64 {
	if sigma <= 0 {
		return 0
	}
	exponent := -0.5 * math.Pow((x-mean)/sigma, 2)
	denominator := sigma * math.Sqrt(2*math.Pi)
	return math.Exp(exponent) / denominator
}

// Задание 2: Функция распределения для нормального закона (аппроксимация)
func NormalCDF(x, mean, sigma float64) float64 {
	if sigma <= 0 {
		if x < mean {
			return 0
		}
		return 1
	}

	// Используем аппроксимацию функции ошибок
	z := (x - mean) / (sigma * math.Sqrt2)

	// Функция ошибок
	t := 1.0 / (1.0 + 0.3275911*math.Abs(z))
	t2 := t * t
	t3 := t2 * t
	t4 := t3 * t
	t5 := t4 * t

	a1 := 0.254829592
	a2 := -0.284496736
	a3 := 1.421413741
	a4 := -1.453152027
	a5 := 1.061405429

	erf := 1.0 - (a1*t+a2*t2+a3*t3+a4*t4+a5*t5)*math.Exp(-z*z)

	if z < 0 {
		erf = -erf
	}

	return 0.5 * (1 + erf)
}

// Задание 3: Моделирование нормального распределения методом обратной функции
// с кусочно-линейной аппроксимацией
func GenerateNormalInverseCDF(mean, sigma, a, b float64, intervals, experiments int) []float64 {
	// Создаем массивы для кусочно-линейной аппроксимации обратной функции
	probabilities := make([]float64, intervals+1)
	values := make([]float64, intervals+1)

	// Заполняем массивы: F(x) = p, x = F^{-1}(p)
	dp := 1.0 / float64(intervals)
	for i := 0; i <= intervals; i++ {
		p := float64(i) * dp

		// Для нахождения x: F(x) = p, используем численный метод
		// Начинаем поиск от a и идем до b
		dx := (b - a) / 1000.0
		bestX := a
		minDiff := math.Abs(NormalCDF(a, mean, sigma) - p)

		for j := 0; j <= 1000; j++ {
			currentX := a + float64(j)*dx
			currentP := NormalCDF(currentX, mean, sigma)
			diff := math.Abs(currentP - p)

			if diff < minDiff {
				minDiff = diff
				bestX = currentX
			}
		}

		probabilities[i] = p
		values[i] = bestX
	}

	// Генерируем случайные числа с нормальным распределением
	randNumbers := make([]float64, experiments)

	for i := 0; i < experiments; i++ {
		// Генерируем равномерно распределенное число Y ∈ (0,1)
		Y := rand.Float64()

		// Определяем номер интервала
		j := int(math.Floor(Y * float64(intervals)))
		if j >= intervals {
			j = intervals - 1
		}

		// Линейная интерполяция внутри интервала
		p0 := float64(j) * dp
		p1 := float64(j+1) * dp

		// Интерполяция для нахождения x
		x0 := values[j]
		x1 := values[j+1]

		// Линейная интерполяция: x = x0 + (Y-p0)*(x1-x0)/(p1-p0)
		randNumbers[i] = x0 + (Y-p0)*(x1-x0)/(p1-p0)
	}

	return randNumbers
}

// Функция для построения гистограммы
func BuildHistogram(data []float64, bins int, a, b float64, filename, title string) {
	// Создаем гистограмму
	h, _ := plotter.NewHist(plotter.Values(data), bins)
	h.Normalize(1) // Нормализуем для получения плотности вероятности

	// Создаем график
	p := plot.New()
	p.Title.Text = title
	p.X.Label.Text = "Значение"
	p.Y.Label.Text = "Плотность вероятности"

	p.Add(h)

	// Устанавливаем пределы по X
	p.X.Min = a
	p.X.Max = b

	// Сохраняем в файл
	if err := p.Save(6*vg.Inch, 4*vg.Inch, filename); err != nil {
		fmt.Printf("Ошибка сохранения гистограммы: %v\n", err)
	} else {
		fmt.Printf("Гистограмма сохранена: %s\n", filename)
	}
}

// Функция для расчета экспериментальной плотности вероятности
func CalculateExperimentalPDF(data []float64, bins int, a, b float64) []float64 {
	// Инициализируем массив для подсчета
	counts := make([]int, bins)
	binWidth := (b - a) / float64(bins)

	// Подсчитываем попадания в бины
	for _, value := range data {
		if value >= a && value <= b {
			binIndex := int((value - a) / binWidth)
			if binIndex >= bins {
				binIndex = bins - 1
			}
			counts[binIndex]++
		}
	}

	// Преобразуем в плотность вероятности
	pdf := make([]float64, bins)
	total := float64(len(data))

	for i := 0; i < bins; i++ {
		pdf[i] = float64(counts[i]) / (total * binWidth)
	}

	return pdf
}

// Функция для расчета теоретической плотности вероятности в центрах бинов
func CalculateTheoreticalPDF(bins int, a, b, mean, sigma float64) []float64 {
	pdf := make([]float64, bins)
	binWidth := (b - a) / float64(bins)

	for i := 0; i < bins; i++ {
		x := a + (float64(i)+0.5)*binWidth
		pdf[i] = NormalPDF(x, mean, sigma)
	}

	return pdf
}

// Функция для расчета средней квадратичной погрешности (RMSE)
func CalculateRMSE(experimental, theoretical []float64) float64 {
	if len(experimental) != len(theoretical) {
		return 0
	}

	var sum float64
	n := len(experimental)

	for i := 0; i < n; i++ {
		diff := experimental[i] - theoretical[i]
		sum += diff * diff
	}

	return math.Sqrt(sum / float64(n))
}

func main() {
	// Инициализируем генератор случайных чисел
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Практическая работа №3")
	fmt.Println("Моделирование нормального закона распределения\n")

	// Задание 1: Построение графиков плотности вероятности
	fmt.Println("=== ЗАДАНИЕ 1 ===")
	fmt.Println("Графики плотности вероятности нормального распределения:")

	// Параметры из задания
	params := []struct {
		mean  float64
		sigma float64
		label string
		color color
	}{
		{10, 2, "M=10, σ=2", color{0, 0, 255}},     // синий
		{10, 1, "M=10, σ=1", color{255, 0, 0}},     // красный
		{10, 0.5, "M=10, σ=1/2", color{0, 128, 0}}, // зеленый
		{12, 1, "M=12, σ=1", color{255, 0, 255}},   // фиолетовый
	}

	// Создаем график
	p1 := plot.New()
	p1.Title.Text = "Плотность вероятности нормального распределения"
	p1.X.Label.Text = "x"
	p1.Y.Label.Text = "f(x)"
	p1.Legend.Top = true

	// Диапазон для построения графиков
	xMin := 5.0
	xMax := 17.0

	// Добавляем кривые для каждого набора параметров
	for _, param := range params {
		pts := make(plotter.XYs, 200)
		for i := range pts {
			x := xMin + (xMax-xMin)*float64(i)/199.0
			pts[i].X = x
			pts[i].Y = NormalPDF(x, param.mean, param.sigma)
		}

		line, err := plotter.NewLine(pts)
		if err != nil {
			fmt.Printf("Ошибка создания линии: %v\n", err)
			continue
		}

		// Задаем цвет линии
		line.Color = param.color
		line.Width = vg.Points(1.5)

		p1.Add(line)
		p1.Legend.Add(param.label, line)
	}

	// Сохраняем график
	if err := p1.Save(10*vg.Inch, 6*vg.Inch, "task1_normal_pdf.png"); err != nil {
		fmt.Printf("Ошибка сохранения графика: %v\n", err)
	} else {
		fmt.Println("График плотности вероятности сохранен: task1_normal_pdf.png")
	}

	// Задание 2: Функция распределения
	fmt.Println("\n=== ЗАДАНИЕ 2 ===")
	fmt.Println("Функция распределения для M=10, σ=2")

	mean := 10.0
	sigma := 2.0
	a := 0.0
	b := 20.0

	// Создаем график функции распределения
	p2 := plot.New()
	p2.Title.Text = "Функция распределения F(x) для M=10, σ=2"
	p2.X.Label.Text = "x"
	p2.Y.Label.Text = "F(x)"

	ptsCDF := make(plotter.XYs, 200)
	for i := range ptsCDF {
		x := a + (b-a)*float64(i)/199.0
		ptsCDF[i].X = x
		ptsCDF[i].Y = NormalCDF(x, mean, sigma)
	}

	lineCDF, err := plotter.NewLine(ptsCDF)
	if err != nil {
		fmt.Printf("Ошибка создания линии: %v\n", err)
	} else {
		lineCDF.Color = color{0, 0, 255}
		lineCDF.Width = vg.Points(2)
		p2.Add(lineCDF)
	}

	// Сохраняем график
	if err := p2.Save(10*vg.Inch, 6*vg.Inch, "task2_cdf_function.png"); err != nil {
		fmt.Printf("Ошибка сохранения графика: %v\n", err)
	} else {
		fmt.Println("График функции распределения сохранен: task2_cdf_function.png")
	}

	// Задание 3: Моделирование нормального распределения
	fmt.Println("\n=== ЗАДАНИЕ 3 ===")
	fmt.Println("Моделирование методом обратной функции с кусочно-линейной аппроксимацией")

	intervals := 100
	experimentCounts := []int{1000, 10000, 100000, 1000000}

	// Массивы для хранения сгенерированных данных
	generatedData := make([][]float64, len(experimentCounts))

	for idx, N := range experimentCounts {
		fmt.Printf("\nГенерация N=%d... ", N)
		start := time.Now()

		data := GenerateNormalInverseCDF(mean, sigma, a, b, intervals, N)
		generatedData[idx] = data

		// Вычисляем статистики
		var sum, sumSq float64
		for _, val := range data {
			sum += val
			sumSq += val * val
		}

		dataMean := sum / float64(N)
		dataVariance := sumSq/float64(N) - dataMean*dataMean
		dataSigma := math.Sqrt(math.Abs(dataVariance))

		elapsed := time.Since(start)

		fmt.Printf("завершено за %v\n", elapsed)
		fmt.Printf("  Выборочное среднее: %.4f (теоретическое: %.2f)\n", dataMean, mean)
		fmt.Printf("  Выборочное СКО: %.4f (теоретическое: %.2f)\n", dataSigma, sigma)
	}

	// Задание 4: Построение гистограмм
	fmt.Println("\n=== ЗАДАНИЕ 4 ===")
	fmt.Println("Гистограммы относительных частот на 100 интервалах")

	histogramExperiments := []int{100, 1000, 10000, 100000}

	// Генерируем дополнительные данные для N=100
	smallData := GenerateNormalInverseCDF(mean, sigma, a, b, intervals, 100)

	// Строим гистограммы
	for _, N := range histogramExperiments {
		var data []float64

		// Выбираем данные соответствующего размера
		switch N {
		case 100:
			data = smallData
		case 1000:
			if len(generatedData[0]) >= 1000 {
				data = generatedData[0][:1000]
			} else {
				data = GenerateNormalInverseCDF(mean, sigma, a, b, intervals, 1000)
			}
		case 10000:
			if len(generatedData[1]) >= 10000 {
				data = generatedData[1][:10000]
			} else {
				data = GenerateNormalInverseCDF(mean, sigma, a, b, intervals, 10000)
			}
		case 100000:
			if len(generatedData[2]) >= 100000 {
				data = generatedData[2][:100000]
			} else {
				data = GenerateNormalInverseCDF(mean, sigma, a, b, intervals, 100000)
			}
		}

		filename := fmt.Sprintf("task4_histogram_N%d.png", N)
		title := fmt.Sprintf("Нормальное распределение, N=%d", N)

		BuildHistogram(data, 100, a, b, filename, title)

		// Выводим дополнительную статистику
		var sum float64
		for _, val := range data {
			sum += val
		}
		dataMean := sum / float64(len(data))
		fmt.Printf("  N=%6d: среднее = %.4f\n", N, dataMean)
	}

	// Задание 5: Расчет средней квадратичной погрешности
	fmt.Println("\n=== ЗАДАНИЕ 5 ===")
	fmt.Println("Расчет RMSE между экспериментальным и теоретическим распределениями")

	// Количество бинов для гистограммы
	bins := 100

	// Вычисляем теоретическую плотность вероятности
	theoreticalPDF := CalculateTheoreticalPDF(bins, a, b, mean, sigma)

	// Массивы для хранения RMSE
	rmseValues := make([]float64, len(histogramExperiments))

	// Создаем график для зависимости RMSE от N
	p5 := plot.New()
	p5.Title.Text = "Зависимость RMSE от числа экспериментов"
	p5.X.Label.Text = "Число экспериментов (N)"
	p5.Y.Label.Text = "RMSE"
	p5.X.Scale = plot.LogScale{}
	p5.Y.Scale = plot.LogScale{}

	// Создаем точки для графика
	rmsePoints := make(plotter.XYs, len(histogramExperiments))

	fmt.Println("\nРасчет RMSE:")
	for i, N := range histogramExperiments {
		var data []float64

		// Выбираем данные соответствующего размера
		switch N {
		case 100:
			data = smallData
		case 1000:
			data = GenerateNormalInverseCDF(mean, sigma, a, b, intervals, 1000)
		case 10000:
			data = GenerateNormalInverseCDF(mean, sigma, a, b, intervals, 10000)
		case 100000:
			data = GenerateNormalInverseCDF(mean, sigma, a, b, intervals, 100000)
		}

		// Вычисляем экспериментальную плотность вероятности
		experimentalPDF := CalculateExperimentalPDF(data, bins, a, b)

		// Вычисляем RMSE
		rmse := CalculateRMSE(experimentalPDF, theoreticalPDF)
		rmseValues[i] = rmse
		rmsePoints[i].X = float64(N)
		rmsePoints[i].Y = rmse

		fmt.Printf("  N=%6d: RMSE = %.6f\n", N, rmse)
	}

	// Добавляем точки на график
	scatter, err := plotter.NewScatter(rmsePoints)
	if err != nil {
		fmt.Printf("Ошибка создания scatter plot: %v\n", err)
	} else {
		scatter.GlyphStyle.Color = color{255, 0, 0}
		scatter.GlyphStyle.Radius = vg.Points(3)
		p5.Add(scatter)

		// Добавляем линию тренда
		line, err := plotter.NewLine(rmsePoints)
		if err != nil {
			fmt.Printf("Ошибка создания линии: %v\n", err)
		} else {
			line.Color = color{0, 0, 255}
			line.Width = vg.Points(1)
			p5.Add(line)
		}
	}

	// Сохраняем график
	if err := p5.Save(10*vg.Inch, 6*vg.Inch, "task5_rmse_vs_N.png"); err != nil {
		fmt.Printf("Ошибка сохранения графика: %v\n", err)
	} else {
		fmt.Println("\nГрафик зависимости RMSE от N сохранен: task5_rmse_vs_N.png")
	}

	// Дополнительный анализ: сравнение для N=100000
	fmt.Println("\n=== ДОПОЛНИТЕЛЬНЫЙ АНАЛИЗ ===")

	// Берем данные для N=100000
	N := 100000
	data := GenerateNormalInverseCDF(mean, sigma, a, b, intervals, N)

	// Создаем график сравнения теоретического и экспериментального распределений
	p6 := plot.New()
	p6.Title.Text = fmt.Sprintf("Сравнение распределений (N=%d)", N)
	p6.X.Label.Text = "x"
	p6.Y.Label.Text = "Плотность вероятности"
	p6.Legend.Top = true

	// Теоретическое распределение (гладкая кривая)
	theoryPts := make(plotter.XYs, 200)
	for i := range theoryPts {
		x := a + (b-a)*float64(i)/199.0
		theoryPts[i].X = x
		theoryPts[i].Y = NormalPDF(x, mean, sigma)
	}

	theoryLine, err := plotter.NewLine(theoryPts)
	if err != nil {
		fmt.Printf("Ошибка создания линии: %v\n", err)
	} else {
		theoryLine.Color = color{0, 0, 255}
		theoryLine.Width = vg.Points(2)
		p6.Add(theoryLine)
		p6.Legend.Add("Теоретическое", theoryLine)
	}

	// Экспериментальное распределение (гистограмма)
	expPDF := CalculateExperimentalPDF(data, 50, a, b)
	expPts := make(plotter.XYs, 50)
	binWidth := (b - a) / 50.0

	for i := 0; i < 50; i++ {
		expPts[i].X = a + (float64(i)+0.5)*binWidth
		expPts[i].Y = expPDF[i]
	}

	expLine, err := plotter.NewLine(expPts)
	if err != nil {
		fmt.Printf("Ошибка создания линии: %v\n", err)
	} else {
		expLine.Color = color{255, 0, 0}
		expLine.Width = vg.Points(1)
		expLine.Dashes = []vg.Length{vg.Points(2), vg.Points(2)}
		p6.Add(expLine)
		p6.Legend.Add("Экспериментальное", expLine)
	}

	// Сохраняем график сравнения
	if err := p6.Save(10*vg.Inch, 6*vg.Inch, "comparison_theory_vs_exp.png"); err != nil {
		fmt.Printf("Ошибка сохранения графика: %v\n", err)
	} else {
		fmt.Println("График сравнения сохранен: comparison_theory_vs_exp.png")
	}

	fmt.Println("\n=== ПРАКТИЧЕСКАЯ РАБОТА ЗАВЕРШЕНА ===")
	fmt.Println("Созданы файлы:")
	fmt.Println("1. task1_normal_pdf.png - плотности вероятности для разных параметров")
	fmt.Println("2. task2_cdf_function.png - функция распределения")
	fmt.Println("3. task4_histogram_N*.png - гистограммы для разных N")
	fmt.Println("4. task5_rmse_vs_N.png - зависимость RMSE от N")
	fmt.Println("5. comparison_theory_vs_exp.png - сравнение теоретического и экспериментального")
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
