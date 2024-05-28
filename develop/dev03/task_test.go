package main

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestGetColumn(t *testing.T) {
	tests := []struct {
		line      string
		column    int
		delimiter string
		expected  string
	}{
		{"1,dog,200", 2, ",", "dog"},
		{"1,dog,200", 1, ",", "1"},
		{"1 dog 200", 2, " ", "dog"},
		{"1 dog 200", 3, " ", "200"},
		{"1,dog,200", 4, ",", ""},
	}

	for _, test := range tests {
		result := getColumn(test.line, test.column, test.delimiter)
		if result != test.expected {
			t.Errorf("getColumn(%q, %d, %q) = %q; want %q", test.line, test.column, test.delimiter, result, test.expected)
		}
	}
}

func TestCompareNumeric(t *testing.T) {
	tests := []struct {
		a        string
		b        string
		reverse  bool
		expected bool
	}{
		{"10", "20", false, true},
		{"20", "10", false, false},
		{"10", "20", true, false},
		{"20", "10", true, true},
		{"a", "b", false, true},
		{"a", "b", true, false},
	}

	for _, test := range tests {
		result := compareNumeric(test.a, test.b, test.reverse)
		if result != test.expected {
			t.Errorf("compareNumeric(%q, %q, %v) = %v; want %v", test.a, test.b, test.reverse, result, test.expected)
		}
	}
}

func TestUniqueLines(t *testing.T) {
	tests := []struct {
		lines    []string
		expected []string
	}{
		{
			[]string{"1,cat", "2,dog", "1,cat"},
			[]string{"1,cat", "2,dog"},
		},
		{
			[]string{"1,cat", "2,dog", "3,fish"},
			[]string{"1,cat", "2,dog", "3,fish"},
		},
	}

	for _, test := range tests {
		result := uniqueLines(test.lines)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("uniqueLines(%v) = %v; want %v", test.lines, result, test.expected)
		}
	}
}

func TestSortLines(t *testing.T) {
	tests := []struct {
		input    string
		config   sortConfig
		expected []string
	}{
		{
			"3,cat,100\n1,dog,200\n2,fish,50\n4,bird,150\n5,cat,75",
			sortConfig{column: 1, numeric: true, reverse: false, unique: false, delimiter: ","},
			[]string{"1,dog,200", "2,fish,50", "3,cat,100", "4,bird,150", "5,cat,75"},
		},
		{
			"3,cat,100\n1,dog,200\n2,fish,50\n4,bird,150\n5,cat,75",
			sortConfig{column: 2, numeric: false, reverse: false, unique: false, delimiter: ","},
			[]string{"4,bird,150", "3,cat,100", "5,cat,75", "1,dog,200", "2,fish,50"},
		},
		{
			"3,cat,100\n1,dog,200\n2,fish,50\n4,bird,150\n5,cat,75",
			sortConfig{column: 3, numeric: true, reverse: true, unique: false, delimiter: ","},
			[]string{"1,dog,200", "4,bird,150", "3,cat,100", "5,cat,75", "2,fish,50"},
		},
	}

	for _, test := range tests {
		lines := strings.Split(test.input, "\n")
		if test.config.unique {
			lines = uniqueLines(lines)
		}

		sort.Slice(lines, func(i, j int) bool {
			a := getColumn(lines[i], test.config.column, test.config.delimiter)
			b := getColumn(lines[j], test.config.column, test.config.delimiter)
			if test.config.numeric {
				return compareNumeric(a, b, test.config.reverse)
			}
			if test.config.reverse {
				return b < a
			}
			return a < b
		})

		if !reflect.DeepEqual(lines, test.expected) {
			t.Errorf("sortLines(%v, %v) = %v; want %v", test.input, test.config, lines, test.expected)
		}
	}
}
