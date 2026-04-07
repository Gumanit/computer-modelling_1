package main

import (
	"fmt"
	"math"
	"math/rand/v2"
)

var num int = 1000000
var x_min, x_max float64 = -1, 1

func main() {
	summa := 0.0
	sumSqLf := 0.0
	for i := 0; i < num; i++ {
		x := x_min + (x_max-x_min)*rand.Float64()
		fx := (1 + math.Pow(x, 2)) / (1 + math.Pow(x, 4))
		summa += fx
		sumSqLf += math.Pow((x_max-x_min)*fx, 2)
	}
	w := (x_max - x_min) * (summa / float64(num))
	disp := (sumSqLf - float64(num)*w*w) / (float64(num) - 1)
	err := math.Sqrt(disp / float64(num))
	halfWidth := 1.96 * err
	integValue := 2.2214415
	fmt.Printf("Значение оценки w = %.7f\n", w)
	fmt.Printf("Истинное значение интеграла = %v\n", integValue)
	fmt.Printf("Погрешность = %.7f\n", w-integValue)
	fmt.Printf("Доверительный интервал: [%.7f, %.7f]\n", w-halfWidth, w+halfWidth)
	fmt.Printf("Значение дисперсии S^2 = %.7f\n", disp)
}
