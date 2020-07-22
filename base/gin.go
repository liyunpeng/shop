package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"

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
		baseGroup.POST("/container", CaculateContainerArrangement)
		baseGroup.POST("/container1", CaculateContainerArrangement1)
	}

	r.Run()
}

/*
最小模型
{
  "p1num": 8,
  "p1hourtime": 2.166666667,
  "p2hourtime": 2.6,
  "p4hourtime": 2.6
}


{
  "p1num": 8,
   amd 总核数：  ->  amd剩余核数 ：
   P2 32核 时间:
   p4 32核 时间：
  "p1hourtime": 2.166666667,
  "p2hourtime": 2.6,
  "p4hourtime": 2.6
}

{
  "p1num": 8,
   amd 总核数：  ->  amd剩余核数 ：
   P2 32核 时间:
   p4 32核 时间：
  "p1hourtime": ,
  "p2hourtime": 1.89,   23核
  "p4hourtime": 7.5    7.5个/小时  21个核  6个benchy
}

{
	"p1hourtime": 2.34,
	"p1core" : 2,
	"p2hourtime": 2.15,
	"p2core" :  23,
	"p4hourtime": 0.8,
	"p4core" :  21
}


{
	"p1num": 8,
	"p1hourtime": 0.85,
	"p4hourtime": 7.5
}
	"p2hourtime": 1.89,

n1= 4 , n2= 3.9687848 , n4= 1.6890516 , count= 1.7094017


*/


type InputParam struct {
	P1Num      float32     `json:"p1num"`
	P1hourtime float32 `json:"p1hourtime"`
	P1core float32 `json:"p1core"`
	P2hourtime float32 `json:"p2hourtime"`
	P2core float32 `json:"p2core"`
	P4hourtime float32 `json:"p4hourtime"`
	P4core float32 `json:"p4core"`
}

const HOURS =  504 //   24*21  HOURS

/*
{
	"p1hourtime": 0.85,
	"p1core" : 2,
	"p2hourtime": 1.89,
	"p2core" :  23
	"p4hourtime": 7.5
	"p2core" :  21
}

 */
func CaculateContainerArrangement(c *gin.Context) {
	inputParam := new(InputParam)
	c.ShouldBindJSON(inputParam)

	p1Speed := 1/ inputParam.P1hourtime
	p2Speed := 1/ inputParam.P2hourtime
	p4Speed := 1/ inputParam.P4hourtime

	p1core := inputParam.P1core
	p2core := inputParam.P2core
	p4core := inputParam.P4core

	fmt.Println(" p1speed=", p1Speed, "p2speed=", p2Speed, ", p4speed=", p4Speed)
	fmt.Println(" p1core=", p1core, "p2core=", p2core, ", p4core=", p4core)
	var n1 float32
	var n2 float32
	var n4 float32
	var coresum int

	var finaln1 int
	var finaln2 int
	var finaln4 int
	var max int
	max = 0
	for n1 =2; n1 < 128/p1core; n1 = n1+2{

		n2 = (128 - p1core *n1 - p4core * ( n1 * p1Speed / p4Speed)) /p2core
		n4 = (128 - p1core *n1 - p2core * ( n1 * p1Speed / p2Speed)) / p4core
		fmt.Println("n1=", n1, ", n2=", n2, ", n4=", n4, ", p1count=", n1*p1Speed, ", p2count=", n2*p2Speed, ", p4count=", n4*p4Speed)
		if (max < int(n1 * p1Speed) )  && ( n2 > 1 ) && ( n4 > 1 ) {
			max = int( n1 * p1Speed)
			finaln1 = int(n1)
			finaln2 = int(n2)
			finaln4 = int(n4)
			coresum = int (p1core) * finaln1 + int (p2core) * finaln2 + int (p4core)  * finaln4
			fmt.Println("max=", max, ", n1=", n1, ", n2=", n2, ", n4=", n4, ", P1count=", n1*p1Speed, ", P2count=",  ", coresum=",coresum)
		}
	}


	mes := Message{
		P1Speed: p1Speed,
		P2Speed: p2Speed,
		P4Speed: p4Speed,

		P1Machines: finaln1,   //
		P2Machines: finaln2,
		P4Machines: finaln4,
		P1Output:   max,
		CoreSum: coresum,
		//OutputPermachine: outputPermachine,
	}

	JSON(c, "配比数据输出", mes)

}



