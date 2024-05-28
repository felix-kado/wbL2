package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestMyCut(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		fields     string
		delimiter  string
		separated  bool
		wantOutput string
	}{
		{
			name:       "DefaultDelimiter",
			input:      "col1\tcol2\tcol3\nval1\tval2\tval3",
			fields:     "1,3",
			delimiter:  "\t",
			separated:  false,
			wantOutput: "col1\tcol3\nval1\tval3\n",
		},
		{
			name:       "CustomDelimiter",
			input:      "col1,col2,col3\nval1,val2,val3",
			fields:     "2",
			delimiter:  ",",
			separated:  false,
			wantOutput: "col2\nval2\n",
		},
		{
			name:       "SeparatedOnly",
			input:      "col1\tcol2\tcol3\nval1\tval2\tval3\nnodelimiterline",
			fields:     "1,3",
			delimiter:  "\t",
			separated:  true,
			wantOutput: "col1\tcol3\nval1\tval3\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("./mycut", "-f", tt.fields, "-d", tt.delimiter)
			if tt.separated {
				cmd.Args = append(cmd.Args, "-s")
			}

			// Подготовка ввода
			cmd.Stdin = strings.NewReader(tt.input)

			// Захват вывода
			var out bytes.Buffer
			cmd.Stdout = &out

			err := cmd.Run()
			if err != nil {
				t.Fatalf("Ошибка выполнения команды: %v", err)
			}

			// Проверка результата
			gotOutput := out.String()
			if gotOutput != tt.wantOutput {
				t.Errorf("Ожидаемый вывод: %q, полученный вывод: %q", tt.wantOutput, gotOutput)
			}
		})
	}
}

func TestMissingFieldsFlag(t *testing.T) {
	cmd := exec.Command("./mycut", "-d", "\t")
	cmd.Stdin = strings.NewReader("col1\tcol2\tcol3\nval1\tval2\tval3")

	err := cmd.Run()
	if err == nil {
		t.Fatal("Ожидалась ошибка при отсутствии флага -f, но её не произошло")
	}
}

func TestInvalidFieldIndex(t *testing.T) {
	cmd := exec.Command("./mycut", "-f", "10", "-d", "\t")
	cmd.Stdin = strings.NewReader("col1\tcol2\tcol3\nval1\tval2\tval3")

	// Захват вывода
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		t.Fatalf("Ошибка выполнения команды: %v", err)
	}

	// Проверка результата
	gotOutput := out.String()
	wantOutput := "\n\n"
	if gotOutput != wantOutput {
		t.Errorf("Ожидаемый вывод: %q, полученный вывод: %q", wantOutput, gotOutput)
	}
}
