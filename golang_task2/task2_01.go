package main

import "fmt"


func main() {
	var num int = 1
	test(&num)
	fmt.Println("test: " ,num)
}

func test(numTest *int) { 
	 *numTest = *numTest + 10
}