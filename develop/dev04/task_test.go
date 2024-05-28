package main

import (
	"reflect"
	"testing"
)

// TestAnagram tests the Anagram function
func TestAnagram(t *testing.T) {
	tests := []struct {
		input    []string
		expected map[string][]string
	}{
		{
			input: []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик"},
			expected: map[string][]string{
				"пятак":  {"пятак", "пятка", "тяпка"},
				"листок": {"листок", "слиток", "столик"},
			},
		},
		{
			input: []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "мяу"},
			expected: map[string][]string{
				"пятак":  {"пятак", "пятка", "тяпка"},
				"листок": {"листок", "слиток", "столик"},
			},
		},
		{
			input: []string{"стопка", "кастоп", "кастоп"},
			expected: map[string][]string{
				"кастоп": {"кастоп", "стопка"},
			},
		},
		{
			input: []string{"ПЯТАК", "пятак", "ПЯТКА", "тяпка"},
			expected: map[string][]string{
				"пятак": {"пятак", "пятка", "тяпка"},
			},
		},
		{
			input:    []string{"один", "два", "три"},
			expected: map[string][]string{},
		},
	}

	for _, test := range tests {
		result := Anagram(&test.input)
		for key, val := range *result {
			if !reflect.DeepEqual(*val, test.expected[key]) {
				t.Errorf("For key %s, expected %v, but got %v", key, test.expected[key], *val)
			}
		}
	}
}
