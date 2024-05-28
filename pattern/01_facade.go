package pattern

/*
	Реализовать паттерн «фасад».
Объяснить применимость паттерна, его плюсы и минусы,а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Facade_pattern
*/
import "fmt"

// Подсистема инвентаря
type Inventory interface {
	CheckStock(productID string) bool
}

type inventoryImpl struct{}

func (i *inventoryImpl) CheckStock(productID string) bool {
	fmt.Println("Проверка наличия продукта:", productID)
	// do something

	return true
}

// Подсистема биллинга
type Billing interface {
	CreateInvoice(orderID string) bool
}

type billingImpl struct{}

func (b *billingImpl) CreateInvoice(orderID string) bool {
	fmt.Println("Создание счета для заказа:", orderID)
	// do something
	return true
}

// Подсистема доставки
type Shipping interface {
	ArrangeShipping(orderID string) bool
}

type shippingImpl struct{}

func (s *shippingImpl) ArrangeShipping(orderID string) bool {
	fmt.Println("Организация доставки для заказа:", orderID)
	// do something
	return true
}

type SalesFacade struct {
	inventory Inventory
	billing   Billing
	shipping  Shipping
}

func NewSalesFacade() *SalesFacade {
	return &SalesFacade{
		inventory: &inventoryImpl{},
		billing:   &billingImpl{},
		shipping:  &shippingImpl{},
	}
}

func (sf *SalesFacade) PlaceOrder(productID string, orderID string) bool {
	if !sf.inventory.CheckStock(productID) {
		fmt.Println("Продукт отсутствует на складе.")
		return false
	}
	if !sf.billing.CreateInvoice(orderID) {
		fmt.Println("Не удалось создать счет.")
		return false
	}
	if !sf.shipping.ArrangeShipping(orderID) {
		fmt.Println("Не удалось организовать доставку.")
		return false
	}
	fmt.Println("Заказ успешно размещен!")
	return true
}

func theoretical_main() {
	facade := NewSalesFacade()
	productID := "P12345"
	orderID := "O67890"

	success := facade.PlaceOrder(productID, orderID)
	if success {
		fmt.Println("Обработка заказа завершена.")
	} else {
		fmt.Println("Обработка заказа не удалась.")
	}
}
