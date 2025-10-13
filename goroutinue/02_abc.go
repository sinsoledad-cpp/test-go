package main

import (
	"fmt"
	"sync"
)

// 三个goroutinue交替打印abc 10次

func main() {
	cha := make(chan struct{})
	chb := make(chan struct{})
	chc := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			<-chc
			fmt.Printf("第%02d次打印\n", i)
			fmt.Println("a")
			cha <- struct{}{}
		}
		<-chc
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			<-cha
			fmt.Println("b")
			chb <- struct{}{}
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			<-chb
			fmt.Println("c")
			chc <- struct{}{}
		}
	}()
	chc <- struct{}{}
	wg.Wait()
	close(cha)
	close(chb)
	close(chc)
}
