package designpattern

import "fmt"

// 这个工厂中所有产品要实现的方法， 这里产品名叫Shape
type Shape interface {
	Draw()
}

type Rectangle struct {
}
func (this Rectangle) Draw() {
	fmt.Println("Inside Rectangle::draw() method.")
}

type Square struct {
}
func (this Square) Draw() {
	fmt.Println("Inside Square ::draw() method.")
}

type Circle struct {
}
func (this Circle) Draw() {
	fmt.Println("Inside Circle  ::draw() method.")
}

type ShapeFactory struct {
}
func (this ShapeFactory) getShape(shapeType string) Shape {

	if shapeType == "" {
		return nil
	}
	if shapeType == "CIRCLE" {
		return Circle{}
	} else if shapeType == "RECTANGLE" {
		return Rectangle{}
	} else if shapeType == "SQUARE" {
		return Square{}
	}
	return nil
}