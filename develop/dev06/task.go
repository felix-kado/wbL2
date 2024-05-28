package main

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

type Config struct {
	Fields    []int
	Delimiter string
	Separated bool
}

func parseFlags() Config {
	var fields []int
	var delimiter string
	var separated bool

	pflag.IntSliceVarP(&fields, "fields", "f", nil, "выбрать поля (колонки)")
	pflag.StringVarP(&delimiter, "delimiter", "d", "\t", "использовать другой разделитель")
	pflag.BoolVarP(&separated, "separated", "s", false, "только строки с разделителем")

	pflag.Parse()

	if len(fields) == 0 {
		fmt.Fprintln(os.Stderr, "Необходимо указать поля для выбора с помощью флага -f")
		os.Exit(1)
	}

	return Config{
		Fields:    fields,
		Delimiter: delimiter,
		Separated: separated,
	}
}

func main() {
	config := parseFlags()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, config.Delimiter)
		if config.Separated && len(parts) == 1 {
			continue
		}

		var result []string
		for _, idx := range config.Fields {
			idx-- // Преобразование к нулевой индексации
			if idx >= 0 && idx < len(parts) {
				result = append(result, parts[idx])
			}
		}
		fmt.Println(strings.Join(result, config.Delimiter))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка чтения:", err)
	}
}
