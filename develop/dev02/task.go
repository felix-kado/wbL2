package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func UnpackString(s string) (string, error) {
	// Обозначу, что, в целом, можно заранее алоцировать память
	// можно первый раз пройтись сложить все числа + одиночные символы, и методом grow алоцировать
	// но тут зависит от данных, если у нас часто будут строки вида a1000b2000x2222x4544x222223v3333
	// (где конкатенация пачки символов почти всегда будет приводить к аллокации) то да,
	// конечно пройтись пару раз чем делать на каждый символ алокацию

	var result strings.Builder
	var escaped bool
	runes := []rune(s)

	i, j := 0, 1

	if len(runes) > 0 && unicode.IsDigit(runes[0]) {
		return "", fmt.Errorf("start with non escape number")
	}
	if len(runes) > 0 && runes[len(runes)-1] == '\\' &&
		(len(runes) == 1 || runes[len(runes)-2] != '\\') {
		return "", fmt.Errorf("end with one backslash")
	}

	for i < len(runes) {
		char := runes[i]
		// Обработка escape последовательностей
		if char == '\\' && !escaped {
			escaped = true
			i++
			if i < len(runes) {
				j = i + 1
			}
			continue
		}

		// Если нашли число
		if j < len(runes) && unicode.IsDigit(runes[j]) {
			// для определения грацины многоразрядоного числа
			for j < len(runes) && unicode.IsDigit(runes[j]) {
				j++
			}
			count, err := strconv.Atoi(string(runes[i+1 : j]))
			if err != nil {
				return "", err
			}
			result.WriteString(strings.Repeat(string(char), count))
		} else {
			result.WriteRune(char)
		}
		i = j
		j = i + 1
		escaped = false
	}

	// Обработка последнего символа
	if i == len(runes)-1 && !unicode.IsDigit(runes[i]) {

		result.WriteRune(runes[i])
	}

	return result.String(), nil
}

func main() {
	str := `a12f2\13\\4abc\\5`
	unpackedString, err := UnpackString(str)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Unpacked string:", unpackedString)
	}
}
