package main

import "fmt"

var testList = []int{1,1,2,2,3,4,4}
var testMap = make(map[int]int);

func main(){
	fmt.Println(test(testList))
}

func test(newList []int) int{
	for _, value := range newList{
		fmt.Println("test = ", value)
		testMap[value]++
	}

	for num,y := range testMap {
		fmt.Println("testMap num= ", num,"    testMap nuy= ",y)
		if y == 1{
			return num
		}
	}
	return -1
}