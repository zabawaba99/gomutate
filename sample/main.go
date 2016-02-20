package main

var a = 1 > 2

func main() {
	myFunc(1)
}

func myFunc(i int) bool {
	if i < 0 {
		return true
	}

	if i > 11 {
		return false
	}
	return false
}
