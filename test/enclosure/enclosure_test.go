package enclosure

import (
	"fmt"
	"testing"
)

func adder() func(int) int {
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}

//利用闭包实现斐波拉契数列
func fibonacci() func() int {
	b0 := 0
	b1 := 1
	return func() int {
		tmp := b0 + b1
		b0 = b1
		b1 = tmp
		return b1
	}

}

func TestEnclosure(t *testing.T) {
	myAdder := adder()

	// 从1加到3
	for i := 1; i <= 3; i++ {
		myAdder(i)
	}

	fmt.Println("myAdder(0)=", myAdder(0))
	fmt.Println("myAdder(10)", myAdder(10))
	//	输出：
	//	myAdder(0)= 6
	//	myAdder(10) 16

	myFibonacci := fibonacci()
	for i := 1; i <= 5; i++ {
		fmt.Println(myFibonacci())
	}
	// 输出：
	//	1
	//	2
	//	3
	//	5
	//  8

}
