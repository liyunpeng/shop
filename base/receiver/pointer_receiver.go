package main

import "fmt"

type introduction interface {
	intro()
}
type stut1 struct {
	name  string
	grade int
}
type cat struct {
	name string
	age  int
}

//值接收器
func (s stut1) intro() {
	fmt.Printf("%s's grade is %d\n", s.name, s.grade)
}

//指针接收器
func (c *cat) intro() {
	fmt.Printf("%s is %d years old\n", c.name, c.age)
}
func main() {
	stu := stut1{
		name:  "kyton",
		grade: 96,
	}
	ca := cat{
		name: "tom",
		age:  6,
	}
	var s1, c1 introduction
	fmt.Println("原始数据：", stu)
	s1 = stu
	s1.intro()
	s1 = &stu
	s1.intro()
	fmt.Println("原始传递：", ca)
	//c1 = ca
	/*Cannot use 'ca' (type cat) as type introduction Type does not
	implement 'introduction' as 'intro' method has a pointer receiver
	*/
	c1 = &ca
	c1.intro()
	(&ca).intro()
}
