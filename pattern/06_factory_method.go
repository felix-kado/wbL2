package pattern

/*
	Реализовать паттерн «фабричный метод».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Factory_method_pattern
*/

/*
Применимость:
Когда создание объектов включает сложную логику, которую нужно изолировать.
Когда нужно избежать жесткой привязки к конкретным классам.

Плюсы
Уменьшает зависимость от конкретных классов.
Легко добавлять новые типы продуктов.

Минусы:
Усложнение кода -- Множество классов для каждого продукта.

*/

import "fmt"

// Интерфейс Product
type Product1 interface {
	GetName() string
	GetCategory() string
}

// Электронный продукт
type Electronic struct {
	name     string
	category string
}

func (e Electronic) GetName() string {
	return e.name
}

func (e Electronic) GetCategory() string {
	return e.category
}

// Продукт одежды
type Clothing struct {
	name     string
	category string
}

func (c Clothing) GetName() string {
	return c.name
}

func (c Clothing) GetCategory() string {
	return c.category
}

// Интерфейс Creator
type Creator interface {
	CreateProduct(name string) Product1
}

// ElectronicCreator создает электронные продукты
type ElectronicCreator struct{}

func (ec ElectronicCreator) CreateProduct(name string) Product1 {
	return &Electronic{
		name:     name,
		category: "Electronics",
	}
}

// ClothingCreator создает продукты одежды
type ClothingCreator struct{}

func (cc ClothingCreator) CreateProduct(name string) Product1 {
	return &Clothing{
		name:     name,
		category: "Clothing",
	}
}

func theoretical_main_6() {
	var creator Creator

	// Создание электронного продукта
	creator = ElectronicCreator{}
	electronic := creator.CreateProduct("Smartphone")
	fmt.Printf("Product: %s, Category: %s\n", electronic.GetName(), electronic.GetCategory())

	// Создание продукта одежды
	creator = ClothingCreator{}
	clothing := creator.CreateProduct("T-Shirt")
	fmt.Printf("Product: %s, Category: %s\n", clothing.GetName(), clothing.GetCategory())
}
