package main

import (
	"fmt"
	"time"
)

func main() {
	testChannel()
	time.Sleep(1 * time.Second)
}

func testChannel() {
	ch := make(chan int, 10)

	go send(ch)
	go accept8(ch)

}

func send(ch1 chan int) {
	for i := 0; i < 100; i++ {
		fmt.Println("send : ", i)
		ch1 <- i
	}
}

func accept8(ch1 chan int) {
	for v := range ch1 {
		fmt.Println("accept8 : ", v)
	}
}