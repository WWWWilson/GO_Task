package main

import (
	"fmt"
	"strconv"
)

func main() {
	fmt.Println("test",test2([]int{1,2,3}))
}

func test(num []int) []int {
	n := len(num)

	for i := n-1; i >= 0; i-- {
		num[i]++
		if num[i] < 10 {
			return num
		} 
		num[i] = 0;
	}

	return append([]int{1},num...)
}

func test2 (num []int) []int{
	if len(num) == 0 {
		return []int{1}
	}

	if len(num) <= 18{
		n := int64(0)
		for i := 0; i < len(num); i++ {
			n = n*10+int64(num[i])
			fmt.Println("n : ",n*10," int64 : ",int64(num[i]))
			fmt.Println("n2 : ",n)
		}
		n++

		numStr := strconv.FormatInt(n,10);
		result := make([]int, len(numStr))
		for  i,char:= range numStr {
			result[i] = int(char - '0')
		}
		return result
	}
	return test2(num)
}