package cmd

import (
	"strings"
	"testing"
)

func TestCleanAndValidateUUID(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantCleaned string
		wantErr     bool
	}{
		{
			name:        "Valid UUID",
			input:       "550e8400-e29b-41d4-a716-446655440000",
			wantCleaned: "550e8400-e29b-41d4-a716-446655440000",
			wantErr:     false,
		},
		{
			name:        "Valid UUID with trailing comma",
			input:       "550e8400-e29b-41d4-a716-446655440000,",
			wantCleaned: "550e8400-e29b-41d4-a716-446655440000",
			wantErr:     false,
		},
		{
			name:        "Valid UUID with spaces",
			input:       "  550e8400-e29b-41d4-a716-446655440000  ",
			wantCleaned: "550e8400-e29b-41d4-a716-446655440000",
			wantErr:     false,
		},
		{
			name:        "Valid UUID with space and comma",
			input:       " 550e8400-e29b-41d4-a716-446655440000, ",
			wantCleaned: "550e8400-e29b-41d4-a716-446655440000",
			wantErr:     false,
		},
		{
			name:        "Invalid UUID with character L",
			input:       "550e8400-e29b-41d4-a716-44665544000L",
			wantCleaned: "550e8400-e29b-41d4-a716-44665544000L",
			wantErr:     true,
		},
		{
			name:        "Username instead of UUID",
			input:       "jules_engineer",
			wantCleaned: "jules_engineer",
			wantErr:     true,
		},
		{
			name:        "Invalid UUID with punctuation in middle",
			input:       "550e8400-e29b-41d4,a716-446655440000",
			wantCleaned: "550e8400-e29b-41d4,a716-446655440000",
			wantErr:     true,
		},
		{
			name:        "Empty input",
			input:       "",
			wantCleaned: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cleanAndValidateUUID(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("cleanAndValidateUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantCleaned {
				t.Errorf("cleanAndValidateUUID() got = %v, want %v", got, tt.wantCleaned)
			}
			if tt.wantErr {
				if !strings.Contains(err.Error(), tt.input) {
					t.Errorf("error message should contain original input %q, but got: %v", tt.input, err)
				}
				if !strings.Contains(err.Error(), tt.wantCleaned) {
					t.Errorf("error message should contain cleaned input %q, but got: %v", tt.wantCleaned, err)
				}
			}
		})
	}
}
