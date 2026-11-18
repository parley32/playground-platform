package main

import (
	"errors"
	"fmt"
	"os/exec"
	"testing"
)

func TestExitCode(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{
			name: "nil error",
			err:  nil,
			want: 0,
		}, {
			name: "child exit error",
			err:  exec.Command("sh", "-c", "exit 3").Run(),
			want: 3,
		}, {
			name: "child exit error, wrapped",
			err:  fmt.Errorf("wait: %w", exec.Command("sh", "-c", "exit 3").Run()),
			want: 3,
		}, {
			name: "not found",
			err:  fmt.Errorf("wait: %w", exec.ErrNotFound),
			want: 127,
		}, {
			name: "usage error",
			err:  errUsage,
			want: 2,
		}, {
			name: "general error",
			err:  errors.New("boom"),
			want: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := exitCode(test.err)
			if test.want != got {
				t.Fatalf("wanted %d, got %d", test.want, got)
			}
		})
	}
}
