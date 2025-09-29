package main

import (
	"fmt"
)

func main(){
	fmt.Println("test =: ",temp("[])())"))
}

func temp (test string) bool{
	 n := len(test)
	if n % 2 == 1 {
		return false
	}

	tempMap := map[byte]byte{
		')':'(',
		']':'[',
		'}':'{',
	}

	stack := make([]byte, 0)

	for i := 0; i < len(test); i++ {
		char := test[i]
		fmt.Println("temp char=: ",char)
		
		if matching,isRight := tempMap[char];isRight {
			fmt.Println("tempMap =: ",matching," isRight = ",isRight)
			if len(stack) == 0 || stack[len(stack) -1] != matching{
				return false
			}
			stack = stack[:len(stack)-1]
		}else{
			fmt.Println("append =: ",char)
			stack = append(stack,char)
		}
	}

	return len(stack) == 0
}