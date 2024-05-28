package main

/*
=== Базовая задача ===

Создать программу печатающую точное время с использованием NTP библиотеки.Инициализировать как go module.
Использовать библиотеку https://github.com/beevik/ntp.
Написать программу печатающую текущее время / точное время с использованием этой библиотеки.

Программа должна быть оформлена с использованием как go module.
Программа должна корректно обрабатывать ошибки библиотеки: распечатывать их в STDERR и возвращать ненулевой код выхода в OS.
Программа должна проходить проверки go vet и golint.
*/

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
	"github.com/spf13/pflag"
)

func main() {
	ntpServer := pflag.String("server", "pool.ntp.org", "Address of the NTP server to query for time")
	pflag.Parse()

	currentTime, err := ntp.Time(*ntpServer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching NTP time from server %s: %v\n", *ntpServer, err)
		os.Exit(1)
	}

	fmt.Printf("Current NTP time from server %s is: %s\n", *ntpServer, currentTime.Format(time.RFC1123))
}
