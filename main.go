package main

import (
	"fmt"
	"math/rand"
)

//var distanceMap = [][]int{
//	{0, 120, 220, 150, 300}, //От города 0 (A)
//	{120, 0, 100, 180, 250}, //От города 1 (B)
//	{220, 100, 0, 80, 120},  //От города 2 (C)
//	{150, 180, 80, 0, 200},  //От города 3 (D)
//	{300, 250, 120, 200, 0}, //От города 4 (E)
//}

//
//func initCitiesMap() map[string]int {
//	cities := make(map[string]int)
//	for i := 0; i < len(InitialDecision); i++ {
//		cities[string(InitialDecision[i])] = i
//	}
//	return cities
//}
//
//func calculateObjFunc(cities map[string]int, path string) (objValue int) {
//	objValue = 0
//	fullPath := path + string(path[0])
//	for i := range fullPath {
//		if i+1 < len(fullPath) {
//			objValue += distanceMap[cities[string(fullPath[i])]][cities[string(fullPath[i+1])]]
//		}
//	}
//	return
//}

func initDec(n int) []int {
	var decision []int
	for i := 0; i < n; i++ {
		decision = append(decision, i)
	}
	return decision
}

// ====== CITY-SWAP АЛГОРИТМ ======

func citySwap(n int, neighborhoodSize int) [][]int {

	if n <= 2 || neighborhoodSize <= 0 {
		return [][]int{}
	}

	initialDecision := initDec(n)

	var neighbors [][]int

	maxPossible := (n - 1) * (n - 2) / 2

	if neighborhoodSize > maxPossible {
		neighborhoodSize = maxPossible
	}

	generatedPairs := make(map[[2]int]bool)

	for len(neighbors) < neighborhoodSize && len(generatedPairs) < maxPossible {
		i := rand.Intn(n-1) + 1
		j := rand.Intn(n-1) + 1

		for i == j {
			j = rand.Intn(n-1) + 1
		}

		small, large := i, j
		if small > large {
			small, large = large, small
		}
		pair := [2]int{small, large}

		if generatedPairs[pair] {
			continue
		}

		generatedPairs[pair] = true

		newPath := make([]int, n)
		copy(newPath, initialDecision)

		newPath[i], newPath[j] = newPath[j], newPath[i]

		neighbors = append(neighbors, newPath)
	}

	return neighbors
}

// ====== 2-OPT АЛГОРИТМ ======

func twoOpt(n int, neighborhoodSize int) [][]int {

	initialDecision := initDec(n)
	if n < 4 {
		return [][]int{}
	}

	var neighbors [][]int

	maxPossible := n * (n - 3) / 2

	if neighborhoodSize > maxPossible {
		neighborhoodSize = maxPossible
	}

	generatedPairs := make(map[[2]int]bool)

	for len(neighbors) < neighborhoodSize && len(generatedPairs) < maxPossible {
		i := rand.Intn(n - 2)

		j := i + 2 + rand.Intn(n-i-2)

		if j == n-1 && i == 0 {
			continue
		}

		pair := [2]int{i, j}
		if generatedPairs[pair] {
			continue
		}
		generatedPairs[pair] = true

		newPath := make([]int, n)

		copy(newPath[0:i+1], initialDecision[0:i+1])

		for k := 0; k <= j-i-1; k++ {
			newPath[i+1+k] = initialDecision[j-k]
		}

		if j+1 < n {
			copy(newPath[j+1:], initialDecision[j+1:])
		}

		neighbors = append(neighbors, newPath)
	}

	return neighbors
}

func main() {
	fmt.Println("City-swap:", citySwap(5, 30))
	fmt.Println("TwoOpt:", twoOpt(6, 30))
}
