package main

import "fmt"

type stut struct {
	name string
	grade int
}
//值接收器
func (s stut) changeName(nName string)  {
	s.name = nName
}
//指针接收器
func (s *stut) changeGrade(nGrade int)  {
	s.grade = nGrade
}
func main() {
	stu := stut{
		name: "kyton",
		grade: 96,
	}
	fmt.Println("原始数据：", stu)
	stu.changeName("newton")
	stu.changeGrade(98)
	fmt.Println("变量传递：", stu)
	(&stu).changeName("Einstein")
	(&stu).changeGrade(100)
	fmt.Println("地址传递", stu)
}
