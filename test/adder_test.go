package test

import (
	"fmt"
	"testing"
)

// Add takes two integers and returns the sum of them
func Add(x, y int) int {
	//slice := make([]int, 100, 200)
	//fmt.Println("slice :", slice)
	return x + y
}

/*
	性能测试的函数名必须以Benchmark为开头，错一个字母都不会执行

 */
func BenchmarkAdd(t *testing.B){
	fmt.Println("B.N=", t.N)

	for i :=0;  i< t.N; i++ {
		_ = Add(1, 2)
	}
}

func heap() ([]byte ) {
	return make([]byte, 1024*10)
}

func BenchmarkHeap(t *testing.B){
	fmt.Println("B.N=", t.N)

	for i :=0;  i< t.N; i++ {
		_ = heap()
	}
}

func ExampleAdd() {
	sum := Add(1, 5)
	fmt.Println(sum)
	// Output: 6
}

func TestAdd(t *testing.T){
	if Add(1, 5)!= 6 {
		t.Fatal("ExampleAdd fatal")
	}
}
