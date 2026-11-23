package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		want    []int
		wantErr bool
	}{
		{
			name:  "Valid input",
			input: []string{"2", "4", "10", "-10"},
			want:  []int{2, 4, 10, -10},
		},
		{
			name:    "Invalid input with non-digit characters",
			input:   []string{"2", "abab", "10", "-10"},
			wantErr: true,
		},
		{
			name:  "Valid empty input",
			input: []string{},
			want:  []int(nil),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parseArgs(test.input)
			if (err != nil) != test.wantErr {
				t.Fatalf("parseArgs() returns err = %v, want %v", err, test.wantErr)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("parseArgs() = %#v, want %#v", got, test.want)
			}
		})
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		name       string
		input      []string
		wantOutput string
		wantErr    bool
	}{
		{
			name:       "sums valid input",
			input:      []string{"2", "4", "5"},
			wantOutput: "Sum: 11\n",
		},
		{
			name:    "rejects missing input",
			input:   []string{},
			wantErr: true,
		},
		{
			name:    "rejects invalid input",
			input:   []string{"2", "not-a-number"},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer

			err := run(test.input, &output)
			if (err != nil) != test.wantErr {
				t.Fatalf("run() error = %v, want error: %v", err, test.wantErr)
			}
			if output.String() != test.wantOutput {
				t.Errorf("run() output = %q, want %q", output.String(), test.wantOutput)
			}
		})
	}
}
