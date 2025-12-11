package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Задание 1: Мультипликативный генератор случайных чисел
func MultiplicativeRNG(a, b, m, x0 int64) func() float64 {
	current := x0
	return func() float64 {
		current = (a*current + b) % m
		return float64(current) / float64(m)
	}
}

// Функция для генерации равномерно распределенных чисел в интервале [min, max]
func UniformDistribution(generator func() float64, min, max float64) float64 {
	return min + (max-min)*generator()
}

// Функция для генерации экспоненциально распределенных чисел
func ExponentialDistribution(generator func() float64, lambda float64) float64 {
	return -math.Log(1-generator()) / lambda
}

// Структура для события (приход или завершение обработки)
type Event struct {
	Time    float64
	Type    string // "arrival" или "departure"
	Request int
}

// Структура для системы массового обслуживания
type QueueingSystem struct {
	ArrivalTimes  []float64
	ServiceTimes  []float64
	BufferSize    int
	CurrentTime   float64
	ServerBusy    bool
	Queue         []float64       // Время прихода заявок в буфер
	DepartureTime float64         // Время завершения текущей обработки
	Stats         map[int]float64 // Статистика по времени нахождения в буфере
	Counts        map[int]int     // Количество случаев нахождения в буфере
	TotalRequests int
	Processed     int
}

// Конструктор системы
func NewQueueingSystem(arrivalTimes, serviceTimes []float64, bufferSize int) *QueueingSystem {
	return &QueueingSystem{
		ArrivalTimes:  arrivalTimes,
		ServiceTimes:  serviceTimes,
		BufferSize:    bufferSize,
		Stats:         make(map[int]float64),
		Counts:        make(map[int]int),
		CurrentTime:   0,
		ServerBusy:    false,
		Queue:         make([]float64, 0),
		DepartureTime: math.MaxFloat64,
	}
}

// Основная функция имитационного моделирования (Задание 3)
func (qs *QueueingSystem) Simulate() {
	arrivalIndex := 0
	serviceIndex := 0

	// Обрабатываем все заявки
	for arrivalIndex < len(qs.ArrivalTimes) || len(qs.Queue) > 0 || qs.ServerBusy {
		// Определяем следующее событие
		nextArrival := math.MaxFloat64
		if arrivalIndex < len(qs.ArrivalTimes) {
			nextArrival = qs.ArrivalTimes[arrivalIndex]
		}

		nextEventTime := math.Min(nextArrival, qs.DepartureTime)

		// Обновляем текущее время
		qs.CurrentTime = nextEventTime

		// Обрабатываем события, которые происходят в это время
		// 1. Завершение обработки
		if math.Abs(qs.DepartureTime-qs.CurrentTime) < 1e-9 {
			qs.ServerBusy = false
			qs.DepartureTime = math.MaxFloat64
			qs.Processed++

			// Если есть заявки в буфере, начинаем обслуживать следующую
			if len(qs.Queue) > 0 {
				// Достаем заявку из буфера
				arrivalTime := qs.Queue[0]
				qs.Queue = qs.Queue[1:]

				// Начинаем обслуживание
				qs.ServerBusy = true
				if serviceIndex < len(qs.ServiceTimes) {
					qs.DepartureTime = qs.CurrentTime + qs.ServiceTimes[serviceIndex]
					serviceIndex++
				}

				// Регистрируем время нахождения в буфере
				waitTime := qs.CurrentTime - arrivalTime
				queueLength := len(qs.Queue) + 1 // +1 для заявки, которая только что вышла из буфера
				qs.Stats[queueLength] += waitTime
				qs.Counts[queueLength]++
			}
		}

		// 2. Приход новой заявки
		if math.Abs(nextArrival-qs.CurrentTime) < 1e-9 && arrivalIndex < len(qs.ArrivalTimes) {
			qs.TotalRequests++

			if !qs.ServerBusy {
				// Сервер свободен, начинаем обслуживание сразу
				qs.ServerBusy = true
				if serviceIndex < len(qs.ServiceTimes) {
					qs.DepartureTime = qs.CurrentTime + qs.ServiceTimes[serviceIndex]
					serviceIndex++
				}
				// Заявка не попадает в буфер, время ожидания = 0
			} else {
				// Сервер занят, заявка идет в буфер
				if len(qs.Queue) < qs.BufferSize {
					qs.Queue = append(qs.Queue, qs.CurrentTime)
				} else {
					// Буфер полон, заявка теряется
					// В данной модели просто пропускаем
				}
			}

			arrivalIndex++
		}
	}
}

// Задание 4: Получение времени нахождения в буфере
func (qs *QueueingSystem) GetBufferTimes() map[int]float64 {
	return qs.Stats
}

// Задание 5: Расчет вероятностей нахождения в буфере
func (qs *QueueingSystem) GetBufferProbabilities() map[int]float64 {
	probs := make(map[int]float64)
	totalTime := qs.CurrentTime

	for queueLength, time := range qs.Stats {
		probs[queueLength] = time / totalTime
	}

	return probs
}

