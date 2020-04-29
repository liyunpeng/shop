package test

import (
	"fmt"
	"github.com/facebookgo/inject"
	"shop/logger"
	"testing"
)

/*
软件构建的核心就是管理复杂度

解耦组件之间的依赖关系，避免手动配置每个组件的依赖关系。
 */
type DBEngine struct {
	Name string
}

type UserDB struct {
	Db *DBEngine `inject:""`
}

type UserService struct {
	Db *UserDB `inject:""`
}

type App struct {
	Name string
	User1 *UserService `inject:""`
}

func (a *App) Create() string {
	return "create app, in db name:" + a.User1.Db.Db.Name+" app name :"+ a.Name
}

type Object struct {
	App *App
}

func Init(useInject bool) *Object {
	db := DBEngine{Name: "db1"}

	var app App
	if useInject{
		logger.Info.Println("使用依赖注入:")
		var g inject.Graph
		app = App{Name: "go-app"}
		_ = g.Provide(
			&inject.Object{Value: &app},
			&inject.Object{Value: &db},
		)
		_ = g.Populate()

	}else{
		logger.Info.Println("不使用依赖注入:")
		app = App{
			Name: "go-app",
			User1: &UserService{
				Db: &UserDB{
					Db: &db,
				},
			},
		}

	}
	return &Object{
		App: &app,
	}





}

func TestInject( t *testing.T) {

	logger.InitCustLogger()
	obj := Init(true)
	fmt.Println(obj.App.Create())


	obj1 := Init(false)
	fmt.Println(obj1.App.Create())
}
