package pattern

/*
	Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Strategy_pattern
*/

/*
Реальный пример -- мой L0 я там БД реализовывал черещ стратегию.
 С кешем она или нет не важно в рамках моей программы, какую нужно -- такую и подставляем
 Они все реализуют интерфейс DB

Когда применять
1. Когда у вас есть несколько алгоритмов для выполнения определенной задачи, и вы хотите выбрать один из них в зависимости от условий.
2. Когда вы хотите изолировать алгоритмы от контекста, в котором они применяются, чтобы облегчить тестирование и поддержку.

Плюсы
Легко изменять алгоритмы, не изменяя классы, которые их используют.
Новые алгоритмы можно добавлять без изменений в существующем коде.
Упрощается тестирование, так как они изолированы в своих собственных классах.

Минусы
Ну код более громоздкий, если реализация одна и точно не будет меняться, то это излишне
*/

import "fmt"

// интерфейс для всех алгоритмов скидок
type DiscountStrategy interface {
	ApplyDiscount(price float64) float64
}

// конкретная стратегия, которая не применяет скидку
type NoDiscount struct{}

func (d NoDiscount) ApplyDiscount(price float64) float64 {
	return price
}

// конкретная стратегия, которая применяет процентную скидку
type PercentageDiscount struct {
	percentage float64
}

func (d PercentageDiscount) ApplyDiscount(price float64) float64 {
	return price * (1 - d.percentage/100)
}

// конкретная стратегия, которая применяет фиксированную скидку
type FixedAmountDiscount struct {
	amount float64
}

func (d FixedAmountDiscount) ApplyDiscount(price float64) float64 {
	return price - d.amount
}

// продукт с ценой
type Product2 struct {
	name     string
	price    float64
	strategy DiscountStrategy
}

// устанавливает стратегию скидок для продукта
func (p *Product2) SetDiscountStrategy(strategy DiscountStrategy) {
	p.strategy = strategy
}

// возвращает цену после применения стратегии скидок
func (p *Product2) GetPrice() float64 {
	return p.strategy.ApplyDiscount(p.price)
}

func theoretical_main_7() {
	// Создание продукта
	product := Product2{name: "iphone 2024 pro mini", price: 100.0}

	// Применение стратегии без скидки
	product.SetDiscountStrategy(NoDiscount{})
	fmt.Printf("Цена без скидки: $%.2f\n", product.GetPrice())

	// Применение процентной скидки
	product.SetDiscountStrategy(PercentageDiscount{percentage: 10})
	fmt.Printf("Цена с 10%% скидкой: $%.2f\n", product.GetPrice())

	// Применение фиксированной скидки
	product.SetDiscountStrategy(FixedAmountDiscount{amount: 15})
	fmt.Printf("Цена с $15 скидкой: $%.2f\n", product.GetPrice())
}