func main() {
	fmt.Println("Практическая работа №4")
	fmt.Println("Имитационное моделирование вычислительных систем\n")

	rand.Seed(time.Now().UnixNano())

	// Параметры системы
	numRequests := 10000
	bufferSize := 10

	// ========== ЗАДАНИЕ 1 ==========
	fmt.Println("=== ЗАДАНИЕ 1 ===")
	fmt.Println("Генерация последовательностей случайных чисел с мультипликативным методом")

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
	for i := 0; i < numRequests; i++ {
		interArrivalTimes[i] = UniformDistribution(randTZ, TZmin, TZmax)
	}

	// Генерируем времена обработки
	serviceTimes := make([]float64, numRequests)
	for i := 0; i < numRequests; i++ {
		serviceTimes[i] = UniformDistribution(randTS, TSmin, TSmax)
	}

	fmt.Printf("Сгенерировано %d времен между заявками (TZ ∈ [%.1f, %.1f])\n",
		numRequests, TZmin, TZmax)
	fmt.Printf("Сгенерировано %d времен обработки (TS ∈ [%.1f, %.1f])\n\n",
		numRequests, TSmin, TSmax)

	// ========== ЗАДАНИЕ 2 ==========
	fmt.Println("=== ЗАДАНИЕ 2 ===")
	fmt.Println("Определение времен прихода заявок")

	// Вычисляем времена прихода
	arrivalTimes := make([]float64, numRequests)
	currentTime := 0.0
	for i := 0; i < numRequests; i++ {
		currentTime += interArrivalTimes[i]
		arrivalTimes[i] = currentTime
	}

	fmt.Printf("Первые 5 времен прихода: ")
	for i := 0; i < 5 && i < numRequests; i++ {
		fmt.Printf("%.2f ", arrivalTimes[i])
	}
	fmt.Println("...")
	fmt.Printf("Последние 5 времен прихода: ")
	for i := numRequests - 5; i < numRequests && i >= 0; i++ {
		fmt.Printf("%.2f ", arrivalTimes[i])
	}
	fmt.Println("\n")

	// ========== ЗАДАНИЕ 3 ==========
	fmt.Println("=== ЗАДАНИЕ 3 ===")
	fmt.Println("Разработка программы имитационного моделирования")

	// Создаем систему
	qs := NewQueueingSystem(arrivalTimes, serviceTimes, bufferSize)

	// ========== ЗАДАНИЕ 4 ==========
	fmt.Println("=== ЗАДАНИЕ 4 ===")
	fmt.Println("Определение времени нахождения в буфере")

	// Запускаем имитацию
	startTime := time.Now()
	qs.Simulate()
	simulationTime := time.Since(startTime)

	fmt.Printf("Имитация завершена за %v\n", simulationTime)
	fmt.Printf("Всего заявок: %d, Обработано: %d\n", qs.TotalRequests, qs.Processed)
	fmt.Printf("Общее время моделирования: %.2f сек\n", qs.CurrentTime)

	// Получаем времена нахождения в буфере
	bufferTimes := qs.GetBufferTimes()
	fmt.Println("\nСреднее время нахождения в буфере:")
	for i := 0; i <= bufferSize; i++ {
		if count, exists := qs.Counts[i]; exists && count > 0 {
			avgTime := bufferTimes[i] / float64(count)
			fmt.Printf("  %d заявок в буфере: среднее время = %.4f сек (случаев: %d)\n",
				i, avgTime, count)
		}
	}

	// ========== ЗАДАНИЕ 5 ==========
	fmt.Println("\n=== ЗАДАНИЕ 5 ===")
	fmt.Println("Вероятности нахождения в буфере")

	bufferProbs := qs.GetBufferProbabilities()
	fmt.Println("Вероятности состояний буфера:")
	totalProb := 0.0
	for i := 0; i <= bufferSize; i++ {
		if prob, exists := bufferProbs[i]; exists {
			fmt.Printf("  P(%d) = %.6f\n", i, prob)
			totalProb += prob
		}
	}
	fmt.Printf("Сумма вероятностей: %.6f\n\n", totalProb)

	// ========== ЗАДАНИЕ 6 ==========
	fmt.Println("=== ЗАДАНИЕ 6 ===")
	fmt.Println("Генерация экспоненциально распределенных последовательностей")

	// Параметры экспоненциальных распределений
	lambda := 1.0 / 3.0 // Для входного потока
	mu := 1.0 / 4.0     // Для времени обработки

	// Генерируем экспоненциально распределенные времена
	expInterArrivalTimes := make([]float64, numRequests)
	expServiceTimes := make([]float64, numRequests)

	for i := 0; i < numRequests; i++ {
		expInterArrivalTimes[i] = ExponentialDistribution(rand.Float64, lambda)
		expServiceTimes[i] = ExponentialDistribution(rand.Float64, mu)
	}

	// Вычисляем времена прихода для экспоненциального распределения
	expArrivalTimes := make([]float64, numRequests)
	currentTime = 0.0
	for i := 0; i < numRequests; i++ {
		currentTime += expInterArrivalTimes[i]
		expArrivalTimes[i] = currentTime
	}

	fmt.Printf("Сгенерировано экспоненциальных времен между заявками (λ=%.3f)\n", lambda)
	fmt.Printf("Сгенерировано экспоненциальных времен обработки (μ=%.3f)\n", mu)
	fmt.Printf("Среднее время между заявками: %.4f сек (теоретическое: %.4f)\n",
		mean(expInterArrivalTimes), 1/lambda)
	fmt.Printf("Среднее время обработки: %.4f сек (теоретическое: %.4f)\n\n",
		mean(expServiceTimes), 1/mu)

	// ========== ЗАДАНИЕ 7 ==========
	fmt.Println("=== ЗАДАНИЕ 7 ===")
	fmt.Println("Анализ для экспоненциальных законов распределения")

	// Создаем новую систему с экспоненциальными распределениями
	expQS := NewQueueingSystem(expArrivalTimes, expServiceTimes, bufferSize)

	// Запускаем имитацию
	startTime = time.Now()
	expQS.Simulate()
	simulationTime = time.Since(startTime)

	fmt.Printf("Имитация завершена за %v\n", simulationTime)
	fmt.Printf("Всего заявок: %d, Обработано: %d\n", expQS.TotalRequests, expQS.Processed)
	fmt.Printf("Общее время моделирования: %.2f сек\n", expQS.CurrentTime)

	// Времена нахождения в буфере
	expBufferTimes := expQS.GetBufferTimes()
	fmt.Println("\nСреднее время нахождения в буфере (экспоненциальное):")
	for i := 0; i <= bufferSize; i++ {
		if count, exists := expQS.Counts[i]; exists && count > 0 {
			avgTime := expBufferTimes[i] / float64(count)
			fmt.Printf("  %d заявок в буфере: среднее время = %.4f сек (случаев: %d)\n",
				i, avgTime, count)
		}
	}

	// Вероятности нахождения в буфере
	expBufferProbs := expQS.GetBufferProbabilities()
	fmt.Println("\nВероятности состояний буфера (экспоненциальное):")
	totalProb = 0.0
	for i := 0; i <= bufferSize; i++ {
		if prob, exists := expBufferProbs[i]; exists {
			fmt.Printf("  P(%d) = %.6f\n", i, prob)
			totalProb += prob
		}
	}
	fmt.Printf("Сумма вероятностей: %.6f\n", totalProb)

	// Дополнительный анализ: сравнение равномерного и экспоненциального распределений
	fmt.Println("\n=== СРАВНИТЕЛЬНЫЙ АНАЛИЗ ===")
	fmt.Println("Сравнение характеристик для разных законов распределения:")

	// Вычисляем среднюю длину очереди
	avgQueueUniform := 0.0
	for length, prob := range bufferProbs {
		avgQueueUniform += float64(length) * prob
	}

	avgQueueExp := 0.0
	for length, prob := range expBufferProbs {
		avgQueueExp += float64(length) * prob
	}

	// Вычисляем коэффициент использования сервера
	rhoUniform := mean(serviceTimes) / mean(interArrivalTimes)
	rhoExp := mean(expServiceTimes) / mean(expInterArrivalTimes)

	fmt.Printf("Равномерное распределение:\n")
	fmt.Printf("  Средняя длина очереди: %.4f\n", avgQueueUniform)
	fmt.Printf("  Коэффициент использования сервера (ρ): %.4f\n", rhoUniform)
	fmt.Printf("  Вероятность простоя сервера: %.4f\n", bufferProbs[0])

	fmt.Printf("\nЭкспоненциальное распределение:\n")
	fmt.Printf("  Средняя длина очереди: %.4f\n", avgQueueExp)
	fmt.Printf("  Коэффициент использования сервера (ρ): %.4f\n", rhoExp)
	fmt.Printf("  Вероятность простоя сервера: %.4f\n", expBufferProbs[0])

	// Теоретические значения для M/M/1 системы
	if rhoExp < 1 {
		theoryAvgQueue := rhoExp * rhoExp / (1 - rhoExp)
		fmt.Printf("\nТеоретические значения для M/M/1 (экспоненциальное):\n")
		fmt.Printf("  Теоретическая средняя длина очереди: %.4f\n", theoryAvgQueue)
		fmt.Printf("  Теоретическая вероятность простоя: %.4f\n", 1-rhoExp)
	}

	fmt.Println("\n=== РАБОТА ЗАВЕРШЕНА ===")
}

// Вспомогательная функция для вычисления среднего
func mean(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}
