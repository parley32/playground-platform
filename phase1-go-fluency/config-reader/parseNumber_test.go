package main

import (
	"testing"
)

func TestParseNumber(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		want    int
		wantErr bool
	}{
		{
			name: "parsing correct number",
			text: "42",
			want: 42,
		},
		{
			name: "parsing correct number with whitespaces",
			text: "	42	",
			want: 42,
		},
		{
			name: "parsing empty number returns 0",
			text: "",
			want: 0,
		},
		{
			name: "parsing correct negative number",
			text: "-42",
			want: -42,
		},
		{
			name: "parsing negative number with symbol noise number",
			text: "+-42",
			want: -42,
		},
		{
			name: "parsing positive number with symbol noise",
			text: "-+42",
			want: 42,
		},
		{
			name:    "parsing incorrect number",
			text:    "4-2",
			want:    0,
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := parseNumber([]byte(test.text))
			if (err != nil) != test.wantErr {
				t.Fatalf("parseNumber() error = %v, wantErr %v", err, test.wantErr)
			}
			if result != test.want {
				t.Errorf("parseNumber() = %d, want %d", result, test.want)
			}
		})
	}
}
