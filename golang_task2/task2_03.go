package main

import (
	"fmt"
	"time"
)

func main() {
fmt.Println("go run: ")
	go test1()
	go test2()
	time.Sleep(2 * time.Second)
}

func test1(){
	for i := 2; i < 10; i+=2 {
			fmt.Println("test1: ", i)
			 time.Sleep(100 * time.Millisecond)
		}
}

func test2(){
	for i := 1; i < 10; i+=2 {
			// fmt.Println("test2: ", i)
			 time.Sleep(100 * time.Millisecond)
		}
}