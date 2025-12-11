package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// ========== ОСНОВНЫЕ ФУНКЦИИ ==========

// MultiplicativeRNG - мультипликативный генератор случайных чисел
func MultiplicativeRNG(a, b, m, x0 int64) func() float64 {
	current := x0
	return func() float64 {
		current = (a*current + b) % m
		return float64(current) / float64(m)
	}
}

// UniformDistribution - равномерное распределение в [min, max]
func UniformDistribution(generator func() float64, min, max float64) float64 {
	return min + (max-min)*generator()
}

// ExponentialDistribution - экспоненциальное распределение
func ExponentialDistribution(generator func() float64, lambda float64) float64 {
	u := generator()
	for u == 0 || u == 1 {
		u = generator() // избегаем крайних значений
	}
	return -math.Log(1-u) / lambda
}

// ========== МОДЕЛЬ СИСТЕМЫ МАССОВОГО ОБСЛУЖИВАНИЯ ==========

// Event - событие в системе
type Event struct {
	Time    float64
	Type    string // "arrival" - приход, "departure" - уход
	Request int
}

// QueueingSystem - система массового обслуживания
type QueueingSystem struct {
	ArrivalTimes  []float64   // Времена прихода заявок
	ServiceTimes  []float64   // Времена обслуживания
	BufferSize    int         // Размер буфера
	CurrentTime   float64     // Текущее время моделирования
	ServerBusy    bool        // Занят ли сервер
	Queue         []float64   // Очередь (время прихода каждой заявки)
	DepartureTime float64     // Время завершения текущего обслуживания
	Statistics    *Statistics // Статистика
	TotalRequests int         // Всего поступило заявок
	Processed     int         // Обработано заявок
	Lost          int         // Потеряно заявок (при переполнении буфера)
}

// Statistics - статистика работы системы
type Statistics struct {
	BufferTimes     map[int]float64 // Суммарное время с заданным числом заявок в буфере
	BufferCounts    map[int]int     // Количество случаев
	WaitingTimes    []float64       // Время ожидания для каждой заявки
	QueueLengths    []int           // Длина очереди в моменты событий
	BusyTimes       []float64       // Временные интервалы, когда сервер был занят
	IdleTimes       []float64       // Временные интервалы, когда сервер был свободен
	LastEventTime   float64         // Время последнего события (для накопления статистики)
	LastQueueLength int             // Длина очереди в момент последнего события
	LastServerState bool            // Состояние сервера в момент последнего события
}

// NewQueueingSystem - создание новой системы
func NewQueueingSystem(arrivalTimes, serviceTimes []float64, bufferSize int) *QueueingSystem {
	return &QueueingSystem{
		ArrivalTimes:  arrivalTimes,
		ServiceTimes:  serviceTimes,
		BufferSize:    bufferSize,
		CurrentTime:   0,
		ServerBusy:    false,
		Queue:         make([]float64, 0),
		DepartureTime: math.MaxFloat64,
		Statistics: &Statistics{
			BufferTimes:  make(map[int]float64),
			BufferCounts: make(map[int]int),
			WaitingTimes: make([]float64, 0),
			QueueLengths: make([]int, 0),
			BusyTimes:    make([]float64, 0),
			IdleTimes:    make([]float64, 0),
		},
	}
}

// updateStatistics - обновление статистики
func (qs *QueueingSystem) updateStatistics(eventTime float64) {
	// Обновляем статистику по времени в буфере
	timeDiff := eventTime - qs.Statistics.LastEventTime
	queueLength := len(qs.Queue)
	if qs.ServerBusy {
		queueLength++ // Учитываем заявку на сервере
	}

	qs.Statistics.BufferTimes[queueLength] += timeDiff
	qs.Statistics.BufferCounts[queueLength]++

	// Обновляем статистику загрузки сервера
	if qs.Statistics.LastServerState {
		qs.Statistics.BusyTimes = append(qs.Statistics.BusyTimes, timeDiff)
	} else {
		qs.Statistics.IdleTimes = append(qs.Statistics.IdleTimes, timeDiff)
	}

	qs.Statistics.LastEventTime = eventTime
	qs.Statistics.LastQueueLength = queueLength
	qs.Statistics.LastServerState = qs.ServerBusy
}

