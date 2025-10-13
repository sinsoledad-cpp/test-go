package main

import (
	"fmt"
	"sync"
)

// 用不超过10个goroutine不重复的打印slice中的100个元素
// 容量为10的有缓冲channel实现
// 每次启动10个，累计启动100个goroutine,且无序打印

func main() {
	var wg sync.WaitGroup
	// 准备数据
	s := make([]int, 100)
	for i := 0; i < 100; i++ {
		s[i] = i
	}
	ch := make(chan struct{}, 10)
	// 启动100个goroutine
	for i := 0; i < 100; i++ {
		wg.Add(1)
		ch <- struct{}{}
		idx := i
		go func(idx int) {
			defer wg.Done()
			fmt.Println(s[idx])
			<-ch // 释放一个槽位
		}(idx)
	}
	close(ch)
	wg.Wait()
}
