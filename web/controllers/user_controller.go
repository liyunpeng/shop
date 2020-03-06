package controllers
import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"shop/models"
	"shop/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"net/http"
	validates "shop/validate"
)

type UserGController struct {
	// context is auto-binded by Iris on each request,
	// remember that on each incoming request iris creates a new UserGController each time,
	// so all fields are request-scoped by-default, only dependency injection is able to set
	// custom fields like the Service which is the same for all requests (static binding)
	// and the Session which depends on the current context (dynamic binding).
	// iris会为每一个request请求自动绑定到iris.context
	// 每一个请求的用户都会为之创建一个新的UserGController 控制器
	//静态绑定： 服务是对应所有请求的
	// 动态绑定： session 依赖于context的当前状态
	Ctx iris.Context

	// Our UserService, it's an interface which
	// is binded from the main application.
	// 服务是主程序绑定过来的
	Service services.UserService

	// Session, binded using dependency injection from the main.go.
	// 使用main.go里的依赖注入， 将session绑定到该控制器
	Session *sessions.Session
}

const usergIDKey = "UserID"

func (c *UserGController) getCurrentUserID() int64 {
	userID := c.Session.GetInt64Default(usergIDKey, 0)
	return userID
}

func (c *UserGController) isLoggedIn() bool {
	return c.getCurrentUserID() > 0
}

func (c *UserGController) logout() {
	c.Session.Destroy()
}

var registerStaticView1 = mvc.View{
	Name: "register.html",
	Data: iris.Map{"Title": "User Registration"},
}

// GetRegister handles GET: http://localhost:8080/user/register.
func (c *UserGController) GetRegister() mvc.Result {
	if c.isLoggedIn() {
		c.logout()
	}

	return registerStaticView1
}

// PostRegister handles POST: http://localhost:8080/user/register.
func (c *UserGController) PostRegister() mvc.Result {
	fmt.Println("响应注册提交")
	var (
		username  = c.Ctx.FormValue("username")
		password  = c.Ctx.FormValue("password")
	)

	// create the new user, the password will be hashed by the service.
	//  创建一个用户， 密码会被被service通过摘要算法计算出一个哈希值，将此哈希值存入数据库用户表
	//u, err := c.Service.InsertUserg(datamodels.User{
	//TODO : InsertUserg 返回值， 新增记录的id如何获取
	 c.Service.InsertUser(&models.User{
		Username:  username,
		Password: password,
	})

	// set the user's id to this session even if err != nil,
	// the zero id doesn't matters because .getCurrentUserID() checks for that.
	// If err != nil then it will be shown, see below on mvc.Response.Err: err.
	//c.Session.Set(usergIDKey, u.ID)

	return mvc.Response{
		// if not nil then this error will be shown instead.
		//Err: err,
		Path: "/self",
		// When redirecting from POST to GET request you -should- use this HTTP status code,
		// however there're some (complicated) alternatives if you
		// search online or even the HTTP RFC.
		// Status "See Other" RFC 7231, however iris can automatically fix that
		// but it's good to know you can set a custom code;
		// Code: 303,
	}
}

var loginStaticView1 = mvc.View{
	Name: "login.html",
	Data: iris.Map{"Title": "User Login"},
}

// GetLogin handles GET: http://localhost:8080/user/login.
func (c *UserGController) GetLogin() mvc.Result {
	//if c.isLoggedIn() {
	//	// if it's already logged in then destroy the previous session.
	//	c.logout()
	//}

	return loginStaticView1
}

// PostLogin handles POST: http://localhost:8080/user/login.
func (c *UserGController) PostLogin() mvc.Result {
	var (
		username = c.Ctx.FormValue("username")
		password = c.Ctx.FormValue("password")
	)

	user, found := c.Service.GetByUsernameAndPassword(username, password)

	if !found {
		c.Service.CreateUserTable()
		fmt.Println("进入注册页面")
		return mvc.Response{
			Path: "/userg/register",
		}
	}
	fmt.Println("user.Username=", user.Username)
	c.Session.Set(usergIDKey, user.Username)

	return mvc.Response{
		Path: "/self",
	}
}

