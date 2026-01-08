package compare_concurrency

import (
	"testing"
)

/*
* @Author: zouyx
* @Email: 1003941268@qq.com
* @Date:   2026 2026/1/8 上午10:28
* @Package:
 */

func makeRandomSlice(n int) []int {
	s := make([]int, n)
	for i := 0; i < n; i++ {
		s[i] = n - i
	}
	return s
}
func getData() []int {
	return makeRandomSlice(10000)
}

func Benchmark_normal(b *testing.B) {
	data := getData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sequentialMergesort(data)
	}
}

func Benchmark_concurrency_V1(b *testing.B) {
	data := getData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parallelMergesortV1(data)
	}
}

func Benchmark_concurrency_V2(b *testing.B) {
	data := getData()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parallelMergesortV2(data)
	}
}
