package main

import (
	"fmt"
)

// 同时开启100个协程(分别为1号协程 2号协程 ... 100号协程，
// 1号协程只打印尾数为1的数字，2号协程只打印尾数为2的数，以此类推)
//	，请顺序打印1-1000整数以及对应的协程号；

func main() {
	done := make(chan struct{})
	m := make(map[int]chan int, 100)

	for i := 0; i < 100; i++ {
		m[i] = make(chan int)
	}

	for i := 0; i < 100; i++ {
		go func(i int) {
			for {
				v, ok := <-m[i]
				if !ok {
					return // 通道关闭，退出协程
				}
				fmt.Printf("协程%02d 打印:%d\n", i, v)
				done <- struct{}{} // 打印完通知主协程
			}
		}(i)
	}

	for i := 1; i <= 1000; i++ {
		idx := i % 100
		m[idx] <- i
		<-done // 等待它打印完，再发送下一个
	}
	//time.Sleep(1 * time.Second)
	// 关闭所有通道
	for i := 0; i < 100; i++ {
		close(m[i])
	}
}
