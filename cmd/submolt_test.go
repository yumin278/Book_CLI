package cmd

import "testing"

func TestStripSubmoltPrefix(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "no prefix", input: "memory", want: "memory"},
		{name: "with m/ prefix", input: "m/memory", want: "memory"},
		{name: "multi-word no prefix", input: "book-club", want: "book-club"},
		{name: "multi-word with prefix", input: "m/book-club", want: "book-club"},
		{name: "empty string", input: "", want: ""},
		{name: "prefix only", input: "m/", want: ""},
		{name: "unrelated prefix is preserved", input: "x/memory", want: "x/memory"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripSubmoltPrefix(tt.input)
			if got != tt.want {
				t.Errorf("stripSubmoltPrefix(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
