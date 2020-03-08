package models

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	"shop/config"
	_ "shop/config"
	"shop/validates"
	"time"
)
var Enforcer *casbin.Enforcer
var err error
var c *gormadapter.Adapter
var DB *gorm.DB
func Register(conf *config.Config) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容
		}
	}()
	mysqlConf := conf.Mysql

	DBConn := "root:123456@/gotest?charset=utf8&parseTime=True&loc=Local"
	DB, err = gorm.Open(
		"mysql",  DBConn)
	//"mysql", "root:password@tcp(192.168.0.220:31111)/gotest?charset=utf8&parseTime=True&loc=Local")
	if err == nil {
		DB.DB().SetMaxIdleConns(mysqlConf.MaxIdle)
		DB.DB().SetMaxOpenConns(mysqlConf.MaxOpen)
		DB.DB().SetConnMaxLifetime(time.Duration(300) * time.Minute)
		err = DB.DB().Ping()
		fmt.Println("成功连接数据库")
	} else {
		fmt.Println("没有连接到数据库 err= ", err)
		panic("数据库错误")
	}
	driverName := "mysql"
	casbinConn  :="root:123456@(127.0.0.1:3306)/"
	c, err = gormadapter.NewAdapter(driverName, casbinConn) // Your driver and data source.
	if err != nil {
		color.Red(fmt.Sprintf("NewAdapter 错误: %v", err))
	}

	Enforcer, err =  casbin.NewEnforcer("./config/rbac_model.conf", c)
	if err != nil {
		color.Red(fmt.Sprintf("NewEnforcer 错误: %v", err))
	}
	_ = Enforcer.LoadPolicy()

	//a := gormadapter.NewAdapter("mysql", DBConn)
	//e := casbin.NewEnforcer("rabc.conf", a)
	////从DB加载策略
	//e.LoadPolicy()
	//
	////获取router路由对象
	//r := gin.New()
	////使用自定义拦截器中间件
	//r.Use(LanjieqiHandler(e))
}


//拦截器
//func LanjieqiHandler(e *casbin.Enforcer) gin.HandlerFunc {
//	return func(c *iris.Context) {
//		//获取请求的URI
//		obj := c.Request.URL.RequestURI()
//		//获取请求方法
//		act := c.Request.Method
//		c.
//		//获取用户的角色
//		sub := "admin"
//
//		//判断策略中是否存在
//		if e.Enforce(sub, obj, act) {
//			fmt.Println("通过权限")
//			c.Next()
//		} else {
//			fmt.Println("权限没有通过")
//			c.Abort()
//		}
//	}
//}Ò



func CreateSystemData(perms []*validates.PermissionRequest) {
	//if rc.App.CreateSysData == 1 {
	if true {
		permIds := CreateSystemAdminPermission(perms) //初始化权限
		role := CreateSystemAdminRole(permIds)        //初始化角色
		if role.ID != 0 {
			CreateSystemAdmin(role.ID) //初始化管理员
		}
	}
}

func GetAll(string, orderBy string, offset, limit int) *gorm.DB {
	if len(orderBy) > 0 {
		DB.Order(orderBy + "desc")
	} else {
		DB.Order("created_at desc")
	}
	if len(string) > 0 {
		DB.Where("name LIKE  ?", "%"+string+"%")
	}
	if offset > 0 {
		DB.Offset((offset - 1) * limit)
	}
	if limit > 0 {
		DB.Limit(limit)
	}
	return DB
}