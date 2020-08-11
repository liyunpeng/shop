package main

import (
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type RandName struct {
	first   string // first name dict 姓氏字典
	gender  string // gender 性别
	namenum int    // name character count 名字个数
	bdict   string // boy's name dict 男名字典
	gdict   string // girl's name dict 女名字典
	mdict   string // 中性用名
}

func (n RandName) writeFile(fname, content string) (dic []string) {
	fileInfo, err := os.Stat(fname)
	if err != nil {
		if os.IsExist(err) {
			log.Println(fileInfo)
		}
	}

	 err = ioutil.WriteFile(fname, []byte(content), os.ModeAppend )
	if err != nil {
		log.Fatal(err)
	}

	//dic = strings.Split(string(data), ",")
	//
	return
}

func (n RandName) readFile(fname string) (dic []string) {
	fileInfo, err := os.Stat(fname)
	if err != nil {
		if os.IsExist(err) {
			log.Println(fileInfo)
		}
	}

	data, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}

	dic = strings.Split(string(data), ",")

	return
}

var size = int64(1024 * 1024)

func ExampleTruncate() {
	f, err := os.Create("foobar1.bin")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := f.Truncate(size); err != nil {
		log.Fatal(err)
	}
}

func writeFile1( filename , wireteString string){
	var f *os.File
	var err1 error

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		f, err1 = os.Create(filename)
	}else{
		f, err1 = os.OpenFile(filename, os.O_APPEND, 0666) //打开文件
	}

	_, err1 = io.WriteString(f, wireteString)

	if err1 != nil {
		fmt.Println("error", err1)
		return
	}
	f.Sync()
	f.Close()
}



func (n *RandName) getFrist() (first string) {
	if n.first == "" {
		fdic := n.readFile("firstName.dat")

		rand := rand.New(rand.NewSource(time.Now().UnixNano()))

		res := rand.Intn(len(fdic))
		first = fdic[res]
	} else {
		first = n.first
	}

	return
}

func (n *RandName) getName() (boy string) {

	var fname string
	if n.gender == "女" {
		fname = "gname.dat"
	} else {
		fname = "bname.dat"
	}

	dic := n.readFile(fname)
	dicLen := len(dic)

	if dicLen <= 0 {
		log.Fatal("名字库内容为空")
	}

	if n.namenum == 0 {
		nseeder := rand.New(rand.NewSource(time.Now().UnixNano()))
		n.namenum = nseeder.Intn(100)%2 + 1
	}

	for index := 0; index < int(n.namenum); index++ {
		rand := rand.New(rand.NewSource(time.Now().UnixNano()))
		res := rand.Intn(dicLen)
		boy = boy + dic[res]
	}

	return
}

/**
 * 主入口函数
 */
func main() {

	u := RandName{first: "张", gender: "女"}
	for index := 0; index < 10; index++ {
		randName := u.getFrist() + u.getName()
		log.Println(randName)
	}

}