func CaculateContainerArrangement1(c *gin.Context) {
	inputParam := new(InputParam)
	c.ShouldBindJSON(inputParam)

	p1Speed := 1/ inputParam.P1hourtime
	p2Speed := 1/ inputParam.P2hourtime
	p4Speed := 1/ inputParam.P4hourtime

	p1core := inputParam.P1core
	p2core := inputParam.P2core
	p4core := inputParam.P4core

	fmt.Println(" p1speed=", p1Speed, "p2speed=", p2Speed, ", p4speed=", p4Speed)
	fmt.Println(" p1core=", p1core, "p2core=", p2core, ", p4core=", p4core)
	var n1 float32
	var n2 float32
	var n4 float32
	var coresum float32
	var finalCoresum float32
	var finaln1 float32

	for n1 =2; n1 < 128/p1core; n1++{
		coresum = n1* p1core + n2 * p2core + n4 * p4core
		if coresum > 128 {
			break;
		}
		finalCoresum = coresum

		finaln1 = n1
		n2 = n1* p1Speed/ p2Speed

		n4 = n1* p1Speed/ p4Speed

		fmt.Println("n1=", n1, ", n2=", n2, ", n4=", n4, ", p1count=", n1*p1Speed, ", p2count=", n2*p2Speed, ", p4count=", n4*p4Speed, ", coresum=", coresum)

	}


	mes := Message1{
		P1Speed: p1Speed,
		P2Speed: p2Speed,
		P4Speed: p4Speed,

		P1num: finaln1,
		P2num: n2,
		P4num: n4,
		CoreSum: int(finalCoresum),
		P1SpeedCount: finaln1*p1Speed,
		P2SpeedCount: n2*p2Speed,
		P4SpeedCount: n4*p4Speed,
		//OutputPermachine: outputPermachine,
	}

	JSON(c, "配比数据输出", mes)

}



func CaculateArrangement(c *gin.Context) {
	inputParam := new(InputParam)
	c.ShouldBindJSON(inputParam)

	p1Num := inputParam.P1Num
	//p1Speed := 1/ inputParam.P1hourtime
	//p2Speed := 1/ inputParam.P2hourtime
	//
	//p4Speed := 1/ inputParam.P4hourtime


	p1Speed := inputParam.P1hourtime
	p2Speed := inputParam.P2hourtime

	p4Speed := inputParam.P4hourtime
	//p4Core := inputParam.P4core
	//p1Core := inputParam.P1core
	//p2Core := inputParam.P2core
	p1Output := decimal.NewFromFloat(float64(p1Num)).Mul(decimal.NewFromFloat(float64(HOURS))).Mul(decimal.NewFromFloat(float64(p1Speed))).BigInt().Int64()
	p2Machines  := decimal.NewFromFloat(float64(p1Output)).Div(decimal.NewFromFloat(HOURS).Mul(decimal.NewFromFloat(float64(p2Speed)))).BigInt().Int64()
	p4Machines  := decimal.NewFromFloat(float64(p1Output)).Div(decimal.NewFromFloat(HOURS).Mul(decimal.NewFromFloat(float64(p4Speed)))).BigInt().Int64()
	outputPermachine := decimal.NewFromFloat(float64(p1Output)).Div(decimal.NewFromFloat(float64(p2Machines)).Add(decimal.NewFromFloat(float64(p4Machines))).Add(decimal.NewFromFloat(1))).BigInt().Int64()

	//cores :=  decimal.NewFromFloat(float64(p1Num)).Mul(decimal.NewFromFloat(float64(HOURS))).Mul(decimal.NewFromFloat(float64(p1Speed))).BigInt().Int64()

	mes := Message{
		//P1Speed: p1Speed,
		//P2Speed: p2Speed,
		//P4Speed: p4Speed,

		P1Machines: 1,   //
		P2Machines:  int(p2Machines),
		P4Machines: int(p4Machines),
		P1Output:   int(p1Output),
		OutputPermachine: outputPermachine,
	}

	JSON(c, "配比数据输出", mes)

}

type Message struct {

	P1Speed float32 `json:"p1speed"`
	P2Speed float32 `json:"p2speed"`
	P4Speed float32 `json:"p4speed"`


	P1Machines int `json:"p1machines"`
	P2Machines int `json:"p2machines"`
	P4Machines int `json:"p4machines"`
	P1Output   int `json:"p1output"`
	CoreSum int `json:"coresum"`
	OutputPermachine int64 `json:"outputPermachine"`
}


type Message1 struct {

	P1Speed float32 `json:"p1speed"`
	P2Speed float32 `json:"p2speed"`
	P4Speed float32 `json:"p4speed"`

	P1num float32  `json:"p1num"`
	P2num float32 `json:"p2num"`
	P4num float32 `json:"p4num"`

	CoreSum int `json:"coresum"`
	P1SpeedCount  float32 `json:"p1speedcount"`
	P2SpeedCount  float32 `json:"p2speedcount"`
	P4SpeedCount  float32 `json:"p4speedcount"`
}