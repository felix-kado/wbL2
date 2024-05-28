package pattern

import "fmt"

/*
	Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern
*/

type Visitor interface {
	VisitBook(book *Book)
	VisitElectronics(electronics *Electronics)
}
type Element interface {
	Accept(visitor Visitor)
}
type Book struct {
	Title  string
	Author string
	Price  float64
}

func (b *Book) Accept(visitor Visitor) {
	visitor.VisitBook(b)
}

type Electronics struct {
	Name  string
	Brand string
	Price float64
}

func (e *Electronics) Accept(visitor Visitor) {
	visitor.VisitElectronics(e)
}

type DiscountVisitor struct {
	DiscountRate float64
}

func (d *DiscountVisitor) VisitBook(book *Book) {
	book.Price -= book.Price * d.DiscountRate
	fmt.Printf("Книга %s теперь стоит %.2f\n", book.Title, book.Price)
}

func (d *DiscountVisitor) VisitElectronics(electronics *Electronics) {
	electronics.Price -= electronics.Price * d.DiscountRate
	fmt.Printf("Электроника %s теперь стоит %.2f\n", electronics.Name, electronics.Price)
}

type PrintVisitor struct{}

func (p *PrintVisitor) VisitBook(book *Book) {
	fmt.Printf("Книга: %s от %s, Цена: %.2f\n", book.Title, book.Author, book.Price)
}

func (p *PrintVisitor) VisitElectronics(electronics *Electronics) {
	fmt.Printf("Электроника: %s от %s, Цена: %.2f\n", electronics.Name, electronics.Brand, electronics.Price)
}

func theoretical_main_3() {
	book := &Book{Title: "The Go Programming Language", Author: "Alan A. A. Donovan", Price: 50.00}
	electronics := &Electronics{Name: "Smartphone", Brand: "BrandX", Price: 300.00}

	discountVisitor := &DiscountVisitor{DiscountRate: 0.10}
	printVisitor := &PrintVisitor{}

	products := []Element{book, electronics}

	for _, product := range products {
		product.Accept(printVisitor)
		product.Accept(discountVisitor)
		product.Accept(printVisitor)
	}
}
