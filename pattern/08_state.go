package pattern

/*
	Реализовать паттерн «состояние».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/State_pattern
*/

/*
Паттерн Состояние полезен, когда объект должен менять свое поведение в зависимости от своего состояния,
 а количество возможных состояний известно и ограничено. Такой, конечный автомат по сути

Недавно исходники sunc mutex читал, вполне себе автомат конечных состояний.

Плюсы:
Упрощает управление состояниями --  легко добавлять новые состояния и изменять поведение без изменения кода контекста.
Разделение логики по состояниям делает код более поддерживаемым и модульным.

Минусы:
Иногда избыточен, для простых объектов с небольшим количеством состояний этот паттерн может быть избыточным и усложнять код
*/

import "fmt"

// Интерфейс состояния
type State interface {
	ExecuteCommand(*Spacecraft)
}

// Космический корабль представляет контекст
type Spacecraft struct {
	currentState State
}

func (sc *Spacecraft) SetState(state State) {
	sc.currentState = state
}

func (sc *Spacecraft) ExecuteCommand() {
	sc.currentState.ExecuteCommand(sc)
}

// Состояние на земле
type GroundedState struct{}

func (s *GroundedState) ExecuteCommand(sc *Spacecraft) {
	fmt.Println("Космический корабль на земле. Подготовка к запуску...")
	sc.SetState(&LaunchingState{})
}

// Состояние запуска
type LaunchingState struct{}

func (s *LaunchingState) ExecuteCommand(sc *Spacecraft) {
	fmt.Println("Космический корабль запускается. Идет обратный отсчет...")
	sc.SetState(&OrbitingState{})
}

// Состояние на орбите
type OrbitingState struct{}

func (s *OrbitingState) ExecuteCommand(sc *Spacecraft) {
	fmt.Println("Космический корабль на орбите. Выполнение орбитальных маневров...")
	sc.SetState(&LandingState{})
}

// Состояние посадки
type LandingState struct{}

func (s *LandingState) ExecuteCommand(sc *Spacecraft) {
	fmt.Println("Космический корабль совершает посадку. Подготовка к приземлению...")
	sc.SetState(&GroundedState{})
}

func theoretical_main_8() {
	spacecraft := &Spacecraft{currentState: &GroundedState{}}

	spacecraft.ExecuteCommand() // Космический корабль на земле. Подготовка к запуску...
	spacecraft.ExecuteCommand() // Космический корабль запускается. Идет обратный отсчет...
	spacecraft.ExecuteCommand() // Космический корабль на орбите. Выполнение орбитальных маневров...
	spacecraft.ExecuteCommand() // Космический корабль совершает посадку. Подготовка к приземлению...
	spacecraft.ExecuteCommand() // Космический корабль на земле. Подготовка к запуску...
}
