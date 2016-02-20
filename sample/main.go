package main

var a = 1 > 2

func main() {
	myFunc()
}

func myFunc() {
	i := 1
	if i < 0 {
		i = 1
		i++
	}

	if i > 11 {
		i = 123
		i--
	}
}
