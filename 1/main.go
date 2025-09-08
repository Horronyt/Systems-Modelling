package main

import (
	"fmt"
	"math"
	"math/rand"
)

// Задание 4: функция CALC_INTEGRAL

func CALC_INTEGRAL(a, b float64, f func(float64) float64, expNmb int) float64 {
	var IntegralValue float64

	var m int = 0
	var xMin, xMax, yMin, yMax float64 = a, b, 0, f(b)
	var p, x, y float64
	for i := 0; i < expNmb; i++ {
		p = rand.Float64()
		x = (xMax-xMin)*p + xMin
		p = rand.Float64()
		y = (yMax-yMin)*p + yMin
		if f(x) > y {
			m++
		}
	}

	IntegralValue = (float64(m) / float64(expNmb)) * (b - a) * f(b)

	return IntegralValue
}

// Задание 1: функция CALC_PI

func CALC_PI(x0, y0, r0 float64, expNmb int) float64 {
	var Pi float64

	var m int = 0
	var p, xp, yp float64
	for i := 0; i < expNmb; i++ {
		p = rand.Float64()
		xp = x0 - r0 + 2*r0*p
		p = rand.Float64()
		yp = y0 - r0 + 2*r0*p
		if (xp-x0)*(xp-x0)+(yp-y0)*(yp-y0) < r0*r0 {
			m++
		}
	}

	Pi = 4 * (float64(m) / float64(expNmb))

	return Pi
}

