package designpattern

import (
	"fmt"
	"shop/logger"
	"testing"
)

type Values map[string][]string
type Request struct {
	Method   string
	Url      string
	PostForm Values
}

func MakeRequest() *Request {
	r := Request{}
	r.PostForm = make(Values)
	return &r
}

func SendRequst(request *Request,mux *Mux){
	mux.root.Process(request)
}

func (this *Request) SetValues(key string, values ...string) {
	this.PostForm[key] = append(this.PostForm[key], values...)
}

func (this *Request) Post(url string,mux *Mux) {
	this.Method = "POST"
	this.Url=url
	logger.Info.Println("请求中的method这里被当作责任链上责任者的名字来用， 这个责任者的名字=",
		this.Method)
	SendRequst(this,mux)
}

// 定义一个IProcess接口，责任链的每一个处理处理者都应该实现这个接口。
type IProcess interface {
	Process(request *Request)
}

// 创建一个Processer结构类型实现这个接口，
// 这个结构类型也作为父类， 被RootProcesser和PostProcesser结构继承。
type HandlersCollection map[string]IProcess
type Processer struct {
	Handlers HandlersCollection
}

func (this *Processer) SetHandler(key string, process IProcess) {
	logger.Info.Println("SetHandler用于本处理者设置责任的下一个处理者，这里下一个处理者的名字=", key)

	this.Handlers[key] = process
	logger.Info.Println("下个处理者可以有很多个，每个责任处理者都有自己的handlers集合, 表示下一个处理者的集合， handlers=",
		this.Handlers)

}

func (this *Processer) Init() {
	this.Handlers = make(HandlersCollection)
}

//RootProcesser和PostProcesser都继承基础处理者。
//  RootProcesser作为根节点，
// 它的职责是把收到的请求，根据请求的方法传递给相应方法的处理节点，这里只有Post方法。
// 因为setHandler只设置了一个键值对
type RootProcesser struct {
	Processer
}

func (this *RootProcesser) Process(request *Request) {
	logger.Info.Println(" 每个责任者都要实现process, 实现方式有两种，一种是转交给下个处理者， 一种是自己真正完成这个请求处理")
	logger.Info.Println("RootProcesser 处理者处理一个请求， 该请求为", request)
	logger.Info.Println("RootProcesser 处理者处理方式，是把请求转交给责任链上的下一个处理者")
	this.Handlers[request.Method].Process(request)
	logger.Info.Println("责任链的处理者还是把处理结果返回给上级处理者， 这里")
}

func newRootProcesser() *RootProcesser {
	root := new(RootProcesser)
	root.Processer.Init()
	return root
}

type PostProcesser struct {
	Processer
	/*
	PostProcesser收到请求后，根据Url和多路复用器，将请求传递给相应Url的处理函数。
	处理函数在Mux结构中注册，以便被调用。
	责任处理者需要调用的函数都封装在mux里面。 责任者类可以通过成员引用到这个mux.
	*/
	PMux *Mux
}

func newPostProcesser(mux *Mux) *PostProcesser {
	post := new(PostProcesser)
	post.Processer.Init()
	post.PMux = mux
	return post
}

func (this *PostProcesser) Process(request *Request) {
	logger.Info.Println("PostProcesser处理者处理一个请求，该请求为", request)
	logger.Info.Println("PostProcesser处理者没有把请求转给下一个处理者，而是实实在在的为请求做事情")

	this.PMux.mux[request.Url](request)
}

//一个map来保存处理器的引用。
//模拟一个Request的发送。
type HandlerFunc func(request *Request)
type muxEntry map[string]HandlerFunc

// mux表示一个责任链
type Mux struct {
	mux  muxEntry
	root IProcess
}

// 增加一个函数名键值对
func (this *Mux) Handle(url string, handlerFunc HandlerFunc) {
	logger.Info.Println("一个责任链可以接受很多请求，一个请求用一个键值对表示，key为请求的url=", url,
		"value=", handlerFunc)
	this.mux[url] = handlerFunc
}

func (this *Mux) SetRootProcess(root IProcess) {
	this.root = root;
}

// 设置好责任链上的每个处理者
func NewMux() *Mux {
	mux := Mux{}
	mux.mux = make(muxEntry)

	// 本程序设计， 责任链中只有两个责任处理者
	// 责任的根处理者， 即责任链的第一个处理者
	root := newRootProcesser()

	// 载新建一个责任的处理者
	post := newPostProcesser(&mux)
	// 设置责任链额第二个处理者
	root.SetHandler("POST", post)

	// 设置好责任链的第一个处理者, 即责任链的根处理者
	mux.root = root;
	return &mux
}

//Login 函数为 HandlerFunc类型， 由此可以被map引用到
func Login(request *Request){
	username := request.PostForm["username"][0]
	password := request.PostForm["password"][0]

	fmt.Println(username)
	fmt.Println(password)
}

func TestChainResponsibity(t *testing.T){
	logger.InitCustLogger()


	mux := NewMux()
	mux.Handle("login",Login)


	req := MakeRequest()
	req.Method="POST"
	req.SetValues("username","111")
	req.SetValues("password","222")
	req.Post("login",mux)
}