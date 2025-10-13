package main

import (
	"fmt"
	"sync"
)

//使用两个Goroutine，向标准输出中按顺序按顺序交替打出字母与数字，输出是a1b2c3

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	strch := make(chan struct{})
	numch := make(chan struct{})

	go func() {
		defer wg.Done()
		for i := 'a'; i <= 'z'; i++ {
			<-numch
			fmt.Printf("%c", i)
			strch <- struct{}{}
		}
	}()

	go func() {
		defer wg.Done()
		for i := 1; i <= 26; i++ {
			<-strch
			fmt.Printf("%d", i)
			numch <- struct{}{}
		}
	}()

	numch <- struct{}{}

	wg.Wait()
	close(strch)
	close(numch)

}
