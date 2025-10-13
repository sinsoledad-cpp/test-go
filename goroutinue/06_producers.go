package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// 使用go实现1000个并发控制并设置执行超时时间1秒
// 创建 1000 个协程，并且进行打印
// 总共超时时间 1s，1s 没执行完就超时，使用 ctx 进行控制

func main() {
	tasks := make(chan int, 1000)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		idx := i
		wg.Add(1)
		tasks <- i
		go func(idx int) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				fmt.Printf("goroutine id: %d\n", idx)
			}
		}(idx)
	}
	wg.Wait()
	close(tasks)

}
