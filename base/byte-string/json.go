package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

// 你可以使用 string 来存储数值数据，在 decode 时再决定按 int 还是 float 使用
// 将数据转为 decode 为 string
func f1() {
	var data = []byte(`{"status": 200}`)
	var result map[string]interface{}
	var decoder = json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&result); err != nil {
		log.Fatalln(err)
	}
	var status uint64
	err := json.Unmarshal([]byte(result["status"].(json.Number).String()), &status);
	checkError(err)
	fmt.Println("Status value: ", status)
}

func f2() {
	var data = []byte(`{"status": 200}`)
	var result struct {
		Status uint64 `json:"status"`
	}

	err := json.NewDecoder(bytes.NewReader(data)).Decode(&result)
	checkError(err)
	fmt.Printf("err:, %+v, Result: %+v", err, result)
}


func main() {
	var data = []byte(`{"status": 200}`)
	var result map[string]interface{}

	if err := json.Unmarshal(data, &result); err != nil {
		log.Fatalln(err)
	}

	/*
	在 encode/decode JSON 数据时，Go 默认会将数值当做 float64 处理
	 */
	fmt.Printf("%T\n", result["status"])    // float64
	var status = result["status"].(int)    // 类型断言错误
	fmt.Println("Status value: ", status)
}

func checkError(e  error) {

}