// Simulate - основная функция имитационного моделирования (Задание 3)
func (qs *QueueingSystem) Simulate() {
	// Инициализация
	qs.Statistics.LastEventTime = 0
	qs.Statistics.LastQueueLength = 0
	qs.Statistics.LastServerState = false

	arrivalIndex := 0
	serviceIndex := 0

	// Основной цикл имитации
	for arrivalIndex < len(qs.ArrivalTimes) || len(qs.Queue) > 0 || qs.ServerBusy {
		// Определяем время следующего события
		nextArrival := math.MaxFloat64
		if arrivalIndex < len(qs.ArrivalTimes) {
			nextArrival = qs.ArrivalTimes[arrivalIndex]
		}

		nextEventTime := math.Min(nextArrival, qs.DepartureTime)

		// Обновляем статистику до момента события
		qs.updateStatistics(nextEventTime)
		qs.CurrentTime = nextEventTime

		// Обработка завершения обслуживания
		if math.Abs(qs.DepartureTime-qs.CurrentTime) < 1e-9 {
			qs.processDeparture(&serviceIndex)
		}

		// Обработка прихода новой заявки
		if math.Abs(nextArrival-qs.CurrentTime) < 1e-9 && arrivalIndex < len(qs.ArrivalTimes) {
			qs.processArrival(&arrivalIndex, &serviceIndex)
		}
	}
}

// processDeparture - обработка завершения обслуживания
func (qs *QueueingSystem) processDeparture(serviceIndex *int) {
	qs.ServerBusy = false
	qs.DepartureTime = math.MaxFloat64
	qs.Processed++

	// Если есть заявки в очереди, начинаем обслуживать следующую
	if len(qs.Queue) > 0 {
		// Извлекаем заявку из очереди
		arrivalTime := qs.Queue[0]
		qs.Queue = qs.Queue[1:]

		// Регистрируем время ожидания
		waitTime := qs.CurrentTime - arrivalTime
		qs.Statistics.WaitingTimes = append(qs.Statistics.WaitingTimes, waitTime)

		// Начинаем обслуживание
		qs.ServerBusy = true
		if *serviceIndex < len(qs.ServiceTimes) {
			qs.DepartureTime = qs.CurrentTime + qs.ServiceTimes[*serviceIndex]
			(*serviceIndex)++
		}
	}
}

// processArrival - обработка прихода заявки
func (qs *QueueingSystem) processArrival(arrivalIndex, serviceIndex *int) {
	qs.TotalRequests++

	if !qs.ServerBusy {
		// Сервер свободен, начинаем обслуживание сразу
		qs.ServerBusy = true
		if *serviceIndex < len(qs.ServiceTimes) {
			qs.DepartureTime = qs.CurrentTime + qs.ServiceTimes[*serviceIndex]
			(*serviceIndex)++
		}
		// Время ожидания = 0
		qs.Statistics.WaitingTimes = append(qs.Statistics.WaitingTimes, 0)
	} else {
		// Сервер занят, заявка идет в очередь
		if len(qs.Queue) < qs.BufferSize {
			qs.Queue = append(qs.Queue, qs.CurrentTime)
		} else {
			// Буфер переполнен, заявка теряется
			qs.Lost++
		}
	}

	(*arrivalIndex)++
}

// ========== ФУНКЦИИ ДЛЯ ЗАДАНИЙ 4 И 5 ==========

// GetBufferTimes - получение времени нахождения в буфере (Задание 4)
func (qs *QueueingSystem) GetBufferTimes() map[int]float64 {
	return qs.Statistics.BufferTimes
}

// GetBufferProbabilities - расчет вероятностей нахождения в буфере (Задание 5)
func (qs *QueueingSystem) GetBufferProbabilities() map[int]float64 {
	probs := make(map[int]float64)
	totalTime := qs.CurrentTime

	for queueLength, time := range qs.Statistics.BufferTimes {
		probs[queueLength] = time / totalTime
	}

	return probs
}

// GetAverageWaitingTime - среднее время ожидания
func (qs *QueueingSystem) GetAverageWaitingTime() float64 {
	if len(qs.Statistics.WaitingTimes) == 0 {
		return 0
	}

	var sum float64
	for _, wt := range qs.Statistics.WaitingTimes {
		sum += wt
	}
	return sum / float64(len(qs.Statistics.WaitingTimes))
}

// GetServerUtilization - коэффициент использования сервера
func (qs *QueueingSystem) GetServerUtilization() float64 {
	if qs.CurrentTime == 0 {
		return 0
	}

	var busyTime float64
	for _, t := range qs.Statistics.BusyTimes {
		busyTime += t
	}
	return busyTime / qs.CurrentTime
}

