package main

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type Config struct {
	Timeout time.Duration
	Host    string
	Port    string
}

func parseConfig() (*Config, error) {
	timeoutFlag := flag.Duration("timeout", 10*time.Second, "Connection timeout")
	flag.Parse()
	if len(flag.Args()) < 2 {
		return nil, fmt.Errorf("usage: go-telnet [--timeout=10s] host port")
	}
	host := flag.Arg(0)
	port := flag.Arg(1)

	return &Config{
		Timeout: *timeoutFlag,
		Host:    host,
		Port:    port,
	}, nil
}

func main() {
	config, err := parseConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	address := net.JoinHostPort(config.Host, config.Port)

	// Подключение к серверу
	conn, err := net.DialTimeout("tcp", address, config.Timeout)
	if err != nil {
		fmt.Printf("Failed to connect: %v\n", err)
		return
	}
	defer conn.Close()
	fmt.Printf("Connected to %s\n", address)

	// Каналы для завершения работы
	done := make(chan struct{})

	// Чтение из сокета и запись в STDOUT
	go func() {
		if _, err := io.Copy(os.Stdout, conn); err != nil {
			fmt.Printf("Error reading from connection: %v\n", err)
		}
		fmt.Println("Connection closed by server")
		done <- struct{}{}
	}()

	// Чтение из STDIN и запись в сокет
	go func() {
		if _, err := io.Copy(conn, os.Stdin); err != nil {
			fmt.Printf("Error writing to connection: %v\n", err)
		}
		done <- struct{}{}
	}()

	// Ожидание завершения одной из горутин
	<-done
	fmt.Println("Exiting...")
}
