package main

import "testing"

func TestMyFuncA(t *testing.T) {
	if myFunc(1) {
		t.FailNow()
	}
}

func TestMyFuncB(t *testing.T) {
	if !myFunc(-1) {
		t.FailNow()
	}
}
