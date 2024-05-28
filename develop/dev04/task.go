package main

import (
	"fmt"
	"sort"
	"strings"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func Anagram(list *[]string) *map[string]*[]string {
	// Создаем мапу для хранения множеств анаграмм
	anagramMap := make(map[string][]string)

	// Если не делать проверку на уникальность строк после приведения к нижнему регистру,
	// то слова "Куб" и "куб" будут дублироваться в результирующей мапе.

	unique := make(map[string]struct{})
	lowerList := make([]string, 0, len(*list))

	for _, word := range *list {
		lowerWord := strings.ToLower(word)
		if _, ok := unique[lowerWord]; !ok {
			lowerList = append(lowerList, lowerWord)
			unique[lowerWord] = struct{}{}
		}
	}

	// Сортируем символы внутри слова
	for _, word := range lowerList {
		sortedWord := sortString(word)
		anagramMap[sortedWord] = append(anagramMap[sortedWord], word)
	}

	// Создаем результативную мапу
	result := make(map[string]*[]string)

	for _, words := range anagramMap {
		if len(words) > 1 {
			sort.Strings(words)
			result[words[0]] = &words
		}
	}

	return &result
}

// sortString сортирует символы в строке
func sortString(s string) string {
	runes := []rune(s)
	sort.Slice(runes, func(i, j int) bool {
		return runes[i] < runes[j]
	})
	return string(runes)
}
func main() {
	words := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стопка", "Кастоп", "кастоп"}
	anagrams := Anagram(&words)

	for key, val := range *anagrams {
		fmt.Printf("%s: %v\n", key, *val)
	}
}
