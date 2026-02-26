package models

import "testing"

func TestParseTaskID(t *testing.T) {
	tests := []struct {
		input   string
		want    int64
		wantErr bool
	}{
		{"19", 19, false},
		{"1", 1, false},
		{"GHST-19", 19, false},
		{"ghst-19", 19, false},
		{"Ghst-5", 5, false},
		{"  GHST-42  ", 42, false},
		{"  7  ", 7, false},
		{"abc", 0, true},
		{"GHST-", 0, true},
		{"GHST-abc", 0, true},
		{"", 0, true},
		{"FOO-19", 0, true},
	}

	for _, tt := range tests {
		got, err := ParseTaskID(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseTaskID(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("ParseTaskID(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}
