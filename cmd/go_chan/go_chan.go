package main

import "fmt"

func main() {
	go fmt.Println("goroutine")
	fmt.Println("main")
}
