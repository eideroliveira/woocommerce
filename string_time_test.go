package woocommerce

import (
	"encoding/json"
	"testing"
	"time"
)

func TestStringTime_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
	}{
		{
			name:    "RFC3339",
			json:    `"2026-01-19T10:00:00Z"`,
			wantErr: false,
		},
		{
			name:    "YYYY-MM-DDTHH:MM:SS",
			json:    `"2026-01-19T10:00:00"`,
			wantErr: false,
		},
		{
			name:    "YYYY-MM-DD",
			json:    `"2026-01-19"`,
			wantErr: false,
		},
		{
			name:    "DD/MM/YYYY",
			json:    `"19/01/2026"`,
			wantErr: false,
		},
		{
			name:    "DD-MM-YYYY",
			json:    `"19-01-2026"`,
			wantErr: true, // Currently fails?
		},
		{
			name:    "YYYY/MM/DD",
			json:    `"2026/01/19"`,
			wantErr: true, // Currently fails?
		},
        {
            name: "DD/MM/YYYY HH:MM",
            json: `"19/01/2026 10:00"`,
            wantErr: true,
        },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var st StringTime
			err := json.Unmarshal([]byte(tt.json), &st)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringTime.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				parsed := time.Time(st)
				if parsed.IsZero() {
					t.Error("Parsed time is zero")
				}
				// Verify sanity of parsed date for 19/01/2026
				if tt.json == `"19/01/2026"` || tt.json == `"2026-01-19"` {
					if parsed.Year() != 2026 || parsed.Month() != 1 || parsed.Day() != 19 {
						t.Errorf("Parsed date mismatch: got %v", parsed)
					}
				}
			}
		})
	}
}
