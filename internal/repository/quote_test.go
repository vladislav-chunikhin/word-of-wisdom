package repository

import (
	"testing"
)

func TestGetQuote(t *testing.T) {
	tests := []struct {
		name     string
		quotes   []string
		expected bool
	}{
		{
			name:     "Single quote",
			quotes:   []string{"The only limit to our realization of tomorrow is our doubts of today."},
			expected: true,
		},
		{
			name: "Multiple quotes",
			quotes: []string{
				"The only limit to our realization of tomorrow is our doubts of today.",
				"Do not wait to strike till the iron is hot; but make it hot by striking.",
			},
			expected: true,
		},
		{
			name:     "No quotes",
			quotes:   []string{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Quote{quotes: tt.quotes}
			quote, err := q.GetQuote()

			if tt.expected && err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if tt.expected && quote == "" {
				t.Errorf("expected a quote, got an empty string")
			}

			if !tt.expected && err == nil {
				t.Errorf("expected an error, got none")
			}
		})
	}
}
