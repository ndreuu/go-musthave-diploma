package service

import "testing"

func TestIsValidLuhn(t *testing.T) {
	tests := []struct {
		name   string
		number string
		want   bool
	}{
		{
			name:   "valid order number",
			number: "9278923470",
			want:   true,
		},
		{
			name:   "invalid checksum",
			number: "9278923471",
			want:   false,
		},
		{
			name:   "empty string",
			number: "",
			want:   false,
		},
		{
			name:   "contains letters",
			number: "92789abc70",
			want:   false,
		},
		{
			name:   "another valid number",
			number: "12345678903",
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidLuhn(tt.number)
			if got != tt.want {
				t.Fatalf("IsValidLuhn(%q) = %v, want %v", tt.number, got, tt.want)
			}
		})
	}
}
