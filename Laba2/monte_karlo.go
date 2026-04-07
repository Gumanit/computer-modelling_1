package main

import (
	"fmt"
	"math"
	"math/rand/v2"
)

var num int = 1000000
var x_min, x_max float64 = -1, 1
var y_min, y_max = 0.0, 1.3

func main() {

	var up, down float64

	for i := 0; i < num; i++ {
		x := x_min + rand.Float64()*(x_max-x_min)
		y := y_min + rand.Float64()*(y_max-y_min)

		funcVal := (1 + math.Pow(x, 2)) / (1 + math.Pow(x, 4))
		if y > funcVal {
			up++
		} else {
			down++
		}
	}

	rectArea := (x_max - x_min) * (y_max - y_min)
	w := down / (up + down)
	I := rectArea * w
	sqrtVal := math.Sqrt((w * (1.0 - w)) / float64(num))
	delta_I := rectArea * 1.96 * sqrtVal
	fmt.Printf("Общее количество точек = %v\n", num)
	fmt.Printf("Доля точек ниже графика функции: %v\n", w)
	fmt.Printf("I = %v\n", I)
	fmt.Printf("Доверительный интервал: [%v, %v], delta I = %v\n", I-delta_I, I+delta_I, delta_I)
	fmt.Printf("I точное = %v\n", 2.2214415)
}
