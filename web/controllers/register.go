package controllers

import (
	"github.com/gorilla/securecookie"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"github.com/kataras/iris/v12/sessions/sessiondb/redis"
	"shop/client"
	"shop/util"
	"time"
)

func RegisterControllers( app *iris.Application) {
	hashKey := []byte("the-big-and-secret-fash-key-here")
	blockKey := []byte("lot-secret-of-characters-big-too")
	secureCookie := securecookie.New(hashKey, blockKey)

	sessManager := sessions.New(sessions.Config{
		Cookie:  util.COOKEI_NAME,
		Expires: 24 * time.Hour,
		Encode:       secureCookie.Encode,
		Decode:       secureCookie.Decode,
		AllowReclaim: true,
	})
	db := redis.New(redis.Config{
		Network:   "tcp",
		//Addr:      "127.0.0.1:6379",
		Addr:      "192.168.0.223:6379",
		Timeout:   time.Duration(30) * time.Second,
		MaxActive: 10,
		Password:  "",
		Database:  "",
		Prefix:    "",
		Delim:     "-",
		Driver:    redis.Redigo(), // redis.Radix() can be used instead.
	})
	sessManager.UseDatabase(db)

	etcdApp := mvc.New(app.Party("/etcd"))
	etcdApp.Register(client.EtcdClientInsance)
	etcdApp.Handle(new(EtcdController))

	home := mvc.New(app.Party("/"))
	home.Register(
		sessManager.Start,
	)
	home.Handle(new(HomeController))

	self := mvc.New(app.Party("/self"))
	self.Register(
		sessManager.Start,
	)
	self.Handle(new(SelfController))

	shopCar := mvc.New(app.Party("/shopcar"))
	shopCar.Handle(new(ShopCarController))

	goodsDetail := mvc.New(app.Party("/goodsdetail"))
	goodsDetail.Handle(new(GoodsDetailController))

	assort := mvc.New(app.Party("/assort"))
	assort.Handle(new(AssortController))

	buy := mvc.New(app.Party("/goodsdetail/buy"))
	buy.Register(
		sessManager.Start,
	)
	buy.Handle(new(BuyController))

	order := mvc.New(app.Party("/order"))
	order.Register(
		sessManager.Start,
	)
	order.Handle(new(OrderController))

	user := mvc.New(app.Party("/user"))
	userService := client.NewUserService()
	user.Register(
		userService,
		sessManager.Start,
	)
	user.Handle(new(UserGController))
}
