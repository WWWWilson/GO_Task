package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	testAtomic()
}

func testAtomic() {
	var counter int32     // 共享计数器
	var wg sync.WaitGroup // 等待组，用于等待所有协程完成

	for i := 0; i < 10; i++ {
		wg.Add(1) // 为每个协程增加等待计数
		go func(id int) {
			defer wg.Done() // 协程结束时减少等待计数
			for j := 0; j < 1000; j++ {
				fmt.Println("counter : ", counter)
				atomic.AddInt32(&counter, 1)
				fmt.Println("counter 释放锁i: ", counter)
			}
		}(i)
	}

	wg.Wait() // 等待所有协程完成

	// 输出最终结果
	finalValue := atomic.LoadInt32(&counter)
	fmt.Printf("\nFinal counter value: %d\n", finalValue)
	fmt.Printf("Correct: %t\n", finalValue == 10000)
}