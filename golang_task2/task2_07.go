package main

import (
	"fmt"
	"time"
)

func main()  {
	channel1()
	
}

func channel1(){
    ch1 := make(chan int)
	go sendOneToTen(ch1)
	go accept(ch1)
	time.Sleep(1 * time.Second)
}


func sendOneToTen(ch chan int)  {
	for i := 1; i < 11; i++ {
		fmt.Println("sendOneToTen : " , i)
		ch <- i;
	}
	close(ch)
}

func accept(ch chan int)  {
	for v := range ch {
		fmt.Println("accept : " ,v)
	}
}