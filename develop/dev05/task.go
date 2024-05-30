package main

/*
Утилита grep

Реализовать утилиту фильтрации по аналогии с консольной утилитой (man grep — смотрим описание и основные параметры).

Реализовать поддержку утилитой следующих ключей:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", напечатать номер строки
*/

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type grepConfig struct {
	after      int
	before     int
	context    int
	count      int
	ignoreCase bool
	invert     bool
	fixed      bool
	lineNum    bool
}

func parseFlags() (grepConfig, string) {
	after := flag.Int("A", 0, "печатать +N строк после совпадения")
	before := flag.Int("B", 0, "печатать +N строк до совпадения")
	context := flag.Int("C", 0, "(A+B) печатать ±N строк вокруг совпадения")
	count := flag.Int("c", -1, "ограничение на количество совпадений")
	ignoreCase := flag.Bool("i", false, "игнорировать регистр")
	invert := flag.Bool("v", false, "вместо совпадения, исключать")
	fixed := flag.Bool("F", false, "точное совпадение со строкой, не паттерн")
	lineNum := flag.Bool("n", false, "печатать номер строки")

	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("Pattern is missing")
		os.Exit(1)
	}

	pattern := flag.Arg(0)

	if *context > 0 {
		*after = *context
		*before = *context
	}

	return grepConfig{
		after:      *after,
		before:     *before,
		context:    *context,
		count:      *count,
		ignoreCase: *ignoreCase,
		invert:     *invert,
		fixed:      *fixed,
		lineNum:    *lineNum,
	}, pattern
}

func main() {
	config, pattern := parseFlags()
	lines := readLines()

	// Копия строк в нижнем регистре для поиска
	searchLines := make([]string, len(lines))
	copy(searchLines, lines)

	if config.ignoreCase {
		pattern = strings.ToLower(pattern)
		for i := range searchLines {
			searchLines[i] = strings.ToLower(searchLines[i])
		}
	}

	matchInxs := []int{}

	if config.fixed {
		for i, line := range searchLines {
			if (!config.invert && line == pattern) || (config.invert && line != pattern) {
				matchInxs = append(matchInxs, i)
			}
		}
	} else {
		for i, line := range searchLines {
			matched, err := regexp.MatchString(pattern, line)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			if (config.invert && !matched) || (!config.invert && matched) {
				matchInxs = append(matchInxs, i)
			}
		}
	}

	limit := config.count
	for _, strInx := range matchInxs {
		if limit == 0 {
			break
		}
		start := max(0, strInx-config.before)
		end := min(len(lines), strInx+config.after+1)

		for i := start; i < end; i++ {
			if config.lineNum {
				fmt.Printf("%d: %s\n", i+1, lines[i])
			} else {
				fmt.Println(lines[i])
			}
		}
		if limit > 0 {
			limit--
		}
	}
}

func readLines() []string {
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	return lines
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
