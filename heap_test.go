package heapq

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"testing"
)

func TestCapacity(t *testing.T) {
	hq := NewHeap[int](100_000, func(a, b *int) bool {
		return b == nil || a != nil && *a < *b
	})

	for idx := 0; idx < 100_000; idx++ {
		err := hq.Push(idx)
		if err != nil {
			t.Fatalf("%v", err)
		}
	}

	for idx := 0; idx < 100_000; idx++ {
		el, err := hq.Pop()
		if err != nil {
			t.Fatalf("%v", err)
		}
		if el != idx {
			t.Fatal("incorrect element")
		}
	}
}

func TestOrdering(t *testing.T) {
	hq := NewHeap[int](10, func(a, b *int) bool {
		return b == nil || a != nil && *a < *b
	})

	hq.Push(10)
	hq.Push(50)
	hq.Push(100)
	hq.Push(5)
	hq.Push(25)
	hq.Push(75)
	hq.Push(150)

	if e, _ := hq.Peek(); e != 5 {
		t.Fatal("incorrect smallest")
	}

	expected := []int{5, 10, 25, 50, 75}
	for idx := 0; idx < len(expected); idx++ {
		e, _ := hq.Pop()
		if e != expected[idx] {
			t.Fatal("incorrect smallest")
		}
	}

	hq.Push(50)
	hq.Push(10)
	hq.Push(1000)

	expected = []int{10, 50, 100, 150, 1000}
	for idx := 0; idx < len(expected); idx++ {
		e, _ := hq.Pop()
		if e != expected[idx] {
			t.Fatal("incorrect smallest")
		}
	}
}

func benchmarkTopNBaseline(topn int, maxsize int, b *testing.B) {
	var arr []int = make([]int, maxsize, maxsize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// random init
		for idx := 0; idx < len(arr); idx++ {
			arr[idx] = rand.Int()
		}
		b.StartTimer()
		sort.Ints(arr[:])
	}
}

func benchmarkTopN(topn int, maxsize int, b *testing.B) {
	var arr []int = make([]int, maxsize, maxsize)
	hq := NewHeap[int](topn, func(a, b *int) bool { return b == nil || a != nil && *a < *b })
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// random init
		for idx := 0; idx < len(arr); idx++ {
			arr[idx] = rand.Int()
		}
		b.StartTimer()
		for idx := 0; idx < len(arr); idx++ {
			if hq.size >= topn {
				hq.Pop()
			}
			hq.Push(arr[idx])
		}
		for idx := 0; idx < topn; idx++ {
			hq.Pop()
		}
	}
}

func BenchmarkFullSortBaseline(b *testing.B) {
	benchmarkTopNBaseline(100_000, 100_000, b)
}

func BenchmarkFullSortHeapq(b *testing.B) {
	benchmarkTopN(100_000, 100_000, b)
}

func BenchmarkTop100SortHeapq(b *testing.B) {
	benchmarkTopN(100, 100_000, b)
}

func TestTopRuntimes(t *testing.T) {
	file, err := os.Create("bench.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	fmt.Fprintln(file, "queue size,sort.Ints,heap topN,heap baseline")
	for topn := 1; topn < 1000; topn += 1 {
		fnBaseline := func(b *testing.B) {
			benchmarkTopNBaseline(topn, 1000, b)
		}
		fnHeap := func(b *testing.B) {
			benchmarkTopN(topn, 1000, b)
		}
		fnHeapBaseline := func(b *testing.B) {
			benchmarkTopN(1000, 1000, b)
		}
		rBaseline := testing.Benchmark(fnBaseline)
		rHeap := testing.Benchmark(fnHeap)
		rHeapBaseline := testing.Benchmark(fnHeapBaseline)

		fmt.Fprintln(file, topn, ",",
			int(rBaseline.T)/rBaseline.N, ",", int(rHeap.T)/rHeap.N, ",",
			int(rHeapBaseline.T)/rHeapBaseline.N)
	}
}
