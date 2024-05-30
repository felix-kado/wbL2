package main

/*
Необходимо реализовать свой собственный UNIX-шелл-утилиту с поддержкой ряда простейших команд:

- cd <args> - смена директории (в качестве аргумента могут быть то-то и то)
- pwd - показать путь до текущего каталога
- echo <args> - вывод аргумента в STDOUT
- kill <args> - "убить" процесс, переданный в качесте аргумента (пример: такой-то пример)
- ps - выводит общую информацию по запущенным процессам в формате *такой-то формат*

Так же требуется поддерживать функционал fork/exec-команд

Дополнительно необходимо поддерживать конвейер на пайпах (linux pipes, пример cmd1 | cmd2 | .... | cmdN).

*Шелл — это обычная консольная программа, которая будучи запущенной, в интерактивном сеансе выводит некое приглашение
в STDOUT и ожидает ввода пользователя через STDIN. Дождавшись ввода, обрабатывает команду согласно своей логике
и при необходимости выводит результат на экран. Интерактивный сеанс поддерживается до тех пор, пока не будет введена команда выхода (например \quit).
*/

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

// Интерфейс Command для всех команд шелла для реализации паттерна
type Command interface {
	Execute(args []string) error
}

type CDCommand struct{}

func (c *CDCommand) Execute(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("cd: отсутствует аргумент")
	}
	return os.Chdir(args[0])
}

type PWDCommand struct{}

func (p *PWDCommand) Execute(args []string) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Println(dir)
	return nil
}

type EchoCommand struct{}

func (e *EchoCommand) Execute(args []string) error {
	fmt.Println(strings.Join(args, " "))
	return nil
}

type KillCommand struct{}

func (k *KillCommand) Execute(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("kill: отсутствует аргумент")
	}
	pid, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("kill: неверный PID")
	}
	return syscall.Kill(pid, syscall.SIGKILL)
}

type PSCommand struct{}

func (p *PSCommand) Execute(args []string) error {
	cmd := exec.Command("ps")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type Shell struct {
	commands map[string]Command
}

func NewShell() *Shell {
	return &Shell{
		commands: map[string]Command{
			"cd":   &CDCommand{},
			"pwd":  &PWDCommand{},
			"echo": &EchoCommand{},
			"kill": &KillCommand{},
			"ps":   &PSCommand{},
		},
	}
}

func (s *Shell) ExecuteCommand(input string) error {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return nil
	}
	cmdName := parts[0]
	args := parts[1:]

	cmd, found := s.commands[cmdName]
	if found {
		return cmd.Execute(args)
	}

	// Попытка выполнить как системную команду
	return s.forkExecCommand(cmdName, args)
}

func (s *Shell) forkExecCommand(name string, args []string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	shell := NewShell()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "exit" {
			break
		}
		err := shell.ExecuteCommand(input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Ошибка:", err)
		}
	}
}
