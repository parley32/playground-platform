package main

import (
	"bytes"
	"context"
	"io"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name    string
		inputs  []string
		wantErr bool
	}{
		{
			name:   "regular inputs",
			inputs: []string{"workerPool", "2", "5"},
		},
		{
			name:    "zero workers",
			inputs:  []string{"workerPool", "0", "5"},
			wantErr: true,
		},
		{
			name:    "negative jobs",
			inputs:  []string{"workerPool", "3", "-5"},
			wantErr: true,
		},
		{
			name:   "more workers than jobs",
			inputs: []string{"workerPool", "5", "3"},
		},
		{
			name:    "invalid worker number",
			inputs:  []string{"workerPool", "abab", "3"},
			wantErr: true,
		},
		{
			name:    "missing input",
			inputs:  []string{"workerPool", "2"},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := run(test.inputs, io.Discard); (err != nil) != test.wantErr {
				t.Fatalf("run() returns err = %v, want %v", err, test.wantErr)
			}
		})
	}
}

func TestRunResults(t *testing.T) {
	var output bytes.Buffer
	want := []string{"Result: 1", "Result: 2", "Result: 3", "Result: 4", "Result: 5"}

	if err := run([]string{"workerPool", "2", "5"}, &output); err != nil {
		t.Fatalf("run() returned an unexpected error: %v", err)
	}

	got := strings.Split(strings.TrimSpace(output.String()), "\n")
	sort.Strings(got)

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("run() results = %v, want %v", got, want)
	}
}

func TestGetReadyForWorkStopsOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	jobs := make(chan int)
	results := make(chan int, 1)

	done := make(chan struct{})
	go func() {
		defer close(done)
		getReadyForWork(ctx, jobs, results)
	}()

	cancel()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("worker did not exit after cancel")
	}
}
