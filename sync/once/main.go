package main

import "sync"

func merge(s []int, middle int) {
	left := append([]int(nil), s[:middle]...)  // 复制左边
	right := append([]int(nil), s[middle:]...) // 复制右边

	i, j := 0, 0
	k := 0

	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			s[k] = left[i]
			i++
		} else {
			s[k] = right[j]
			j++
		}
		k++
	}

	// 把剩余部分补上
	for i < len(left) {
		s[k] = left[i]
		i++
		k++
	}
	for j < len(right) {
		s[k] = right[j]
		j++
		k++
	}
}
func sequentialMergesort(s []int) {
	if len(s) <= 1 {
		return
	}

	middle := len(s) / 2
	sequentialMergesort(s[:middle]) // First half
	sequentialMergesort(s[middle:]) // Second half
	merge(s, middle)                // Merges the two halves
}

func parallelMergesortV1(s []int) {
	if len(s) <= 1 {
		return
	}

	middle := len(s) / 2

	var wg sync.WaitGroup
	wg.Add(2)

	go func() { // Spins up the first half of the work in a goroutine
		defer wg.Done()
		parallelMergesortV1(s[:middle])
	}()

	go func() { // Spins up the second half of the work in a goroutine
		defer wg.Done()
		parallelMergesortV1(s[middle:])
	}()

	wg.Wait()
	merge(s, middle) // Merges the halves
}
func main() {

}
