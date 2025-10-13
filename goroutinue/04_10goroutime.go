package main

import (
	"fmt"
	"sync"
)

// 用不超过10个goroutine不重复的打印slice中的100个元素
// 创建10个无缓冲channel和10个goroutine
// 固定10个goroutine,且顺序打印

func main() {
	// 准备100个元素
	s := make([]int, 100)
	for i := 0; i < 100; i++ {
		s[i] = i
	}

	var wg sync.WaitGroup
	m := make(map[int]chan int, 10)
	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		idx := i
		m[idx] = make(chan int)
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			for v := range m[idx] {
				fmt.Printf("goroutine %d print %d\n", idx, v)
				done <- struct{}{} // 打印完通知主协程
			}
		}(idx)
	}

	for i := 0; i < 100; i++ {
		idx := i % 10
		m[idx] <- s[i]
		<-done
	}
	close(done)
	for i := 0; i < 10; i++ {
		close(m[i])
	}
	wg.Wait()
}
