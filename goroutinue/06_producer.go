package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	var cond sync.Cond
	cond.L = new(sync.Mutex)
	msgCh := make(chan int, 5)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rand.Seed(time.Now().UnixNano())

	// 生产者
	producer := func(ctx context.Context, out chan<- int, idx int) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				// 每次生产者退出，都唤醒一个消费者处理，防止最后有消费者线程死锁
				// 生产者比消费者多，所以cond.Signal()就可以。不然的话建议Broadcast()
				cond.Broadcast()
				fmt.Println("producer finished")
				return
			default:
				cond.L.Lock()
				for len(msgCh) == 5 {
					cond.Wait()
				}
				num := rand.Intn(500)
				out <- num
				fmt.Printf("producer: %d, msg: %d\n", idx, num)
				cond.Signal()
				cond.L.Unlock()
			}
		}
	}

	// 消费者
	consumer := func(ctx context.Context, in <-chan int, idx int) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				// 消费者可以选择继续消费直到channel为空
				for len(msgCh) > 0 {
					select {
					case num := <-in:
						fmt.Printf("consumer %d, msg: %d\n", idx, num)
					default:
						// 如果channel已经空了，跳出循环
						break
					}
				}
				fmt.Println("consumer finished")
				return
			default:
				cond.L.Lock()
				for len(msgCh) == 0 {
					cond.Wait()
				}
				num := <-in
				fmt.Printf("consumer %d, msg: %d\n", idx, num)
				cond.Signal()
				cond.L.Unlock()
			}
		}
	}

	// 启动生产者和消费者
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go producer(ctx, msgCh, i+1)
	}
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go consumer(ctx, msgCh, i+1)
	}

	// 模拟程序运行一段时间
	wg.Wait()
	close(msgCh)
	fmt.Println("all finished")
}
