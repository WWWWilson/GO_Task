package main

func main() {
}

func test(nums []int, target int) []int {
	numMap := make(map[int]int)

	for i, num := range nums {
		complement := target - num
		if idx, exist := numMap[complement]; exist {
			return []int{idx, i}
		}

		numMap[num] = i
	}
	return nil
}