// GetMe handles GET: http://localhost:8080/user/me.
func (c *UserGController) GetMe() mvc.Result {
	//if !c.isLoggedIn() {
	//	// if it's not logged in then redirect user to the login page.
	//	return mvc.Response{Path: "/user/login"}
	//}
	//
	//u, found := c.Service.GetByID(c.getCurrentUserID())
	//if !found {
	//	// if the  session exists but for some reason the user doesn't exist in the "database"
	//	// then logout and re-execute the function, it will redirect the client to the
	//	// /user/login page.
	//	c.logout()
	//	return c.GetMe()
	//}

	return mvc.View{
		Name: "user/me.html",
		Data: iris.Map{
			"Title": "Profile of ",
			"User":  "u",
		},
	}
}

// AnyLogout handles All/Any HTTP Methods for: http://localhost:8080/user/logout.
func (c *UserGController) AnyLogout() {
	if c.isLoggedIn() {
		c.logout()
	}

	c.Ctx.Redirect("/user/login")
}

func ApiUserGetAll(c iris.Context) {
	users := models.UserFindAll()
	c.StatusCode(http.StatusOK)
	//v1 := []string{"one", "two", "three"}
	_, _ = c.JSON(ApiResource(true, users, "RepsonseJson message  "))
}


func ApiUserGetById(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	user := models.UserFindById(id)
	ctx.StatusCode(http.StatusOK)
	//v1 := []string{"one", "two", "three"}
	_, _ = ctx.JSON(ApiResource(true, user, "RepsonseJson message  "))
}


func  ApiUserPost(ctx iris.Context) {
	fmt.Println("ApiUserPost is called")
	aul := new(validates.CreateUpdateUserRequest)
	if err := ctx.ReadJSON(aul); err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(false, nil, err.Error()))
		return
	}
	u := models.User{}
	u.Username = aul.Username
	u.Password = aul.Password
	u.Phonenumber = aul.Phonenumber
	u.Level = aul.Level
	models.UserInsert(&u)

	_, _ = ctx.JSON(ApiResource(true, nil, "成功添加数据行"))
}

func  ApiUserUpdate(ctx iris.Context) {
	fmt.Println("ApiUserUpdate is called")
	aul := new(validates.CreateUpdateUserRequest)
	if err := ctx.ReadJSON(aul); err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(false, nil, err.Error()))
		return
	}
	u := models.User{}
	u.Username = aul.Username
	u.Password = aul.Password
	models.UserUpdate(&u)

	_, _ = ctx.JSON(ApiResource(true, nil, "成功修改数据行"))
}

func  ApiUserInsertOrUpdate(ctx iris.Context) {
	fmt.Println("ApiUserInsertOrUpdate is called")
	aul := new(validates.CreateUpdateUserRequest)
	if err := ctx.ReadJSON(aul); err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(false, nil, err.Error()))
		return
	}
	u := models.User{}

	u.Password = aul.Password
	u.Phonenumber = aul.Phonenumber
	u.Level = aul.Level

	if models.UserFindByName(u.Username) != nil {
		models.UserUpdate(&u)
	}else{
		u.Username = aul.Username
		models.UserInsert(&u)
	}

	_, _ = ctx.JSON(ApiResource(true, nil, "成功更新数据行"))
}

func  ApiDatabaseCreate(ctx iris.Context) {
	str := models.UserCreateTable()
	str += models.AuthtokenCreateTable()
	user := models.User{
		Username : "admin@126.com",
		Password : "123",
	}

	models.UserInsert(&user)

	_, _ = ctx.JSON(ApiResource(true, nil, str))

}

func UserLogin(ctx iris.Context) {
	aul := new(validates.LoginRequest)

	if err := ctx.ReadJSON(aul); err != nil {  //{"username":"username", "password","passrd"}
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(false, nil, err.Error()))
		return
	}

	err := validates.Validate.Struct(*aul)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(validates.ValidateTrans) {
			if len(e) > 0 {
				ctx.StatusCode(iris.StatusOK)
				_, _ = ctx.JSON(ApiResource(false, nil, e))
				return
			}
		}
	}

	ctx.Application().Logger().Infof("%s 登录系统", aul.Username)
	ctx.StatusCode(iris.StatusOK)
	response, status, msg := models.CheckLogin(aul.Username, aul.Password)
	_, _ = ctx.JSON(ApiResource(status, response, msg))
	return

}
