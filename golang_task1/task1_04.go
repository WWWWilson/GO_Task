package main

import (
	"fmt"
)

func main(){
	value := test([]string{"flower","flow","flight"})
	fmt.Println("value = " , value)
}

func test (str []string)string{
	if len(str) == 0 {
		return " "
	}
	minLen := len(str[0])
	for i := 1; i < len(str); i++ {
		if len(str[i]) < minLen{
			minLen = len(str[i])
		}
	}

	for i := 0; i < minLen; i++ {
		currentChar := str[0][i]
		fmt.Println("currentChar = " ,currentChar)
		for j := 1; j < len(str); j++ {
			if str[j][i] != currentChar {
				return str[0][:i]
			}
		}
	}
	return str[0][:minLen]
}