package util

import "testing"

func TestGenerateRandomId(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateRandomId(); got != tt.want {
				t.Errorf("GenerateRandomId() = %v, want %v", got, tt.want)
			}
		})
	}
}
