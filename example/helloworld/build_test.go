package main

import "testing"

func TestShouldLookInterface(t *testing.T) {
	desc, _ := shouldLookInterface([]interface{}{1, 2}, "1")
	t.Log(desc)
}
