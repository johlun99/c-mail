package main

import (
	"strings"
	"testing"
)

func TestGreet(t *testing.T) {
	app := NewApp(nil)
	got := app.Greet("Johan")
	if !strings.Contains(got, "Johan") {
		t.Errorf("Greet() = %q, want it to contain the name", got)
	}
}
