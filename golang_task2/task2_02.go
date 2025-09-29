package main

import "fmt"

func main() {
	var sli = make([]int, 3, 3)
	sli[0] = 1;
	sli[1] = 2;
	sli[2] = 3;
	test(&sli)
	fmt.Println("test = " ,sli)
}

func test(nums *[]int) {
	for i := 0; i < len(*nums); i++ {
		(*nums)[i] *= 2
	}

}