package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2025, 7, 5, 2, 15, 0, 0, time.UTC),
			want: "05 Jul 2025 at 02:15",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2025, 7, 5, 2, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "05 Jul 2025 at 01:15",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hd := humanDate(test.tm)
			if hd != test.want {
				t.Errorf("want %q; got %q", test.want, hd)
			}
		})
	}
}
