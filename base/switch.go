package main

import "fmt"

func main() {
	isSpace := func(ch byte) bool {
		switch (ch) {
		case ' ': //error
		case '\t':
			return true
		}
		return false
	}
	fmt.Println(isSpace('\t')) //prints true (ok)
	fmt.Println(isSpace(' '))  //prints false (not ok)  原因是不能顺序执行

	isSpace1 := func(ch byte) bool {
		switch (ch) {
		case ' ', '\t': // 要想顺序执行， 就得放在一起
			return true
		}
		return false
	}
	fmt.Println(isSpace1('\t')) //prints true (ok)
	fmt.Println(isSpace1(' '))  //prints true (ok)
}
