### 解决vue在接收websocket服务数据，vue的onmessage走不到的问题
websocket服务器返回的数据必须是结构化的数据， vue的onmessage才会走到， 否则死活走不到，
res的结构化数据， 必须包括code,msg,data三个部分：
``` 
res := &Response{
	Code: 1,
	Msg:  "success",
	Data: msg,
}
jsonHandler.Send(conn, res)
```
这样在vue的onmessage接受websocket服务数据的函数才会走到：
```
this.websock.onmessage = this.onmessage
onmessage: function (e) {
   ...
}
```

能断点调试，还是断点调试， 能最快定位问题，goland启动断点调试很快， 浏览器启动断电调试也很快。


  
