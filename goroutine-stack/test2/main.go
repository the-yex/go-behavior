package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"
)

func deepCall(n int) {
	if n == 0 {
		return
	}
	var buf [512]byte // 每层占用更多栈空间
	_ = buf
	deepCall(n - 1)
}

func measureStackUsage(label string, count int, useDeep bool) {
	runtime.GC()
	time.Sleep(100 * time.Millisecond)

	var before, after runtime.MemStats
	beforeGoroutine := runtime.NumGoroutine()
	runtime.ReadMemStats(&before)

	done := make(chan struct{})

	for i := 0; i < count; i++ {
		go func() {
			if useDeep {
				deepCall(100)
			}
			<-done
		}()
	}

	time.Sleep(200 * time.Millisecond)
	runtime.ReadMemStats(&after)

	stackInc := int64(after.StackInuse - before.StackInuse)
	grInc := runtime.NumGoroutine() - beforeGoroutine
	avg := int64(0)
	if grInc > 0 {
		avg = stackInc / int64(grInc)
	}

	fmt.Printf("\n【%s】\n", label)
	fmt.Printf("新增 goroutines: %d\n", grInc)
	fmt.Printf("栈增长: %d KB\n", stackInc/1024)
	fmt.Printf("📊 平均栈占用: %d KB\n", avg/1024)

	close(done)
	time.Sleep(100 * time.Millisecond)
}

func main() {
	debug.SetGCPercent(-1) // 手动控制 GC

	fmt.Println("\n=== 阶段0：测量基线 ===")
	fmt.Println("创建空 goroutine，观察当前起始栈大小")
	measureStackUsage("第0批：空 goroutine", 5000, false)

	fmt.Println("\n=== 阶段1：训练 runtime ===")
	fmt.Println("创建大量深栈 goroutine，让 runtime 学习栈使用模式")
	measureStackUsage("第1批：深栈 goroutine", 5000, true)

	fmt.Println("\n>>> 触发 GC，runtime 将分析栈使用情况...")
	runtime.GC()
	time.Sleep(200 * time.Millisecond)

	fmt.Println("\n=== 阶段2：观察适应 ===")
	fmt.Println("再次创建空 goroutine，观察起始栈是否已经调整")
	measureStackUsage("第2批：空 goroutine", 5000, false)

	fmt.Println("\n=== 实验结束 ===")
	fmt.Println("💡 观察：对比第0批和第2批")
	fmt.Println("   - 如果第2批更大 → runtime 学会了使用更大的起始栈")
	fmt.Println("   - 如果相同 → 当前 Go 版本的自适应策略较保守")
}
