package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	//"strings"
	"github.com/shopspring/decimal"
)

type CommonRes struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
	Error interface{} `json:"error"`
}

func JSON(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, CommonRes{
		Code: 0,
		Msg:  msg,
		Data: data,
	})
}

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "test")
	})

	baseGroup := r.Group("/")
	{
		baseGroup.POST("/arr", CaculateArrangement)
	}

	r.Run() // listen and serve on 0.0.0.0:8080
}

/*
{
  "p1num": 8,
  "p1hourtime": 2.166666667,
  "p2hourtime": 2.6,
  "p4hourtime": 2.6
}
*/
type InputParam struct {
	P1Num      int     `json:"p1num"`
	P1hourtime float32 `json:"p1hourtime"`
	P2hourtime float32 `json:"p2hourtime"`
	P4hourtime float32 `json:"p4hourtime"`
}

const HOURS =  504 //   24*21  HOURS

func CaculateArrangement(c *gin.Context) {
	inputParam := new(InputParam)
	c.ShouldBindJSON(inputParam)

	p1Num := inputParam.P1Num
	p1Speed := 1/ inputParam.P1hourtime
	p2Speed := 1/ inputParam.P2hourtime
	p4Speed := 1/ inputParam.P4hourtime

	p1Output := decimal.NewFromFloat(float64(p1Num)).Mul(decimal.NewFromFloat(float64(HOURS))).Mul(decimal.NewFromFloat(float64(p1Speed))).BigInt().Int64()
	p2Machines  := decimal.NewFromFloat(float64(p1Output)).Div(decimal.NewFromFloat(HOURS).Mul(decimal.NewFromFloat(float64(p2Speed)))).BigInt().Int64()
	p4Machines  := decimal.NewFromFloat(float64(p1Output)).Div(decimal.NewFromFloat(HOURS).Mul(decimal.NewFromFloat(float64(p4Speed)))).BigInt().Int64()
	outputPermachine := decimal.NewFromFloat(float64(p1Output)).Div(decimal.NewFromFloat(float64(p2Machines)).Mul(decimal.NewFromFloat(float64(p4Machines)))).BigInt().Int64()

	mes := Message{
		//P1Speed: p1Speed,
		//P2Speed: p2Speed,
		//P4Speed: p4Speed,

		P1Machines: 1,   //
		P2Machines: p2Machines,
		P4Machines: p4Machines,
		P1Output:   p1Output,
		OutputPermachine: outputPermachine,
	}

	JSON(c, "配比数据输出", mes)

}

type Message struct {

	//P1Speed float32 `json:"p1speed"`
	//P2Speed float32 `json:"p2speed"`
	//P4Speed float32 `json:"p4speed"`


	P1Machines int64 `json:"p1machines"`
	P2Machines int64 `json:"p2machines"`
	P4Machines int64 `json:"p4machines"`
	P1Output   int64 `json:"p1output"`
	OutputPermachine int64 `json:"outputPermachine"`
}
