package main

import (
	"fmt"
	"sync"
)

//编写一个程序限制10个goroutine执行，每执行完一个goroutine就放一个新的goroutine进来

func main() {
	var wg sync.WaitGroup
	ch := make(chan struct{}, 10)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		ch <- struct{}{}
		go func() {
			defer wg.Done()
			fmt.Println("goroutine")
			<-ch
		}()
	}
	wg.Wait()
	close(ch)
}
