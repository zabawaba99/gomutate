package main

import "testing"

func TestBoundaryMutation1(t *testing.T) {
	if !boundaryMutation1("foo") {
		t.FailNow()
	}
}
