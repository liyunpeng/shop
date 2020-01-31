package datamodels

import (
	"bytes"
	//"database/sql"
	//"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	//"time"
)

type UserG struct {
	gorm.Model
	Salt      string `gorm:"type:varchar(255)" json:"salt"`
	Username  string `gorm:"type:varchar(32)" json:"username"`
	Password  string `gorm:"type:varchar(200);column:password" json:"-"`
	Languages string `gorm:"type:varchar(200);column:languages" json:"languages"`
}

func (u UserG) TableName() string {
	return "gorm_user"
}


func (u *UserG)  CreateTable(db *gorm.DB) (s string) {
	var buffer bytes.Buffer
	/*
		gorm创建的表名默认为小写开头, 出现大写字符， 则会_分割， 以复数结尾， 可能加s,也可能加es
	*/
	// 分割， 名为credit_cards
	if !db.HasTable("usergs") {
		db.CreateTable(&UserG{})
		buffer.WriteString("gorm_user表创建成功\n")
	} else {
		buffer.WriteString("gorm_user表已存在，不再次创建\n")
	}

	return buffer.String()
}

func (u *UserG) Update(db *gorm.DB) (err error){
	db.Update(u)
	return nil
}

func (u *UserG) Delete(db *gorm.DB) (err error){
	db.Delete(u)
	return nil
}

func (u *UserG) FindById(db *gorm.DB, id int) ( err error){
	db.Where("id =?", id).First(u)
	return  nil
}

func (u *UserG) FindByName(db *gorm.DB, name string) ( err error){
	db.Where("username =?", name).First(u)
	return  nil
}

func (u *UserG) FindAll(db *gorm.DB, id int) ( err error){
	db.Find(u)
	return  nil
}

//func (u *UserG) FindCommentByAuthor(db *gorm.DB, author string)([]Comment){
//	//db := getDb()
//	//defer func() { db.Close() }()
//
//	comments := make([]Comment, 0)
//	/*
//		必须有过append 评论(comment)的帖子(post)才能在related得到评论（comment)。
//		posts和comments表并没有相互关联的列， 如下：
//		mysql> select * from posts;
//		posts表：
//		+----+---------------+---------------+---------------------+
//		| id | content       | author        | createed_at         |
//		+----+---------------+---------------+---------------------+
//		|  1 | content value | author value1 | 0000-00-00 00:00:00 |
//		comments表：
//		mysql> select * from comments;
//		+----+---------------+---------------+---------+---------------------+
//		| id | content       | author        | post_id | createed_at         |
//		+----+---------------+---------------+---------+---------------------+
//		|  1 | content value | author value1 |       2 | 0000-00-00 00:00:00 |
//		gorm维护了这种关联关系。
//	*/
//	db.Where("author =?", author).Last(&u)
//	fmt.Println("post:", u)
//
//	db.Model(&u).Related(&comments)
//	fmt.Println("comments:", comments)
//	return  comments
//}


//////////////////////////////////////////
func (u *UserG) Insert(db *gorm.DB) (err error){

	db.Create(u)


	return nil
}