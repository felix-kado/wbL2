package main

import (
	"testing"
)

func TestUnpackString(t *testing.T) {
	tests := []struct {
		input string
		want  string
		err   bool
	}{
		{`a4bc2d5e`, `aaaabccddddde`, false}, // просто строка без escape
		{`abcd`, `abcd`, false},              // без цифр
		{`45`, ``, true},                     // только цифры
		{`a12da`, `aaaaaaaaaaaada`, false},   // двузначное чилсло
		{`a12`, `aaaaaaaaaaaa`, false},       // двузначное чилсло в конце строки
		{``, ``, false},                      // ничего
		{`qwe\4\5aaa`, `qwe45aaa`, false},    // экранированные символы подряд без повторов
		{`qwe\45\12`, `qwe4444411`, false},   // экранированные символы подряд c повторениями
		{`qwe\\5`, `qwe\\\\\`, false},        // экранированный слеш с повторами
		{`qwe\\`, `qwe\`, false},             // экранированный слеш без повторов
		{`abc\`, ``, true},                   // конец строки экранирует ничто
		{`\`, ``, true},                      // начало строки экранирует ничто
		{`a12f2\13\\4abc\\5\\`, `aaaaaaaaaaaaff111\\\\abc\\\\\\`, false}, // все
	}

	for _, tt := range tests {
		got, err := UnpackString(tt.input)
		if (err != nil) != tt.err {
			t.Errorf(`UnpackString(%q) error = %v, wantErr %v`, tt.input, err, tt.err)
		}
		if got != tt.want {
			t.Errorf(`UnpackString(%q) = %q, want %q`, tt.input, got, tt.want)
		}
	}
}
