package main

import (
	"fmt"
	"sync"
)

func main() {
	testSync()
}

func testSync() {
	var counter int       // 共享计数器
	var mu sync.Mutex     // 互斥锁，用于保护计数器
	var wg sync.WaitGroup // 等待组，用于等待所有协程完成

	for i := 0; i < 10; i++ {
		wg.Add(1) // 为每个协程增加等待计数
		go func(id int) {
			defer wg.Done() // 协程结束时减少等待计数
			for j := 0; j < 1000; j++ {
				fmt.Println("counter : ", counter)
				mu.Lock()   // 获取锁
				counter++   // 在锁的保护下递增计数器
				mu.Unlock() // 释放锁
				fmt.Println("counter 释放锁i: ", counter)
			}
		}(i)
	}

	wg.Wait() // 等待所有协程完成

	// 输出最终结果
	fmt.Printf("\nFinal counter value: %d\n", counter)
	fmt.Printf("Correct: %t\n", counter == 10000)
}