func main() {
	fmt.Println("Практическая работа №7")
	fmt.Println("Имитационное моделирование вычислительных систем\n")

	rand.Seed(time.Now().UnixNano())

	// Параметры системы
	numRequests := 100000
	bufferSize := 10

	// ========== ЗАДАНИЕ 1 ==========
	fmt.Println("=== ЗАДАНИЕ 1 ===")
	fmt.Println("Генерация последовательностей случайных чисел (мультипликативный метод)")

	// Параметры генераторов
	M := int64(1000)
	a_TZ := int64(39)
	a_TS := int64(39)
	b := int64(1)
	x0 := int64(1)

	// Создаем генераторы
	randTZ := MultiplicativeRNG(a_TZ, b, M, x0)
	randTS := MultiplicativeRNG(a_TS, b, M, x0)

	// Параметры распределений
	TZmin := 4.0  // сек
	TZmax := 12.0 // сек
	TSmin := 2.0  // сек
	TSmax := 8.0  // сек

	// Генерируем времена между заявками
	interArrivalTimes := make([]float64, numRequests)
	fmt.Printf("Генерация времен между заявками (TZ ∈ [%.1f, %.1f])... ", TZmin, TZmax)
	start := time.Now()

	for i := 0; i < numRequests; i++ {
		interArrivalTimes[i] = UniformDistribution(randTZ, TZmin, TZmax)
	}

	fmt.Printf("завершено за %v\n", time.Since(start))

	// Генерируем времена обслуживания
	serviceTimes := make([]float64, numRequests)
	fmt.Printf("Генерация времен обслуживания (TS ∈ [%.1f, %.1f])... ", TSmin, TSmax)
	start = time.Now()

	for i := 0; i < numRequests; i++ {
		serviceTimes[i] = UniformDistribution(randTS, TSmin, TSmax)
	}

	fmt.Printf("завершено за %v\n", time.Since(start))

	// Базовая статистика
	avgInterArrival := mean(interArrivalTimes)
	avgService := mean(serviceTimes)
	fmt.Printf("\nСтатистика сгенерированных данных:\n")
	fmt.Printf("  Среднее время между заявками: %.4f сек\n", avgInterArrival)
	fmt.Printf("  Среднее время обслуживания: %.4f сек\n", avgService)
	fmt.Printf("  Коэффициент загрузки (ρ = TS/TZ): %.4f\n", avgService/avgInterArrival)

	// ========== ЗАДАНИЕ 2 ==========
	fmt.Println("\n=== ЗАДАНИЕ 2 ===")
	fmt.Println("Определение времен прихода заявок")

	// Вычисляем времена прихода
	arrivalTimes := make([]float64, numRequests)
	currentTime := 0.0

	for i := 0; i < numRequests; i++ {
		currentTime += interArrivalTimes[i]
		arrivalTimes[i] = currentTime
	}

	fmt.Printf("Время прихода последней заявки: %.2f сек\n", arrivalTimes[numRequests-1])
	fmt.Printf("Интервал моделирования: [0, %.2f] сек\n", arrivalTimes[numRequests-1])

	// ========== ЗАДАНИЕ 3, 4, 5 ==========
	fmt.Println("\n=== ЗАДАНИЯ 3, 4, 5 ===")
	fmt.Println("Имитационное моделирование системы (равномерное распределение)")

	// Создаем и запускаем систему
	qsUniform := NewQueueingSystem(arrivalTimes, serviceTimes, bufferSize)

	fmt.Print("Запуск имитации... ")
	start = time.Now()
	qsUniform.Simulate()
	simulationTime := time.Since(start)

	fmt.Printf("завершено за %v\n", simulationTime)

	// Результаты
	fmt.Printf("\nРезультаты имитации (равномерное распределение):\n")
	fmt.Printf("  Всего заявок: %d\n", qsUniform.TotalRequests)
	fmt.Printf("  Обработано заявок: %d\n", qsUniform.Processed)
	fmt.Printf("  Потеряно заявок: %d (%.2f%%)\n",
		qsUniform.Lost, float64(qsUniform.Lost)/float64(qsUniform.TotalRequests)*100)
	fmt.Printf("  Общее время моделирования: %.2f сек\n", qsUniform.CurrentTime)
	fmt.Printf("  Среднее время ожидания: %.4f сек\n", qsUniform.GetAverageWaitingTime())
	fmt.Printf("  Коэффициент использования сервера: %.4f\n", qsUniform.GetServerUtilization())

	// Задание 4: Времена нахождения в буфере
	bufferTimes := qsUniform.GetBufferTimes()
	fmt.Println("\nСреднее время нахождения в буфере (Задание 4):")
	for i := 0; i <= bufferSize; i++ {
		if count, exists := qsUniform.Statistics.BufferCounts[i]; exists && count > 0 {
			avgTime := bufferTimes[i] / float64(count)
			fmt.Printf("  %d заявок в буфере: среднее время = %.4f сек (случаев: %d)\n",
				i, avgTime, count)
		}
	}

	// Задание 5: Вероятности нахождения в буфере
	bufferProbs := qsUniform.GetBufferProbabilities()
	fmt.Println("\nВероятности нахождения в буфере (Задание 5):")
	totalProb := 0.0
	for i := 0; i <= bufferSize; i++ {
		if prob, exists := bufferProbs[i]; exists {
			fmt.Printf("  P(%d) = %.6f\n", i, prob)
			totalProb += prob
		}
	}
	fmt.Printf("Сумма вероятностей: %.6f\n", totalProb)

	// ========== ЗАДАНИЕ 6 ==========
	fmt.Println("\n=== ЗАДАНИЕ 6 ===")
	fmt.Println("Генерация экспоненциально распределенных последовательностей")

	// Параметры экспоненциальных распределений
	lambda := 1.0 / 3.0 // Для входного потока
	mu := 1.0 / 4.0     // Для времени обслуживания

	// Генерируем экспоненциально распределенные времена
	expInterArrivalTimes := make([]float64, numRequests)
	expServiceTimes := make([]float64, numRequests)

	fmt.Print("Генерация экспоненциальных времен между заявками... ")
	start = time.Now()

	for i := 0; i < numRequests; i++ {
		expInterArrivalTimes[i] = ExponentialDistribution(rand.Float64, lambda)
	}

	fmt.Printf("завершено за %v\n", time.Since(start))

	fmt.Print("Генерация экспоненциальных времен обслуживания... ")
	start = time.Now()

	for i := 0; i < numRequests; i++ {
		expServiceTimes[i] = ExponentialDistribution(rand.Float64, mu)
	}

	fmt.Printf("завершено за %v\n", time.Since(start))

	// Вычисляем времена прихода
	expArrivalTimes := make([]float64, numRequests)
	currentTime = 0.0
	for i := 0; i < numRequests; i++ {
		currentTime += expInterArrivalTimes[i]
		expArrivalTimes[i] = currentTime
	}

	// Статистика
	avgExpInterArrival := mean(expInterArrivalTimes)
	avgExpService := mean(expServiceTimes)
	fmt.Printf("\nСтатистика экспоненциальных данных:\n")
	fmt.Printf("  Среднее время между заявками: %.4f сек (теоретическое: %.4f)\n",
		avgExpInterArrival, 1/lambda)
	fmt.Printf("  Среднее время обслуживания: %.4f сек (теоретическое: %.4f)\n",
		avgExpService, 1/mu)
	fmt.Printf("  Коэффициент загрузки (ρ = TS/TZ): %.4f\n", avgExpService/avgExpInterArrival)

	// ========== ЗАДАНИЕ 7 ==========
	fmt.Println("\n=== ЗАДАНИЕ 7 ===")
	fmt.Println("Имитационное моделирование системы (экспоненциальное распределение)")

	// Создаем и запускаем систему с экспоненциальными распределениями
	qsExponential := NewQueueingSystem(expArrivalTimes, expServiceTimes, bufferSize)

	fmt.Print("Запуск имитации... ")
	start = time.Now()
	qsExponential.Simulate()
	simulationTime = time.Since(start)

	fmt.Printf("завершено за %v\n", simulationTime)

	// Результаты
	fmt.Printf("\nРезультаты имитации (экспоненциальное распределение):\n")
	fmt.Printf("  Всего заявок: %d\n", qsExponential.TotalRequests)
	fmt.Printf("  Обработано заявок: %d\n", qsExponential.Processed)
	fmt.Printf("  Потеряно заявок: %d (%.2f%%)\n",
		qsExponential.Lost, float64(qsExponential.Lost)/float64(qsExponential.TotalRequests)*100)
	fmt.Printf("  Общее время моделирования: %.2f сек\n", qsExponential.CurrentTime)
	fmt.Printf("  Среднее время ожидания: %.4f сек\n", qsExponential.GetAverageWaitingTime())
	fmt.Printf("  Коэффициент использования сервера: %.4f\n", qsExponential.GetServerUtilization())

	// Времена нахождения в буфере
	expBufferTimes := qsExponential.GetBufferTimes()
	fmt.Println("\nСреднее время нахождения в буфере (экспоненциальное):")
	for i := 0; i <= bufferSize; i++ {
		if count, exists := qsExponential.Statistics.BufferCounts[i]; exists && count > 0 {
			avgTime := expBufferTimes[i] / float64(count)
			fmt.Printf("  %d заявок в буфере: среднее время = %.4f сек (случаев: %d)\n",
				i, avgTime, count)
		}
	}

	// Вероятности нахождения в буфере
	expBufferProbs := qsExponential.GetBufferProbabilities()
	fmt.Println("\nВероятности нахождения в буфере (экспоненциальное):")
	totalProb = 0.0
	for i := 0; i <= bufferSize; i++ {
		if prob, exists := expBufferProbs[i]; exists {
			fmt.Printf("  P(%d) = %.6f\n", i, prob)
			totalProb += prob
		}
	}
	fmt.Printf("Сумма вероятностей: %.6f\n", totalProb)

	// ========== СРАВНИТЕЛЬНЫЙ АНАЛИЗ ==========
	fmt.Println("\n=== СРАВНИТЕЛЬНЫЙ АНАЛИЗ ===")
	fmt.Println("Сравнение равномерного и экспоненциального распределений:")

	// Ключевые метрики
	metrics := []struct {
		name        string
		uniform     float64
		exponential float64
	}{
		{"Среднее время ожидания (сек)",
			qsUniform.GetAverageWaitingTime(), qsExponential.GetAverageWaitingTime()},
		{"Коэффициент использования сервера",
			qsUniform.GetServerUtilization(), qsExponential.GetServerUtilization()},
		{"Вероятность потери заявок (%)",
			float64(qsUniform.Lost) / float64(qsUniform.TotalRequests) * 100,
			float64(qsExponential.Lost) / float64(qsExponential.TotalRequests) * 100},
		{"Средняя длина очереди", calculateAverageQueueLength(qsUniform),
			calculateAverageQueueLength(qsExponential)},
	}

	for _, metric := range metrics {
		fmt.Printf("  %s: равномерное=%.4f, экспоненциальное=%.4f\n",
			metric.name, metric.uniform, metric.exponential)
	}

	// Теоретические значения для M/M/1 (экспоненциальное распределение)
	rho := avgExpService / avgExpInterArrival
	if rho < 1 {
		fmt.Println("\nТеоретические значения для M/M/1 системы (экспоненциальное):")
		theoryAvgQueue := rho * rho / (1 - rho)
		theoryAvgWait := rho / (mu * (1 - rho))
		theoryIdleProb := 1 - rho

		fmt.Printf("  Теоретическая средняя длина очереди: %.4f\n", theoryAvgQueue)
		fmt.Printf("  Теоретическое среднее время ожидания: %.4f сек\n", theoryAvgWait)
		fmt.Printf("  Теоретическая вероятность простоя сервера: %.4f\n", theoryIdleProb)

		// Сравнение с экспериментальными значениями
		expAvgQueue := calculateAverageQueueLength(qsExponential)
		expAvgWait := qsExponential.GetAverageWaitingTime()
		expIdleProb := 1 - qsExponential.GetServerUtilization()

		fmt.Printf("\n  Сравнение с экспериментальными значениями:\n")
		fmt.Printf("    Длина очереди: ошибка = %.2f%%\n",
			math.Abs(expAvgQueue-theoryAvgQueue)/theoryAvgQueue*100)
		fmt.Printf("    Время ожидания: ошибка = %.2f%%\n",
			math.Abs(expAvgWait-theoryAvgWait)/theoryAvgWait*100)
		fmt.Printf("    Вероятность простоя: ошибка = %.2f%%\n",
			math.Abs(expIdleProb-theoryIdleProb)/theoryIdleProb*100)
	}

	fmt.Println("\n=== ПРАКТИЧЕСКАЯ РАБОТА ЗАВЕРШЕНА ===")
}

// ========== ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ ==========

// mean - вычисление среднего значения
func mean(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}

	var sum float64
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

// calculateAverageQueueLength - вычисление средней длины очереди
func calculateAverageQueueLength(qs *QueueingSystem) float64 {
	var sum float64
	var count float64

	for queueLength, time := range qs.Statistics.BufferTimes {
		sum += float64(queueLength) * time
		count += time
	}

	if count == 0 {
		return 0
	}
	return sum / count
}
