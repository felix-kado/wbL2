package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
=== Утилита sort ===

Отсортировать строки
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type sortConfig struct {
	column    int
	numeric   bool
	reverse   bool
	unique    bool
	delimiter string
}

func parseFlags() sortConfig {
	column := flag.Int("k", 1, "column to sort by (1-based index)")
	numeric := flag.Bool("n", false, "sort by numeric value")
	reverse := flag.Bool("r", false, "sort in reverse order")
	unique := flag.Bool("u", false, "output unique lines only")
	delimiter := flag.String("t", " ", "column delimiter")
	flag.Parse()

	return sortConfig{
		column:    *column,
		numeric:   *numeric,
		reverse:   *reverse,
		unique:    *unique,
		delimiter: *delimiter,
	}
}

func main() {
	config := parseFlags()
	lines := readLines()

	if config.unique {
		lines = uniqueLines(lines)
	}

	sort.Slice(lines, func(i, j int) bool {
		a := getColumn(lines[i], config.column, config.delimiter)
		b := getColumn(lines[j], config.column, config.delimiter)
		if config.numeric {
			return compareNumeric(a, b, config.reverse)
		}
		if config.reverse {
			return b < a
		}
		return a < b
	})

	for _, line := range lines {
		fmt.Println(line)
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

func getColumn(line string, column int, delimiter string) string {
	fields := strings.Split(line, delimiter)
	if column-1 < len(fields) {
		return fields[column-1]
	}
	return ""
}

func compareNumeric(a, b string, reverse bool) bool {
	aNum, aErr := strconv.ParseFloat(a, 64)
	bNum, bErr := strconv.ParseFloat(b, 64)

	if aErr != nil && bErr != nil {
		if reverse {
			return b < a
		}
		return a < b
	}
	if aErr != nil {
		return false
	}
	if bErr != nil {
		return true
	}
	if reverse {
		return bNum < aNum
	}
	return aNum < bNum
}

func uniqueLines(lines []string) []string {
	uniqueMap := make(map[string]bool)
	var uniqueLines []string
	for _, line := range lines {
		if !uniqueMap[line] {
			uniqueMap[line] = true
			uniqueLines = append(uniqueLines, line)
		}
	}
	return uniqueLines
}
