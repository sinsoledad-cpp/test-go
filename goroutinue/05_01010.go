package main

import (
	"fmt"
	"time"
)

//用单个channel实现0,1的交替打印

func main() {
	done := make(chan struct{})
	go func() {
		for {
			<-done
			fmt.Println(0)
			done <- struct{}{}
		}
	}()
	go func() {
		for {
			<-done
			fmt.Println(1)
			done <- struct{}{}
		}
	}()
	done <- struct{}{}
	time.Sleep(time.Second)
}
