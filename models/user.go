package models

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/jameskeane/bcrypt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"shop/logger"
	"shop/validates"
	"strconv"
	"time"
)

type User struct {
	gorm.Model
	Salt        string `gorm:"type:varchar(255)" json:"salt"`
	Username    string `gorm:"unique_index" json:"username"`
	Name        string `gorm:"not null VARCHAR(191)"`
	Address     string `gorm:"type:varchar(191)"`
	Password    string `gorm:"type:varchar(200);column:password" json:"-"`
	Phonenumber string `gorm:"type:varchar(200);column:phonenumber" json:"phonenumber"`
	Level       string `gorm:"type:varchar(200);column:level" json:"level"`
}

func (u User) TableName() string {
	return "gorm_user"
}

func UserCreateTable() (s string) {
	var buffer bytes.Buffer
	/*
		gorm创建的表名默认为小写开头, 出现大写字符， 则会_分割， 以复数结尾， 可能加s,也可能加es
	*/
	if !DB.HasTable("gorm_user") {
		DB.CreateTable(&User{})
		buffer.WriteString("gorm_user表创建成功\n")
	} else {
		buffer.WriteString("gorm_user表已存在，不再次创建\n")
	}

	return buffer.String()
}

func UserInsert(user *User) {
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(user.Password, salt)

	user.Password = hash
	//user = &User{
	//	Username: aul.Username,
	//	Password: hash,
	//	Name:     aul.Name,
	//}

	DB.Create(user)
}

func UserFindByName(name string) *User {
	user := new(User)
	tx := DB.Where("username =?", name).Find(&user)
	if tx.Error == nil {
		return user
	} else {
		return nil
	}
}

func IsUserExist(name string) bool {
	var count int
	if DB.Model(&User{}).Where("username =?", name).Count(&count).Error != nil {
		panic("IsExist 异常")
	} else {
		if count == 0 {
			return false
		} else if count == 1 {
			return true
		} else {
			panic("IsExist count 异常")
		}
	}
}

func UserUpdate(user *User) (err error) {
	DB.Model(&User{}).Where("username =?", user.Username).Update(user)
	return nil
}

func UserDeleteByName(username string) {
	user := &User{
		Username: username,
	}
	DB.Delete(user)
}

func UserFindById(id int64) *User {
	user := new(User)
	DB.Where("id =?", id).First(user)
	return user
}

func UserFindAll() []*User {
	//var users []*User
	//DB.DB().Ping()
	//usersa := make([]*User, 100)
	//
	//for i := 0; i < 100; i++ {
	//	usersa[i] = new(User)
	//}
	var usersa []*User
	//TODO: find 查找全部咋用，
	DB.Model(&User{}).Find(&usersa)
	//DB.Model(&User{}).First(usersa[0])
	return usersa
}

func IsNotFound(err error) {
	if ok := errors.Is(err, gorm.ErrRecordNotFound); !ok && err != nil {
		color.Red(fmt.Sprintf("error :%v \n ", err))
	}
}

func UserAdminCheckLogin(username string) *User {
	user := new(User)
	IsNotFound(DB.Where("username = ?", username).First(user).Error)

	return user
}

func CheckLogin(username, password string) (response Token, status bool, msg string) {
	user := UserAdminCheckLogin(username)
	if user.ID == 0 {
		msg = "用户不存在"
		return
	} else {
		if ok := bcrypt.Match(password, user.Password); ok {
			token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"exp": time.Now().Add(time.Hour * time.Duration(1)).Unix(),
				"iat": time.Now().Unix(),
			})
			tokenString, _ := token.SignedString([]byte("HS2JDFKhu7Y1av7b"))
			oauthToken := new(OauthToken)
			oauthToken.Token = tokenString
			oauthToken.UserId = user.ID
			oauthToken.Secret = "secret"
			oauthToken.Revoked = false
			oauthToken.ExpressIn = time.Now().Add(time.Hour * time.Duration(1)).Unix()
			oauthToken.CreatedAt = time.Now()
			response = oauthToken.OauthTokenCreate()
			status = true
			msg = "登陆成功"

			return

		} else {
			msg = "用户名或密码错误"
			return
		}
	}
}
func addRoles(uj *validates.CreateUpdateUserRequest, user *User) {
	if len(uj.RoleIds) > 0 {
		userId := strconv.FormatUint(uint64(user.ID), 10)
		if _, err = Enforcer.DeleteRolesForUser(userId); err != nil {
			color.Red(fmt.Sprintf("CreateUserErr:%s \n ", err))
		}

		for _, roleId := range uj.RoleIds {
			roleId := strconv.FormatUint(uint64(roleId), 10)
			if _, err = Enforcer.AddRoleForUser(userId, roleId); err != nil {
				color.Red(fmt.Sprintf("CreateUserErr:%s \n ", err))
			}
		}
	} else {
		logger.Info.Println("没后为用户添加角色")
	}
}

func CreateUser(aul *validates.CreateUpdateUserRequest) (user *User) {
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(aul.Password, salt)

	user = &User{
		Username: aul.Username,
		Password: hash,
		Name:     aul.Name,
	}

	if err := DB.Create(user).Error; err != nil {
		color.Red(fmt.Sprintf("CreateUserErr:%s \n ", err))
	}

	addRoles(aul, user)

	return
}
