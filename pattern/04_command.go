package pattern

import "fmt"

/*
	Реализовать паттерн «комманда».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Command_pattern
*/

/*

Область:
- Когда необходимо ставить задачи в очередь и выполнять их в определенном порядке.
- Когда нужно записывать действия пользователя для анализа или отката.

Плюсы:
Разделение сущностей
Удобство логгирования операции
легко добавлять новые команды
удобно реализовывать отмену

минусы:
Код становиться более грамоздким
Требует больше памяти?

Ну в программах типо фотошопа для последовательностей действий наверняка можно было что-то такое устроить или в играх для последовательностей действий


*/

// Интерфейс дя команды
type Command interface {
	Execute()
}

// Структура которая умеет исполнять
type Receiver struct {
	name string
}

func (r *Receiver) PlaceOrder() {
	fmt.Println("Order placed for", r.name)
}

func (r *Receiver) CancelOrder() {
	fmt.Println("Order canceled for", r.name)
}

func (r *Receiver) UpdateInventory() {
	fmt.Println("Inventory updated for", r.name)
}

// Конкретная команда для размещенеия заказов
type PlaceOrderCommand struct {
	receiver *Receiver
}

func (c *PlaceOrderCommand) Execute() {
	c.receiver.PlaceOrder()
}

// Конкретная команда для отмены заказов
type CancelOrderCommand struct {
	receiver *Receiver
}

func (c *CancelOrderCommand) Execute() {
	c.receiver.CancelOrder()
}

// Конкретная команда для обновления инвенторя
type UpdateInventoryCommand struct {
	receiver *Receiver
}

func (c *UpdateInventoryCommand) Execute() {
	c.receiver.UpdateInventory()
}

// Струтура хранитель команд
type Invoker struct {
	commands []Command
}

func (i *Invoker) StoreCommand(command Command) {
	i.commands = append(i.commands, command)
}

func (i *Invoker) ExecuteCommands() {
	for _, command := range i.commands {
		command.Execute()
	}
	i.commands = nil
}

func theoretical_main_4() {
	receiver := &Receiver{name: "Product A"}

	placeOrderCommand := &PlaceOrderCommand{receiver: receiver}
	cancelOrderCommand := &CancelOrderCommand{receiver: receiver}
	updateInventoryCommand := &UpdateInventoryCommand{receiver: receiver}

	invoker := &Invoker{}
	invoker.StoreCommand(placeOrderCommand)
	invoker.StoreCommand(cancelOrderCommand)
	invoker.StoreCommand(updateInventoryCommand)

	invoker.ExecuteCommands()
}
