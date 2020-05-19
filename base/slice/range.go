package main

import "fmt"

//如果要修改原有元素的值，应该使用索引直接访问：
func main() {
	data := []int{1, 2, 3}
	for i, v := range data {
		data[i] = v * 10
	}
	fmt.Println("data: ", data)    // data:  [10 20 30]
}
