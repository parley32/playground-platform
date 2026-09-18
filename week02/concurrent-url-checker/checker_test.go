package main

import (
	"os"
	"reflect"
	"testing"
)

func TestGetURLList(t *testing.T) {
	fileName := "tempFile"

	_, err := getURLList(fileName)
	if err == nil {
		t.Fatalf("failed to return err with invalid filename")
	}

	file, err := os.CreateTemp("", fileName)
	defer os.Remove(file.Name())
	defer file.Close()

	if err != nil {
		t.Fatalf("failed to create temp file: %q\n", err)
	}

	urlList, err := getURLList(file.Name())
	if err != nil {
		t.Fatalf("getURLList() returned unexpected error: %v", err)
	}
	want := []string{}
	if !reflect.DeepEqual(urlList, want) {
		t.Fatalf("failed to return an empty list when input is empty. Got: %q, want: %q\n", urlList, want)
	}

	file.WriteString("stringA\nstringB")
	urlList, err = getURLList(file.Name())
	if err != nil {
		t.Fatalf("getURLList() returned unexpected error: %v", err)
	}
	if !reflect.DeepEqual(urlList, []string{"stringA", "stringB"}) {
		t.Fatalf("failed to return an array of file strings: %q\n", urlList)
	}
}

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "valid url",
			url:  "https://www.test.com",
			want: true,
		}, {
			name: "empty url",
			url:  "",
			want: false,
		}, {
			name: "empty host",
			url:  "https://",
			want: false,
		}, {
			name: "missing scheme",
			url:  "www.test.com",
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if result := isValidURL(test.url); result != test.want {
				t.Fatalf("got %v, want %v", result, test.want)
			}
		})
	}
}
