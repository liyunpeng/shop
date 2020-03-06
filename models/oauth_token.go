package models

import (
	"bytes"
	"github.com/jinzhu/gorm"
)

type OauthToken struct {
	gorm.Model

	Token     string `gorm:"not null default '' comment('Token') VARCHAR(191)"`
	UserId    uint   `gorm:"not null default '' comment('UserId') VARCHAR(191)"`
	Secret    string `gorm:"not null default '' comment('Secret') VARCHAR(191)"`
	ExpressIn int64  `gorm:"not null default 0 comment('是否是标准库') BIGINT(20)"`
	Revoked   bool
}

type Token struct {
	Token string `json:"access_token"`
}

func (u OauthToken) TableName() string {
	return "auth_token"
}

func AuthtokenCreateTable() (s string) {
	var buffer bytes.Buffer
	if !DB.HasTable("auth_token") {
		DB.CreateTable(&OauthToken{})
		buffer.WriteString("auth_token表创建成功\n")
	} else {
		buffer.WriteString("auth_token表已存在，不再次创建\n")
	}

	return buffer.String()
}


/**
 * oauth_token
 * @method OauthTokenCreate
 */
func (ot *OauthToken) OauthTokenCreate() (response Token) {
	DB.Create(ot)
	response = Token{ot.Token}

	return
}

/**
 * 通过 token 获取 access_token 记录
 * @method GetOauthTokenByToken
 * @param  {[type]}       token string [description]
 */
func GetOauthTokenByToken(token string) (ot *OauthToken) {
	ot = new(OauthToken)
	DB.Where("token =  ?", token).First(&ot)
	return
}

/**
 * 通过 user_id 更新 oauth_token 记录
 * @method UpdateOauthTokenByUserId
 *@param  {[type]}       user  *OauthToken [description]
 */
func UpdateOauthTokenByUserId(userId uint) (ot *OauthToken) {
	DB.Model(ot).Where("revoked = ?", false).
		Where("user_id = ?", userId).
		Updates(map[string]interface{}{"revoked": true})
	return
}
