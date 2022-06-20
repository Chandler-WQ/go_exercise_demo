package main

import (
	"testing"
	"time"
)

func TestDate(t *testing.T) {
	ti := time.Date(2047, 47, 96, 123, 98, 81, 999434, time.UTC)
	d := (1 * time.Hour)

	// Calling Truncate() method
	trunc := ti.Truncate(d)
	t.Logf("%v", trunc)

	round := ti.Round(d)
	t.Logf("%v", round)
}
