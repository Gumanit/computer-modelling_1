package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// InventoryParams содержит все входные параметры модели.
type InventoryParams struct {
	Lambda     float64 // интенсивность потока заказов (заказов в день)
	S          float64 // максимальный уровень запаса (целевой)
	s          float64 // критический уровень (точка заказа)
	R          float64 // цена продажи единицы продукции
	H          float64 // стоимость хранения единицы в день
	K          float64 // фиксированная стоимость размещения заказа поставщику
	L          float64 // время доставки от поставщика (дни)
	T          float64 // горизонт моделирования (дни)
	MeanDemand float64 // средний размер заказа (параметр Пуассона)
	Penalty    float64 // штраф за единицу отложенного заказа в день (0 = не учитывать)
}

// InventorySystem хранит состояние системы управления запасами.
type InventorySystem struct {
	params       InventoryParams
	t            float64    // текущее время
	stock        float64    // свободный запас на складе
	onOrder      float64    // количество товара, заказанного у поставщика (в пути)
	backorders   float64    // суммарное количество отложенных заказов
	totalHold    float64    // общие затраты на хранение
	totalOrder   float64    // общие затраты на заказы поставщику
	totalRevenue float64    // общий доход (признаётся при отгрузке)
	totalPenalty float64    // общий штраф за дефицит (если penalty > 0)
	nextArrival  float64    // время поступления следующего заказа
	nextDelivery float64    // время прибытия следующей поставки от поставщика
	rng          *rand.Rand // локальный генератор случайных чисел
}

// NewInventorySystem создаёт новый экземпляр системы с заданными параметрами и seed.
func NewInventorySystem(params InventoryParams, seed int64) *InventorySystem {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	rng := rand.New(rand.NewSource(seed))

	sys := &InventorySystem{
		params:       params,
		t:            0.0,
		stock:        params.S,
		onOrder:      0.0,
		backorders:   0.0,
		totalHold:    0.0,
		totalOrder:   0.0,
		totalRevenue: 0.0,
		totalPenalty: 0.0,
		nextDelivery: math.MaxFloat64, // нет активного заказа
		rng:          rng,
	}
	// Генерируем время первого заказа
	sys.nextArrival = sys.generateInterarrival()
	return sys
}

// generateInterarrival возвращает промежуток между возникнованием спроса (экспоненциальное).
func (sys *InventorySystem) generateInterarrival() float64 {
	u := sys.rng.Float64()
	// fmt.Printf("Промежутки времени: %v\n", -math.Log(u) / sys.params.Lambda)
	return -math.Log(u) / sys.params.Lambda
}

// generateDemand возвращает размер заказа (распределение Пуассона).
func (sys *InventorySystem) generateDemand() float64 {
	L := math.Exp(-sys.params.MeanDemand)
	k := 0
	p := 1.0
	for p > L {
		k++
		p *= sys.rng.Float64()
	}
	return float64(k - 1)
}

// Run выполняет имитацию до времени T.
func (sys *InventorySystem) Run() {
	for sys.t < sys.params.T {
		// Ближайшее событие: поступление заказа или доставка
		nextEvent := math.Min(sys.nextArrival, sys.nextDelivery)
		if nextEvent > sys.params.T {
			break
		}

		// Начисляем затраты на хранение и штраф за дефицит за прошедший интервал
		elapsed := nextEvent - sys.t
		sys.totalHold += sys.stock * sys.params.H * elapsed
		if sys.params.Penalty > 0 {
			sys.totalPenalty += sys.backorders * sys.params.Penalty * elapsed
		}
		sys.t = nextEvent

		// Обработка события
		if sys.nextArrival <= sys.nextDelivery {
			sys.handleCustomerOrder()
		} else {
			sys.handleDelivery()
		}
	}

	// Учёт затрат за оставшееся время до T
	elapsed := sys.params.T - sys.t
	sys.totalHold += sys.stock * sys.params.H * elapsed
	if sys.params.Penalty > 0 {
		sys.totalPenalty += sys.backorders * sys.params.Penalty * elapsed
	}
}

// handleCustomerOrder обрабатывает поступление заказа
func (sys *InventorySystem) handleCustomerOrder() {
	demand := sys.generateDemand()

	// добавить тут для массива

	// Сколько можем отгрузить немедленно из свободного запаса
	shipped := math.Min(demand, sys.stock)
	sys.stock -= shipped
	sys.totalRevenue += shipped * sys.params.R

	// Остаток добавляем в отложенные заказы
	shortage := demand - shipped
	sys.backorders += shortage

	// Проверка необходимости размещения заказа поставщику
	// Используем "позицию запаса" = stock - backorders + onOrder
	inventoryPosition := sys.stock - sys.backorders + sys.onOrder
	if inventoryPosition < sys.params.s && sys.onOrder == 0 {
		orderQty := sys.params.S - inventoryPosition
		if orderQty > 0 {
			sys.onOrder = orderQty
			sys.totalOrder += sys.params.K
			sys.nextDelivery = sys.t + sys.params.L
		}
	}

	// Генерируем время следующего заказа
	sys.nextArrival = sys.t + sys.generateInterarrival()
}

// handleDelivery обрабатывает прибытие поставки от поставщика.
func (sys *InventorySystem) handleDelivery() {
	delivery := sys.onOrder
	sys.onOrder = 0

	// Сначала удовлетворяем отложенные заказы
	if sys.backorders > 0 {
		fulfilled := math.Min(delivery, sys.backorders)
		sys.backorders -= fulfilled
		delivery -= fulfilled
		// Доход от отложенных заказов признаётся в момент отгрузки
		sys.totalRevenue += fulfilled * sys.params.R
	}

	// Остаток добавляем к свободному запасу
	sys.stock += delivery

	// Сбрасываем время доставки
	sys.nextDelivery = math.MaxFloat64
}

