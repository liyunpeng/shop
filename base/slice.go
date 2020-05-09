package main

import "fmt"

func sliceRange() {
	data := []int{1, 2, 3}
	for _, v := range data {
		v *= 10 //original item is not changed
	}
	fmt.Println("data:", data) //prints data: [1 2 3]
}
func sliceRange1() {
	data := []int{1, 2, 3}
	for i, _ := range data {
		data[i] *= 10
	}
	fmt.Println("data:", data) //prints data: [10 20 30]
}

func get() []byte {
	raw := make([]byte, 10000)
	fmt.Println(len(raw), cap(raw), &raw[0]) //prints: 10000 10000 <byte_addr_x>
	return raw[:3]                           // 虽然返回前三个数据， 但cap计算还是1000
}
func cap1() {
	data := get()
	fmt.Println(len(data), cap(data), &data[0]) //prints: 3 10000 <byte_addr_x>
}

func sliceAppendNew() {

	s1 := []int{1, 2, 3}
	fmt.Println(len(s1), cap(s1), s1) //prints 3 3 [1 2 3]
	s2 := s1[1:]
	fmt.Println(len(s2), cap(s2), s2) //prints 2 2 [2 3]
	for i := range s2 {
		s2[i] += 20
	}
	//still referencing the same array
	fmt.Println(s1) //prints [1 22 23]
	fmt.Println(s2) //prints [22 23]
	s2 = append(s2, 4)
	for i := range s2 {
		s2[i] += 10
	}
	//s1 is now "stale"
	fmt.Println(s1) //prints [1 22 23]
	fmt.Println(s2) //prints [32 33 14]
}


func Slice() {
	fmt.Println("<----------------------- Slice begin ---------------------->")
	sliceRange()
	sliceRange1()
	cap1()
	sliceAppendNew()
	fmt.Println("<----------------------- Slice end ---------------------->")
}


