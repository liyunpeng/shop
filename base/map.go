package main

import (
	"fmt"
	"runtime"
	"strconv"
	"time"
)

const (
	numElements = 100000
)

/*
	在垃圾回收过程中，运行时会扫描包含指针的对象并遍历其指针。
	map[string]int，那么垃圾回收器就不得不在每次垃圾回收过程中检查map中的每个字符串，因为字符串包含指针。
	这个例子中我们向一个map[string]int中写入了一千万个元素，
	然后测量垃圾回收的时间。map是在包的作用域中分配的，以保证它被分配到堆上。
*/
var Pointerfoo = map[string]int{}

func timeGC() {
	t := time.Now()
	runtime.GC()
	/*
		打印垃圾回收运行的时间
	*/
	fmt.Printf("gc took: %s\n", time.Since(t))
}

func pointerMap() {
	for i := 0; i < numElements; i++ {
		Pointerfoo[strconv.Itoa(i)] = i
	}

	for {
		timeGC()
		time.Sleep(1 * time.Second)
	}
	/*
	运行结果:
		gc took: 104.821989ms
		gc took: 98.726321ms
		gc took: 105.524633ms
		gc took: 102.829451ms
		gc took: 102.71908ms
		gc took: 103.084104ms
		gc took: 104.821989ms
	*/
}

/*
因为垃圾回要扫描包含指针的对象并遍历其指针
所以尽量去掉指针，这样能减少垃圾回收器需要遍历的指针数量。
由于字符串包含指针，因此我们可以用map[int]int来实现
*/
var valuefoo = map[int]int{}

func valueMap() {
	for i := 0; i < numElements; i++ {
		valuefoo[i] = i
	}

	for {
		timeGC()
		time.Sleep(1 * time.Second)
	}
	/*
		运行结果:
		gc took: 3.608993ms
		gc took: 3.926913ms
		gc took: 3.955706ms
		gc took: 4.063795ms
		gc took: 3.91519ms
		gc took: 3.75226ms
	*/
}

func main() {

	pointerMap()
	valueMap()
}
