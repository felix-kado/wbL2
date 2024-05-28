package main

import (
	"testing"
	"time"
)

// helper функция для создания канала и закрытия его сразу
func closedChan() <-chan interface{} {
	c := make(chan interface{})
	close(c)
	return c
}

func TestOrEmpty(t *testing.T) {
	result := or()
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}
}

func TestOrSingle(t *testing.T) {
	c := closedChan()
	result := or(c)

	time.Sleep(time.Second * 1)
	select {
	case <-result:
		// Ожидаем, что результат закроется сразу
	default:
		t.Errorf("Expected closed channel, but it was not closed")
	}
}

func TestOrMultiple(t *testing.T) {
	c1 := closedChan()
	c2 := make(chan interface{})
	result := or(c1, c2)

	time.Sleep(time.Second * 1)

	select {
	case <-result:
		// Ожидаем, что результат закроется сразу, так как c1 закрыт
	default:
		t.Errorf("Expected closed channel, but it was not closed")
	}
}

func TestOrAllClosed(t *testing.T) {
	c1 := closedChan()
	c2 := closedChan()
	c3 := closedChan()
	result := or(c1, c2, c3)

	time.Sleep(time.Second * 1)

	select {
	case <-result:
		// Ожидаем, что результат закроется сразу, так как все каналы закрыты
	default:
		t.Errorf("Expected closed channel, but it was not closed")
	}
}

func TestOrMixed(t *testing.T) {
	c1 := make(chan interface{})
	c2 := closedChan()
	c3 := make(chan interface{})
	result := or(c1, c2, c3)

	time.Sleep(time.Second * 1)

	select {
	case <-result:
		// Ожидаем, что результат закроется сразу, так как c2 закрыт
	default:
		t.Errorf("Expected closed channel, but it was not closed")
	}
}
