package pattern

import "fmt"

/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern
*/

/*
Пример использования
Обработка HTTP-запросов -- веб-серверы используют этот паттерн для обработки запросов через цепочку middleware, которая может выполнять аутентификацию, логирование, валидацию данных и другие задачи.

Плюсы:
Все обработчики отделены в понятную структуру, куда можно легко добавить или убрать обработчик


*/
// Интерфейс обработчика
type Handler interface {
	SetNext(handler Handler)
	Handle(request *SaleRequest)
}

type SaleRequest struct {
	ProductID     int
	Quantity      int
	Price         float64
	Discount      float64
	Tax           float64
	IsDiscounted  bool
	IsTaxed       bool
	IsInInventory bool
}

type BaseHandler struct {
	next Handler
}

func (h *BaseHandler) SetNext(handler Handler) {
	h.next = handler
}

func (h *BaseHandler) Handle(request *SaleRequest) {
	if h.next != nil {
		h.next.Handle(request)
	}
}

type DiscountHandler struct {
	BaseHandler
}

func (h *DiscountHandler) Handle(request *SaleRequest) {
	if request.Quantity > 10 {
		request.Discount = request.Price * 0.1
		request.IsDiscounted = true
		fmt.Println("Скидка применена")
	}
	h.BaseHandler.Handle(request)
}

type TaxHandler struct {
	BaseHandler
}

func (h *TaxHandler) Handle(request *SaleRequest) {
	if request.IsDiscounted {
		request.Tax = (request.Price - request.Discount) * 0.2
	} else {
		request.Tax = request.Price * 0.2
	}
	request.IsTaxed = true
	fmt.Println("Налог применен")
	h.BaseHandler.Handle(request)
}

type InventoryHandler struct {
	BaseHandler
	inventory map[int]int
}

func NewInventoryHandler(inventory map[int]int) *InventoryHandler {
	return &InventoryHandler{inventory: inventory}
}

func (h *InventoryHandler) Handle(request *SaleRequest) {
	if h.inventory[request.ProductID] >= request.Quantity {
		request.IsInInventory = true
		fmt.Println("Проверка наличия на складе пройдена")
	} else {
		fmt.Println("Проверка наличия на складе не пройдена")
		return
	}
	h.BaseHandler.Handle(request)
}

func theoretical_main_5() {
	// Настройка склада
	inventory := map[int]int{
		1: 20,
		2: 5,
	}

	// Обработчики
	discountHandler := &DiscountHandler{}
	taxHandler := &TaxHandler{}
	inventoryHandler := NewInventoryHandler(inventory)

	// Настройка цепочки
	discountHandler.SetNext(taxHandler)
	taxHandler.SetNext(inventoryHandler)

	// Создание запроса на продажу
	saleRequest := &SaleRequest{
		ProductID: 1,
		Quantity:  15,
		Price:     100.0,
	}

	// Обработка запроса на продажу
	discountHandler.Handle(saleRequest)

	// Вывод результатов
	fmt.Printf("Итоговый запрос на продажу: %+v\n", saleRequest)
}
