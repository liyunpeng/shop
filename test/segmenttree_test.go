package test

import (
	"fmt"
	"testing"
)

//线段树节点
type Node struct {
	L     int
	R     int
	Left  *Node
	Right *Node
	Cover int
}

//建树
func (this *Node) build(l, r int) *Node {
	this.L = l
	this.R = r
	this.Cover = 0
	mid := (this.L + this.R) >> 1

	if (l + 1) < r {
		if this.L != mid {
			NewLeftNode := Node{}
			this.Left = NewLeftNode.build(this.L, mid)
		}
		if this.R != mid {
			NewRightNode := Node{}
			this.Right = NewRightNode.build(mid, this.R)
		}
		return this
	} else {
		return this
	}
}

//插入
func (this *Node) Insert(l, r int) {
	mid := (this.L + this.R) >> 1
	if l == this.L && r == this.R {
		//若插入的线段正好被包含，结束递归
		this.Cover = 1
	} else if r <= mid {
		if this.Left != nil {
			this.Left.Insert(l, r)
		}
	} else if l >= mid {
		if this.Right != nil {
			this.Right.Insert(l, r)
		}
	} else {
		if this.Left != nil {
			this.Left.Insert(l, mid)
		}
		if this.Right != nil {
			this.Right.Insert(mid, r)
		}
	}
}

//统计
func (this *Node) Sum() int {
	if this.Cover == 1 {
		return this.R - this.L
	} else if this.R-this.L == 1 {
		return 0
	}
	return this.Left.Sum() + this.Right.Sum()
}

//中序遍历
func (this *Node) InorderTraversal() {
	if this != nil {
		this.Left.InorderTraversal()
		fmt.Println(this)
		this.Right.InorderTraversal()
	}
}

func Inserts() {
	Tree := Node{}
	Tree.build(1, 8)
	Tree.Insert(1, 2)
	Tree.Insert(3, 5)
	Tree.Insert(4, 6)
	Tree.Insert(5, 6)
	Tree.Insert(7, 8)
	Tree.InorderTraversal()
	fmt.Println(Tree.Sum())
}

func TestInserts( t *testing.T){
	Inserts()
}