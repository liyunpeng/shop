// file: controllers/user_controller.go

package controllers

import (
	//"datamodels"
	"fmt"
	"shop/datamodels"
	"shop/services"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
)

// UserGController is our /user controller.
// UserGController is responsible to handle the following requests:
// GET  			/user/register
// POST 			/user/register
// GET 				/user/login
// POST 			/user/login
// GET 				/user/me
// All HTTP Methods /user/logout
type UserGController struct {
	// context is auto-binded by Iris on each request,
	// remember that on each incoming request iris creates a new UserGController each time,
	// so all fields are request-scoped by-default, only dependency injection is able to set
	// custom fields like the Service which is the same for all requests (static binding)
	// and the Session which depends on the current context (dynamic binding).
	Ctx iris.Context

	// Our UserService, it's an interface which
	// is binded from the main application.
	Service services.UserGService

	// Session, binded using dependency injection from the main.go.
	Session *sessions.Session
}

const userIDKey1 = "UserID"

func (c *UserGController) getCurrentUserID() int64 {
	userID := c.Session.GetInt64Default(userIDKey1, 0)
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
	//if c.isLoggedIn() {
	//	c.logout()
	//}

	return registerStaticView1
}

// PostRegister handles POST: http://localhost:8080/user/register.
func (c *UserGController) PostRegister() mvc.Result {
	fmt.Println("响应注册提交")
	// get firstname, username and password from the form.
	var (
		//firstname = c.Ctx.FormValue("firstname")
		username  = c.Ctx.FormValue("username")
		password  = c.Ctx.FormValue("password")
	)

	// create the new user, the password will be hashed by the service.
	//u, err := c.Service.InsertUserg(datamodels.UserG{
	//TODO : InsertUserg 返回值， 新增记录的id如何获取
	 c.Service.InsertUserg(datamodels.UserG{
		Username:  username,
		Password: password,
	})

	// set the user's id to this session even if err != nil,
	// the zero id doesn't matters because .getCurrentUserID() checks for that.
	// If err != nil then it will be shown, see below on mvc.Response.Err: err.
	//c.Session.Set(userIDKey1, u.ID)

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

	u, found := c.Service.GetByUsernameAndPassword(username, password)


	if !found {
		c.Service.CreateUsergTable()
		fmt.Println("进入注册页面")
		return mvc.Response{
			Path: "/userg/register",
		}
	}
	fmt.Println("u.Username=", u.Username)
	//c.Session.Set(userIDKey1, u.Username)

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
