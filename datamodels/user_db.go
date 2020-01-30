package datamodels

import (
	"bytes"
	//"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	//"time"
)

func getDb()(db *gorm.DB){
	defer func(){
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容
		}
	}()

	//TODO？ 全局db没用
	/*
		链接localhost数据库， 用户名root, 密码root
	*/
	db, err := gorm.Open(
		"mysql", "root:root@/gotest?charset=utf8&parseTime=True&loc=Local")
	if err == nil {
		fmt.Println("open db sucess", db)

	} else {
		fmt.Println("open db error ", err)
		panic("数据库错误")
	}

	return db
}

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


func (post *UserG)  CreateTable() (s string) {
	db := getDb()
	defer func() { db.Close() }()

	var buffer bytes.Buffer
	/*
		gorm创建的表名默认为小写开头, 出现大写字符， 则会_分割， 以复数结尾， 可能加s,也可能加es
	*/
	// 分割， 名为credit_cards
	if !db.HasTable("usergs") {
		db.CreateTable(&UserG{})
		buffer.WriteString("Post表创建成功\n")
	} else {
		buffer.WriteString("Post表已存在，不再次创建\n")
	}

	return buffer.String()
}



func (u *UserG) Update() (err error){
	db := getDb()
	defer func() { db.Close() }()
	db.Update(u)
	return nil
}

func (u *UserG) Delete() (err error){
	db := getDb()
	defer func() { db.Close() }()
	db.Delete(u)
	return nil
}

func (u *UserG) FindById(id int) ( err error){
	db := getDb()
	defer func() { db.Close() }()
	db.Where("id =?", id).First(u)
	return  nil
}

func (u *UserG) FindAll(id int) ( err error){
	db := getDb()
	defer func() { db.Close() }()
	/*
		获取所有记录
		以下等价于： SELECT * FROM users;
	*/
	db.Find(u)
	return  nil
}

func (u *UserG) FindCommentByAuthor(author string)([]Comment){
	db := getDb()
	defer func() { db.Close() }()

	comments := make([]Comment, 0)
	/*
		必须有过append 评论(comment)的帖子(post)才能在related得到评论（comment)。
		posts和comments表并没有相互关联的列， 如下：
		mysql> select * from posts;
		posts表：
		+----+---------------+---------------+---------------------+
		| id | content       | author        | createed_at         |
		+----+---------------+---------------+---------------------+
		|  1 | content value | author value1 | 0000-00-00 00:00:00 |
		comments表：
		mysql> select * from comments;
		+----+---------------+---------------+---------+---------------------+
		| id | content       | author        | post_id | createed_at         |
		+----+---------------+---------------+---------+---------------------+
		|  1 | content value | author value1 |       2 | 0000-00-00 00:00:00 |
		gorm维护了这种关联关系。
	*/
	db.Where("author =?", author).Last(&u)
	fmt.Println("post:", u)

	db.Model(&u).Related(&comments)
	fmt.Println("comments:", comments)
	return  comments
}


//////////////////////////////////////////
func (u *UserG) Insert() (err error){
	db := getDb()
	defer func() { db.Close() }()

	/*
		这里插入的post记录， 没有append comment记录， 所有post的related方法不会得到comment记录
	*/
	//db.Create(post)

	/*
			等同于，
		    INSERT INTO products (name, code) VALUES ("name", "code") ON CONFLICT;
	*/
	db.Set("gorm:insert_option", "ON CONFLICT").Create(u)

	return nil
}