func main() {

	// Задание 2: расчет значения числа Пи для заданной окружности и различного
	// количества экспериментов

	var SERIA_1, SERIA_2, SERIA_3, SERIA_4, SERIA_5 [5]float64
	for i := 4; i <= 8; i++ {
		SERIA_1[i-4] = CALC_PI(1.0, 2.0, 5.0, int(math.Pow(10, float64(i))))
	}
	for i := 4; i <= 8; i++ {
		SERIA_2[i-4] = CALC_PI(1.0, 2.0, 5.0, int(math.Pow(10, float64(i))))
	}
	for i := 4; i <= 8; i++ {
		SERIA_3[i-4] = CALC_PI(1.0, 2.0, 5.0, int(math.Pow(10, float64(i))))
	}
	for i := 4; i <= 8; i++ {
		SERIA_4[i-4] = CALC_PI(1.0, 2.0, 5.0, int(math.Pow(10, float64(i))))
	}
	for i := 4; i <= 8; i++ {
		SERIA_5[i-4] = CALC_PI(1.0, 2.0, 5.0, int(math.Pow(10, float64(i))))
	}
	fmt.Println("Результаты вычисления значения числа Пи для первой серии экспериментов (SERIA_1):", SERIA_1)
	fmt.Println("Результаты вычисления значения числа Пи для второй серии экспериментов (SERIA_2):", SERIA_2)
	fmt.Println("Результаты вычисления значения числа Пи для третьей серии экспериментов (SERIA_3):", SERIA_3)
	fmt.Println("Результаты вычисления значения числа Пи для четвертой серии экспериментов (SERIA_4):", SERIA_4)
	fmt.Println("Результаты вычисления значения числа Пи для пятой серии экспериментов (SERIA_5):", SERIA_5, "\n")

	// Задание 3: расчет погрешности вычислений значений числа Пи

	var Eps1, Eps2, Eps3, Eps4, Eps5 [5]float64
	for i := 0; i < 5; i++ {
		Eps1[i] = math.Abs((SERIA_1[i] - math.Pi) / math.Pi)
		Eps2[i] = math.Abs((SERIA_2[i] - math.Pi) / math.Pi)
		Eps3[i] = math.Abs((SERIA_3[i] - math.Pi) / math.Pi)
		Eps4[i] = math.Abs((SERIA_4[i] - math.Pi) / math.Pi)
		Eps5[i] = math.Abs((SERIA_5[i] - math.Pi) / math.Pi)
	}
	fmt.Println("Погрешности вычислений значений числа Пи для первой серии экспериментов (SERIA_1):", Eps1)
	fmt.Println("Погрешности вычислений значений числа Пи для второй серии экспериментов (SERIA_2):", Eps2)
	fmt.Println("Погрешности вычислений значений числа Пи для третьей серии экспериментов (SERIA_3):", Eps3)
	fmt.Println("Погрешности вычислений значений числа Пи для четвертой серии экспериментов (SERIA_4):", Eps4)
	fmt.Println("Погрешности вычислений значений числа Пи для пятой серии экспериментов (SERIA_5):", Eps5, "\n")

	var S_e4, S_e5, S_e6, S_e7, S_e8 float64
	S_e4 = (SERIA_1[0] + SERIA_2[0] + SERIA_3[0] + SERIA_4[0] + SERIA_5[0]) / 5
	S_e5 = (SERIA_1[1] + SERIA_2[1] + SERIA_3[1] + SERIA_4[1] + SERIA_5[1]) / 5
	S_e6 = (SERIA_1[2] + SERIA_2[2] + SERIA_3[2] + SERIA_4[2] + SERIA_5[2]) / 5
	S_e7 = (SERIA_1[3] + SERIA_2[3] + SERIA_3[3] + SERIA_4[3] + SERIA_5[3]) / 5
	S_e8 = (SERIA_1[4] + SERIA_2[4] + SERIA_3[4] + SERIA_4[4] + SERIA_5[4]) / 5

	var Eps_S_e4, Eps_S_e5, Eps_S_e6, Eps_S_e7, Eps_S_e8 float64
	Eps_S_e4 = math.Abs((S_e4 - math.Pi) / math.Pi)
	Eps_S_e5 = math.Abs((S_e5 - math.Pi) / math.Pi)
	Eps_S_e6 = math.Abs((S_e6 - math.Pi) / math.Pi)
	Eps_S_e7 = math.Abs((S_e7 - math.Pi) / math.Pi)
	Eps_S_e8 = math.Abs((S_e8 - math.Pi) / math.Pi)
	fmt.Println("Погрешность вычислений для усредненного значения вычисленного числа Пи при ExpNmb=10^4:", Eps_S_e4)
	fmt.Println("Погрешность вычислений для усредненного значения вычисленного числа Пи при ExpNmb=10^5:", Eps_S_e5)
	fmt.Println("Погрешность вычислений для усредненного значения вычисленного числа Пи при ExpNmb=10^6:", Eps_S_e6)
	fmt.Println("Погрешность вычислений для усредненного значения вычисленного числа Пи при ExpNmb=10^7:", Eps_S_e7)
	fmt.Println("Погрешность вычислений для усредненного значения вычисленного числа Пи при ExpNmb=10^8:", Eps_S_e8, "\n")

	// Задание 4

	var INTEGRAL_SERIA_1, INTEGRAL_SERIA_2, INTEGRAL_SERIA_3 [4]float64
	for i := 4; i <= 7; i++ {
		INTEGRAL_SERIA_1[i-4] = CALC_INTEGRAL(0, 2.0,
			func(x float64) float64 {
				return x*x*x + 1
			},
			int(math.Pow(10, float64(i))))
	}
	for i := 4; i <= 7; i++ {
		INTEGRAL_SERIA_2[i-4] = CALC_INTEGRAL(0, 2.0,
			func(x float64) float64 {
				return x*x*x + 1
			},
			int(math.Pow(10, float64(i))))
	}
	for i := 4; i <= 7; i++ {
		INTEGRAL_SERIA_3[i-4] = CALC_INTEGRAL(0, 2.0,
			func(x float64) float64 {
				return x*x*x + 1
			},
			int(math.Pow(10, float64(i))))
	}
	fmt.Println("Результаты нахождения значения определенного интеграла функции y = x^3+1 для первой серии экспериментов:", INTEGRAL_SERIA_1)
	fmt.Println("Результаты нахождения значения определенного интеграла функции y = x^3+1 для второй серии экспериментов:", INTEGRAL_SERIA_2)
	fmt.Println("Результаты нахождения значения определенного интеграла функции y = x^3+1 для третьей серии экспериментов:", INTEGRAL_SERIA_3, "\n")

	// Значение CorrectIntergralValue равно значению интеграла функции y=x^3+1 на промежутке [0;2]
	// К этому значению можно прийти классическими методами расчета интеграла на бумаге
	const CorrectIntergralValue = 6

	var Integral_Eps1, Integral_Eps2, Integral_Eps3 [4]float64
	for i := 0; i < 4; i++ {
		Integral_Eps1[i] = math.Abs((INTEGRAL_SERIA_1[i] - CorrectIntergralValue) / CorrectIntergralValue)
		Integral_Eps2[i] = math.Abs((INTEGRAL_SERIA_2[i] - CorrectIntergralValue) / CorrectIntergralValue)
		Integral_Eps3[i] = math.Abs((INTEGRAL_SERIA_3[i] - CorrectIntergralValue) / CorrectIntergralValue)
	}
	fmt.Println("Погрешности нахождений значения определенного интеграла функции y = x^3+1 для первой серии экспериментов:", Integral_Eps1)
	fmt.Println("Погрешности нахождений значения определенного интеграла функции y = x^3+1 для второй серии экспериментов:", Integral_Eps2)
	fmt.Println("Погрешности нахождений значения определенного интеграла функции y = x^3+1 для третьей серии экспериментов:", Integral_Eps3, "\n")

	var Integral_S_e4, Integral_S_e5, Integral_S_e6, Integral_S_e7 float64
	Integral_S_e4 = (INTEGRAL_SERIA_1[0] + INTEGRAL_SERIA_2[0] + INTEGRAL_SERIA_3[0]) / 3
	Integral_S_e5 = (INTEGRAL_SERIA_1[1] + INTEGRAL_SERIA_2[1] + INTEGRAL_SERIA_3[1]) / 3
	Integral_S_e6 = (INTEGRAL_SERIA_1[2] + INTEGRAL_SERIA_2[2] + INTEGRAL_SERIA_3[2]) / 3
	Integral_S_e7 = (INTEGRAL_SERIA_1[3] + INTEGRAL_SERIA_2[3] + INTEGRAL_SERIA_3[3]) / 3

	var Integral_Eps_S_e4, Integral_Eps_S_e5, Integral_Eps_S_e6, Integral_Eps_S_e7 float64
	Integral_Eps_S_e4 = math.Abs((Integral_S_e4 - CorrectIntergralValue) / CorrectIntergralValue)
	Integral_Eps_S_e5 = math.Abs((Integral_S_e5 - CorrectIntergralValue) / CorrectIntergralValue)
	Integral_Eps_S_e6 = math.Abs((Integral_S_e6 - CorrectIntergralValue) / CorrectIntergralValue)
	Integral_Eps_S_e7 = math.Abs((Integral_S_e7 - CorrectIntergralValue) / CorrectIntergralValue)

	fmt.Println("Погрешность вычислений для усредненного значения найденного значения определенного интеграла функции y = x^3+1 при ExpNmb=10^4:", Integral_Eps_S_e4)
	fmt.Println("Погрешность вычислений для усредненного значения найденного значения определенного интеграла функции y = x^3+1 при ExpNmb=10^5:", Integral_Eps_S_e5)
	fmt.Println("Погрешность вычислений для усредненного значения найденного значения определенного интеграла функции y = x^3+1 при ExpNmb=10^6:", Integral_Eps_S_e6)
	fmt.Println("Погрешность вычислений для усредненного значения найденного значения определенного интеграла функции y = x^3+1 при ExpNmb=10^7:", Integral_Eps_S_e7)
}
