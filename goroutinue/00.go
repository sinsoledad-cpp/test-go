package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"
)

/*
1. 每秒随机生成一批（100个以内）键值对<key,value>，其中key必须是小写字母 a~z的范围，value必须是 0.1~5.1的浮点数
2. 通过chan将键值对传递给下游
3. 下游从channel消费，对<key, value>做全局统计，要求10个并发（即同时处理10个<key, value>），要求最终输出每个key的value最大值、最小值、平均值、和值以及分位值，分位值是加分项
4. 程序收到SIGHUP、SIGINT、SIGQUIT信号，停止生成器、统计器，并打印出统计结果
*/
type KV struct {
	Key   string
	Value float64
}

type Stats struct {
	Max    float64
	Min    float64
	Sum    float64
	Count  int64
	Values []float64
}

func main() {

	const numWorkers int = 10
	dataChan := make(chan KV, 1000)

	var statsMutex sync.Mutex
	statsData := make(map[string]*Stats)

	var workerWg sync.WaitGroup
	var generatorWg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	workerWg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer workerWg.Done()
			for kv := range dataChan {
				statsMutex.Lock()
				stats, ok := statsData[kv.Key]
				if !ok {
					stats = &Stats{
						Max:    kv.Value,
						Min:    kv.Value,
						Sum:    kv.Value,
						Count:  1,
						Values: []float64{kv.Value},
					}
					statsData[kv.Key] = stats
				} else {
					stats.Max = max(stats.Max, kv.Value)
					stats.Min = min(stats.Min, kv.Value)
					stats.Sum += kv.Value
					stats.Count++
					stats.Values = append(stats.Values, kv.Value)
				}
				statsMutex.Unlock()
			}
		}()
	}

	generatorWg.Add(1)
	go func() {
		defer generatorWg.Done()
		source := rand.NewSource(time.Now().UnixNano())
		r := rand.New(source)
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				batchSize := r.Intn(100) + 1
				for i := 0; i < batchSize; i++ {
					key := string(byte(r.Intn(26) + 'a'))
					value := r.Float64()*5.0 + 0.1
					dataChan <- KV{key, value}
				}
			case <-ctx.Done():
				close(dataChan)
				return
			}
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)
	<-signalChan

	cancel()
	workerWg.Wait()

	keys := make([]string, 0, len(statsData))
	for k := range statsData {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		stats := statsData[key]
		if stats.Count == 0 {
			continue
		}

		// 计算平均值
		avg := stats.Sum / float64(stats.Count)

		// 计算分位值
		sort.Float64s(stats.Values)
		n := len(stats.Values)
		q1 := stats.Values[n/4]
		median := stats.Values[n/2]
		q3 := stats.Values[n*3/4]

		fmt.Printf("Key [%s]:\n", key)
		fmt.Printf("  总和 (Sum):    %.2f\n", stats.Sum)
		fmt.Printf("  平均值 (Avg):   %.2f\n", avg)
		fmt.Printf("  最小值 (Min):   %.2f\n", stats.Min)
		fmt.Printf("  最大值 (Max):   %.2f\n", stats.Max)
		fmt.Printf("  计数值 (Count): %d\n", stats.Count)
		fmt.Printf("  分位值 (Quantiles): Q1(25%%)=%.2f, Median(50%%)=%.2f, Q3(75%%)=%.2f\n", q1, median, q3)
		fmt.Println("--------------------")
	}

}
