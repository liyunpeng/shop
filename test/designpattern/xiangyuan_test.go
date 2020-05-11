package designpattern

import (
	"fmt"
	"math/rand"
	"testing"
)

type Shapea interface {
	Draw3()
	SetX(x int)
	SetY(y int)
	SetRadius(radius int)
}

type Circlea struct {
	color  string
	x      int
	y      int
	radius int
}

func (c *Circlea) Circlea(color string) {
	c.color = color
}

func (c *Circlea) Draw3() {
	fmt.Println("Circlea: Draw() [Color :", c.color, "x:", c.x, "y:", c.y, "radius:", c.radius)
}
func (c *Circlea) SetX(x int) {
	c.x = x
}
func (c *Circlea) SetY(y int) {
	c.y = y
}
func (c *Circlea) SetRadius(radius int) {
	c.radius = radius
}


//创建一个工厂，生成基于给定消息的实体类的对象
type ShapeaFactory struct {
	CircleaMap map[string]Shapea
}

func (s *ShapeaFactory) GetCirclea(color string) Shapea {
	if s.CircleaMap == nil {
		s.CircleaMap = make(map[string]Shapea)
	}

	circlea := s.CircleaMap[color]
	if circlea == nil {
		newCirclea := new(Circlea)
		newCirclea.Circlea(color)
		s.CircleaMap[color] = newCirclea
		fmt.Println("Creating Circlea of color : ", color)
		circlea = newCirclea
	}
	return circlea
}


//测试：
type FlyweightPatternDemo struct {
	colors []string
}

func (f *FlyweightPatternDemo) GetRandomColor() string {
	if f.colors == nil {
		f.colors = []string{"Red", "Green", "Blue", "White", "Black"}
	}

	num := rand.Intn(len(f.colors) - 1)
	return f.colors[num]
}

func (f *FlyweightPatternDemo) GetRandomXAndY() int {
	num := rand.Intn(10) * 100
	return num
}
func TestFlyweightPattern(t *testing.T) {
	flyweightPatternDemo := FlyweightPatternDemo{}
	ShapeaFactory := new(ShapeaFactory)
	for i := 0; i < 20; i++ {
		Circlea := ShapeaFactory.GetCirclea(flyweightPatternDemo.GetRandomColor())
		Circlea.SetX(flyweightPatternDemo.GetRandomXAndY())
		Circlea.SetY(flyweightPatternDemo.GetRandomXAndY())
		Circlea.SetRadius(100)
		Circlea.Draw3()
	}
}