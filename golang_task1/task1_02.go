package main

import (
	"fmt"
	"strconv"
)

func main(){
	fmt.Println("temp = ",temp(12332))
}

func temp (x int) bool{
	 if x < 0 || (x % 10 == 0 && x != 0 ){
        return false;
    }

	s := strconv.Itoa(x)
	left,right := 0,len(s)-1
	for left < right {
		if s[left] != s[right]{
			return false
		}
		left++
		right--
	}
	return true
}