// Statistics возвращает итоговые показатели прогона.
func (sys *InventorySystem) Statistics() (profit, avgDailyProfit float64) {
	totalCost := sys.totalOrder + sys.totalHold + sys.totalPenalty
	profit = sys.totalRevenue - totalCost
	avgDailyProfit = profit / sys.params.T
	return
}

// GetDetailed возвращает все компоненты для усреднения.
func (sys *InventorySystem) GetDetailed() (revenue, orderCost, holdCost, penalty float64) {
	return sys.totalRevenue, sys.totalOrder, sys.totalHold, sys.totalPenalty
}

// SimulationResult содержит усреднённые результаты множественных прогонов.
type SimulationResult struct {
	AvgProfit      float64
	AvgRevenue     float64
	AvgOrderCost   float64
	AvgHoldCost    float64
	AvgPenalty     float64
	AvgDailyProfit float64
}

// RunMultipleSimulations выполняет numRuns независимых прогонов и возвращает средние значения.
func RunMultipleSimulations(params InventoryParams, numRuns int) SimulationResult {
	var sumProfit, sumRevenue, sumOrder, sumHold, sumPenalty float64

	for i := 0; i < numRuns; i++ {
		// Для каждого прогона новый seed (0 = авто)
		sys := NewInventorySystem(params, 0)
		sys.Run()
		p, _ := sys.Statistics()
		r, o, h, penalty := sys.GetDetailed()
		sumProfit += p
		sumRevenue += r
		sumOrder += o
		sumHold += h
		sumPenalty += penalty
	}

	n := float64(numRuns)
	return SimulationResult{
		AvgProfit:      sumProfit / n,
		AvgRevenue:     sumRevenue / n,
		AvgOrderCost:   sumOrder / n,
		AvgHoldCost:    sumHold / n,
		AvgPenalty:     sumPenalty / n,
		AvgDailyProfit: (sumProfit / n) / params.T,
	}
}

func main() {
	// Базовые параметры модели
	baseParams := InventoryParams{
		Lambda:     1.0,   // в среднем 1 заказ в день
		S:          100.0, // целевой запас
		s:          50.0,  // точка заказа
		R:          10.0,  // цена за единицу
		H:          0.1,   // стоимость хранения единицы в день
		K:          50.0,  // фикс. стоимость заказа поставщику
		L:          2.0,   // доставка 2 дня
		T:          365.0, // моделируем год
		MeanDemand: 5.0,   // средний размер заказа
		Penalty:    0.5,   // штраф за дефицит у.е. за единицу в день
	}

	fmt.Println("=== Имитационная модель управления запасами (s, S) с отложенными заказами ===")
	fmt.Printf("Параметры: S=%.0f, s=%.0f, L=%.1f, T=%.0f, λ=%.1f, ср.спрос=%.1f\n",
		baseParams.S, baseParams.s, baseParams.L, baseParams.T, baseParams.Lambda, baseParams.MeanDemand)
	fmt.Println()

	// Одиночный прогон
	fmt.Println("--- Одиночный прогон ---")
	single := NewInventorySystem(baseParams, 42)
	single.Run()
	_, o, h, p := single.GetDetailed()

	fmt.Printf("Общие затраты: %.2f\n", o+h+p)
	fmt.Printf("Затраты на заказы: %.2f\n", o)
	fmt.Printf("Затраты на хран.:  %.2f\n", h)
	fmt.Printf("Штраф за дефицит:  %.2f\n", p)
	fmt.Println()

	// Множественные прогоны для получения устойчивых оценок
	fmt.Println("--- Множественные прогоны (100) ---")
	result := RunMultipleSimulations(baseParams, 100)
	fmt.Printf("Средние общие затраты: %.2f\n", result.AvgOrderCost+result.AvgHoldCost+result.AvgPenalty)
	fmt.Printf("Средние затраты на заказы: %.2f\n", result.AvgOrderCost)
	fmt.Printf("Средние затраты на хран.:  %.2f\n", result.AvgHoldCost)
	fmt.Printf("Средний штраф за дефицит:  %.2f\n", result.AvgPenalty)
	fmt.Println()

	// Сравнение различных стратегий (s, S)
	fmt.Println("--- Сравнение стратегий (s, S) ---")
	strategies := []struct{ s, S float64 }{
		{50, 100},
		{50, 150},
		{75, 150},
		{100, 200},
		{150, 300},
		{200, 250},
		{250, 400},
		{10, 500},
		{100, 400},
		{5, 15},
		{10, 20},
		{15, 30},
	}

	var bestS, best_s float64
	var bestCost float64 = math.MaxFloat64

	for _, st := range strategies {
		params := baseParams
		params.s = st.s
		params.S = st.S
		res := RunMultipleSimulations(params, 1)
		fmt.Printf("s=%.0f S=%.0f -> Общие затраты: %.2f, Затраты на заказы: %.2f, Затраты на хран.:  %.2f, Штраф за дефицит:  %.2f\n",
			st.s, st.S, res.AvgOrderCost + res.AvgHoldCost + res.AvgPenalty, res.AvgOrderCost, res.AvgHoldCost, res.AvgPenalty)

		if res.AvgOrderCost + res.AvgHoldCost + res.AvgPenalty < bestCost {
			bestCost = res.AvgOrderCost + res.AvgHoldCost + res.AvgPenalty
			bestS = st.S
			best_s = st.s
		}
	}

	fmt.Printf("\n=== Лучшая стратегия: s=%.0f S=%.0f с затратами %.2f ===\n", best_s, bestS, bestCost)
}
