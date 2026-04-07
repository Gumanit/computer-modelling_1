package main

import (
	"fmt"
	"math"
	"math/rand/v2"
)

var num int = 10000000
var x_min, x_max float64 = -1, 1

func monteKarlo(num int) (float64, float64, []float64) {
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
	I := 2.2214415
	interval := []float64{w - halfWidth, w + halfWidth}
	return w, math.Abs(I - w), interval
}

func main() {
	w, pogr, interval := monteKarlo(num)
	fmt.Printf("Оценка w = %v\n", w)
	fmt.Printf("|I - w| = %.7f\n", pogr)
	fmt.Printf("Доверительный интервал = [%.8f, %.8f]\n", interval[0], interval[1])
}
