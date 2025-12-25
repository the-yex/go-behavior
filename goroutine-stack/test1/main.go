package main

import (
	"fmt"
	"runtime"
	"time"
)

// burnStack 通过递归调用强制使用栈空间
func burnStack(depth int) {
	if depth == 0 {
		return
	}
	var buf [256]byte // 每层递归占用一些栈空间
	_ = buf
	burnStack(depth - 1)
}

func main() {
	runtime.GC() // 触发一次 GC
	time.Sleep(100 * time.Millisecond)

	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	beforeGoroutine := runtime.NumGoroutine()
	fmt.Printf("初始 StackInuse: %d KB\n", before.StackInuse/1024)
	fmt.Printf("初始 Goroutines: %d\n\n", beforeGoroutine)

	const N = 10000
	done := make(chan struct{})

	// 创建 N 个 goroutine，每个都会使用一定的栈空间
	for i := 0; i < N; i++ {
		go func() {
			burnStack(50) // 调用深度 50，强制使用栈
			<-done
		}()
	}

	time.Sleep(200 * time.Millisecond) // 等待所有 goroutine 启动并执行
	runtime.ReadMemStats(&after)

	stackIncrease := int64(after.StackInuse - before.StackInuse)
	goroutineIncrease := runtime.NumGoroutine() - beforeGoroutine

	fmt.Printf("=== 创建 goroutine 后 ===\n")
	fmt.Printf("当前 StackInuse: %d KB\n", after.StackInuse/1024)
	fmt.Printf("当前 Goroutines: %d\n", runtime.NumGoroutine())
	fmt.Printf("栈增长: %d KB\n", stackIncrease/1024)
	fmt.Printf("新增 Goroutines: %d\n\n", goroutineIncrease)

	if goroutineIncrease > 0 {
		avgStackPerGoroutine := stackIncrease / int64(goroutineIncrease)
		fmt.Printf("平均每个 goroutine 栈占用: %d KB\n", avgStackPerGoroutine/1024)
	}

	close(done)
	time.Sleep(100 * time.Millisecond)
}
