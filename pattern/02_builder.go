package pattern

/*
	Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern
*/

/*
Область применения;
1. Когда объект имеет множество параметров, особенно если некоторые из них являются необязательными
2. Если объект строится пошагово, и на каждом этапе требуется выполнение сложных операций.
3. Например когда объект не изменяется после сборки и есть какой-то момент в которой ему стоило бы быть полностью "собранным" за один раз


+ паттерна Строитель

Инкапсуляция логики: Сложная логика конструирования инкапсулируется в одном месте.
Читаемость код: отдельные части процесса конструирования становятся более читаемыми и понятным.
Функицональность: Разные представления объекта могут быть созданы с помощью одного и того же процесса конструирования.

*/

import "fmt"

// Product -- сложный объект, который мы хотии "построить"
type Product struct {
	Name       string
	Category   string
	Price      float64
	Discounted bool
}

// Вывод инфо
func (p *Product) Display() {
	fmt.Printf("Продукт: %s\nКатегория: %s\nЦена: $%.2f\nСкидка: %t\n", p.Name, p.Category, p.Price, p.Discounted)
}

// Builder — это абстрактный интерфейс для постройки продукта
type Builder interface {
	SetName(name string) Builder
	SetCategory(category string) Builder
	SetPrice(price float64) Builder
	ApplyDiscount(discounted bool) Builder
	Build() Product
}

// ProductBuilder — это конкретный строитель для продукта
type ProductBuilder struct {
	product Product
}

// Просто конструктор
func NewProductBuilder() *ProductBuilder {
	return &ProductBuilder{product: Product{}}
}

func (b *ProductBuilder) SetName(name string) Builder {
	b.product.Name = name
	return b
}

func (b *ProductBuilder) SetCategory(category string) Builder {
	b.product.Category = category
	return b
}

func (b *ProductBuilder) SetPrice(price float64) Builder {
	b.product.Price = price
	return b
}

func (b *ProductBuilder) ApplyDiscount(discounted bool) Builder {
	b.product.Discounted = discounted
	return b
}

func (b *ProductBuilder) Build() Product {
	return b.product
}

// директор отвечает за конструирование продукта, используя интерфейс строителя
type Director struct {
	builder Builder
}

// просто конструктор
func NewDirector(builder Builder) *Director {
	return &Director{builder: builder}
}

// конструирует электронный продукт
func (d *Director) ConstructElectronicProduct() Product {
	return d.builder.SetName("Смартфон").
		SetCategory("Электроника").
		SetPrice(699.99).
		ApplyDiscount(true).
		Build()
}

// конструирует книжный продукт
func (d *Director) ConstructBookProduct() Product {
	return d.builder.SetName("Язык Программирования Go").
		SetCategory("Книги").
		SetPrice(39.99).
		ApplyDiscount(false).
		Build()
}

func theoretical_main_1() {
	builder := NewProductBuilder()
	director := NewDirector(builder)

	electronicProduct := director.ConstructElectronicProduct()
	electronicProduct.Display()

	bookProduct := director.ConstructBookProduct()
	bookProduct.Display()
}
