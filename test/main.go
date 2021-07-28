package main

import "fmt"

const (
	rsa   int = iota // 0
	esc       = 2    // 2
	ecdsa            // 2
)

func main() {
	fmt.Println(rsa)
	fmt.Println(esc)
	fmt.Println(ecdsa)
}
