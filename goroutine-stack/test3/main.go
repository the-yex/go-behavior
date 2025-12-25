package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	const batchSize = 10000
	const batches = 3

	fmt.Printf("=== Goroutine 复用实验 ===\n\n")

	for batch := 0; batch < batches; batch++ {
		start := runtime.NumGoroutine()

		// 创建一批 goroutine
		for i := 0; i < batchSize; i++ {
			wg.Add(1)
			go func() {
				time.Sleep(10 * time.Millisecond)
				wg.Done()
			}()
		}

		time.Sleep(20 * time.Millisecond)
		peak := runtime.NumGoroutine()

		// 等待这批 goroutine 完成
		wg.Wait()
		time.Sleep(50 * time.Millisecond)

		end := runtime.NumGoroutine()

		fmt.Printf("批次 %d:\n", batch)
		fmt.Printf("  启动前: %d goroutines\n", start)
		fmt.Printf("  峰值:   %d goroutines\n", peak)
		fmt.Printf("  结束后: %d goroutines\n", end)
		fmt.Printf("  创建数: %d\n", batchSize)
		fmt.Printf("  💡 观察：goroutine 数量回落，说明被回收进入复用池\n\n")

		runtime.GC() // 触发 GC，促进 g 结构回收
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("=== 实验结束 ===")
	fmt.Println("结论：runtime 会回收空闲的 goroutine 结构，放入复用池")
	fmt.Println("     下次创建 goroutine 时优先从复用池获取，而非总是新建")
}
