package main

/*
=== Взаимодействие с ОС ===

Необходимо реализовать собственный шелл

встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах

Реализовать утилиту netcat (nc) клиент
принимать данные из stdin и отправлять в соединение (tcp/udp)
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("shell> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		input = strings.TrimSpace(input)
		args := strings.Split(input, " ")
		if len(args) == 0 {
			continue
		}
		switch args[0] {
		case "cd":
			if len(args) < 2 {
				fmt.Println("cd: missing argument")
			} else {
				err := os.Chdir(args[1])
				if err != nil {
					fmt.Println("cd:", err)
				}
			}
		case "pwd":
			dir, err := os.Getwd()
			if err != nil {
				fmt.Println("pwd:", err)
			} else {
				fmt.Println(dir)
			}
		case "echo":
			fmt.Println(strings.Join(args[1:], " "))
		case "kill":
			if len(args) < 2 {
				fmt.Println("kill: missing argument")
			} else {
				pid := args[1]
				killProcess(pid)
			}
		case "ps":
			psCommand()
		case "exit":
			os.Exit(0)
		default:
			runExternalCommand(args)
		}
	}
}

func killProcess(pid string) {
	cmd := exec.Command("kill", pid)
	err := cmd.Run()
	if err != nil {
		fmt.Println("kill:", err)
	}
}

func psCommand() {
	cmd := exec.Command("ps")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("ps:", err)
	}
}

func runExternalCommand(args []string) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				fmt.Printf("Command exited with status %d\n", status.ExitStatus())
			}
		} else {
			fmt.Println("Error running command:", err)
		}
	}
}
