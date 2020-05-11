package designpattern

import (
	"fmt"
	"testing"
)

/*
模板设计模式定义：
模板方法模式定义了一个算法的步骤，
并允许子类别为一个或多个步骤提供其实践方式。
让子类别在不改变算法架构的情况下，重新定义算法中的某些步骤.

模板设计模式用处：
系统的组件都是按照一定的流程执行, 并且不同的组件的实现方式不同,需要我们延迟到子类来实现.
模板方法模式的关键点还是在将具体操作延迟到子类实现.
*/

type IPeople interface {
	SetName(name string)
	BeforeOut()
	Out()
}

type People struct {
	Specific IPeople
	name     string
}

func (this *People) SetName(name string) {
	this.name = name
}

func (this *People) Out() {
	this.BeforeOut()
	fmt.Println(this.name + " go out...")
}

func (this *People) BeforeOut() {
	if this.Specific == nil {
		return
	}

	/*
		 按道理将这个BeforeOut是要延迟到子类去实现的,
		这个不符合国际惯例啊? 我们还是来看看这个方法的内容吧,
		其实关键点就一句话, this.Specific.BeforeOut(),
		这句话直接调用了我们上面说的那个绑定的实例的BeforeOut方法,
		仔细品味一下, 还是符合惯例的,
		这里我们只不过用了一种折中的方式来实现方法的延迟实现.
	*/
	this.Specific.BeforeOut()
}

/*
用绑定的方式来实现模板方法模式,
这里的这个特殊的字段其实就是我们要绑定的实例,
在我们这个例子中就是男生或者女生了,
*/
type Boy struct {
	People
}

func (_ *Boy) BeforeOut() {
	fmt.Println("get up..")
}

type Girl struct {
	People
}

func (_ *Girl) BeforeOut() {
	fmt.Println("dress up..")
}

func TestMoBan(t *testing.T) {
	var p *People = &People{}

	p.Specific = &Boy{}
	p.SetName("qibin")
	p.Out()

	p.Specific = &Girl{}
	p.SetName("loader")
	p.Out()
}
