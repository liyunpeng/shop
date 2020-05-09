package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func compareString() {

	var str string = "one"
	var in interface{} = "one"
	fmt.Println("str == in:", str == in, reflect.DeepEqual(str, in))
	//prints: str == in: true true
	v1 := []string{"one", "two"}
	v2 := []interface{}{"one", "two"}
	fmt.Println("v1 == v2:", reflect.DeepEqual(v1, v2))
	//prints: v1 == v2: false (not ok)
	data := map[string]interface{}{
		"code":  200,
		"value": []string{"one", "two"},
	}
	encoded, _ := json.Marshal(data)
	var decoded map[string]interface{}
	json.Unmarshal(encoded, &decoded)
	fmt.Println("data == decoded:", reflect.DeepEqual(data, decoded))
	//prints: data == decoded: false (not ok)
}

func byteSlice(b []byte) []byte {
	return b
}

func runeSlice(r []rune) []rune {
	return r
}

/*
Go语言中byte和rune实质上就是uint8和int32类型:
type byte = uint8
type rune = int32
从定义看出:
byte是uint8别名,
rune是int32的别名
*/
func byteRune() {
	b := []byte{0, 1}
	u8 := []uint8{2, 3}
	fmt.Printf("%T %T \n", b, u8)
	fmt.Println(byteSlice(b))
	fmt.Println(byteSlice(u8))
	/*
	   运行结果:
	      []uint8 []uint8
	      [0 1]
	      [2 3]
	*/

	r := []rune{4, 5}
	i32 := []int32{6, 7}
	fmt.Printf("%T %T \n", r, i32)
	fmt.Println(runeSlice(r))
	fmt.Println(runeSlice(i32))
	/*
		[]int32 []int32
		[4 5]
		[6 7]
	*/
}

func changeChar() {
	/*
		go中string是常量，只能用双引号来表示。
		a := "this is string"
		a[0] = 'c' (这个是错误的，会报错)
		如果要做改变某个字符的操作应该先将string转换为byte数组：
	*/
	a := "this is string"
	c := []byte(a)
	c[0] = 'c'
	d := string(c)
	fmt.Println(d)
}

func rune2() {
	/*
		字符串的汉字占三个字节， 字符串末尾没有结束符号， 所以输出len(s)=8

	*/
	s := "go有相"

	/*
		字符串类型转换数组， 不管是字节数组，还是rune数组， 都是做了一次复制，
		对数组修改不能改变原字符串
	*/
	b := []byte(s)
	b[0] = 'b'
	fmt.Println(s)
	/*
		运行结果：
		go有相
	*/
	r := []rune(s)
	r[0] = 'r'
	fmt.Println(s)
	/*
		运行结果：
		go有相
	*/

	/*
		字节数组长度和原字符串长度相同
		rune的则是实际的字符个数, 汉字也是一个字符,
		虽然一个汉字占3个字节, 但rune还是用4个字节存它
	*/
	fmt.Printf("s address=%p, len(s)=%d\n", s, len(s))
	fmt.Printf("byte address=%p, len([]byte(s))=%d \n", b, len(b))
	fmt.Printf("rune address=%p, len([]rune(s))=%d \n", r, len(r))

	/*
		运行结果
		s address=%!p(string=go有相), len(s)=8
		byte address=0xc000385200, len([]byte(s))=8
		rune address=0xc000385210, len([]rune(s))=4
	*/

	for i := 0; i < len(b); i++ {
		fmt.Println("byte[", i, "]=", b[i])
		fmt.Printf("byte[%d]=%c \n", i, b[i])
	}
	/*
			运行结果：
		byte[ 0 ]= 103
		byte[0]=g
		byte[ 1 ]= 111
		byte[1]=o
		byte[ 2 ]= 230
		byte[2]=æ
		byte[ 3 ]= 156
		byte[3]=
		byte[ 4 ]= 137
		byte[4]=
		byte[ 5 ]= 231
		byte[5]=ç
		byte[ 6 ]= 155
		byte[6]=
		byte[ 7 ]= 184
		byte[7]=¸
	*/
	for i := 0; i < len(r); i++ {
		fmt.Println("rune[", i, "]=", r[i])
		fmt.Printf("rune[%d]=%c \n", i, r[i])
	}
	/*
			运行结果
		rune[ 0 ]= 103
		rune[0]=g
		rune[ 1 ]= 111
		rune[1]=o
		rune[ 2 ]= 26377
		rune[2]=有
		rune[ 3 ]= 30456
		rune[3]=相
	*/
}

func String() {
	compareString()

	changeChar()

	byteRune()

	rune2()
}
