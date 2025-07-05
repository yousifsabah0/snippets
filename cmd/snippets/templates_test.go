package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	tm := time.Date(2025, 7, 5, 2, 15, 0, 0, time.UTC)
	hd := humanDate(tm)

	want := "05 Jul 2025 at 02:15"
	if hd != want {
		t.Errorf("got %q; want %q", hd, want)
	}
}
