#### orm利弊
##### 好处
* 防止直接拼接sql语句引入sql注入漏洞
* 方便对modle进行统一管理
* 专注业务，加速开发
##### 坏处
* 开发者与最终的sql语句隔了一层orm，因此可能会不慎引入烂sql
* 依赖于orm的成熟度，无法进行一些「复杂」的查询。当然，复杂的查询一大半都是应该从设计上规避的 

#### gorm类型与数据库表的映射  
一个数据表用一个文件来定义  
文件里要提供相应的数据类型的定义和数据操作的方法  
类型里面需要有个model基本类  
类型里的每个成员都要大写开头， tag标记必须有gorm开头的字符创， tag 设定了在mysql数据表里的数据类型的定义
由于grom是使用的orm映射，所以需要定义要操作的表的model，
在go中需要定义一个struct， struct的名字就是对应数据库中的表名，  
注意gorm查找struct名对应数据库中的表名的时候会默认把你的struct中的大写字母转换为小写并加上“s”，  
所以可以加上 db.SingularTable(true) 让grom转义struct名字的时候不用加上s。  


定义model，即struct, 定于struct时我们可以只定义我们需要从数据库中取回的特定字段：
gorm在转义表名的时候会把struct的大写字母(首字母除外) 替换成“_”，
所以下面的"XzAutoServerConf "会转义成数数据库中对应“xz_auto_server_conf”的表名, 
对应的字段名的查找会先按照tag里面的名称去里面查找，  
如果没有定义标签则按照struct定义的字段查找，  
查找的时候struct字段中的大写会被转义成“ ”,
例“GroupZone”会去查找表中的group_zone字段

```gotemplate
type XzAutoServerConf struct {
    GroupZone string `gorm:"column:group_zone"`
    ServerId int
    OpenTime string
    ServerName string
    Status int
}
```
#### 建议在数据库创建表，不用gorm创建表 
我是提前在数据库中创建好表的然后再用grom去查询的，也可以用gorm去创建表，  
我感觉还是直接在数据库上创建，修改表字段的操作方便，grom只用来查询和更新数据。
假设数据库中的表已经创建好，下面是数据库中的建表语句：
```mysql
CREATE TABLE `xz_auto_server_conf` ( 
    `id` int(11) NOT NULL AUTO_INCREMENT, 
    `group_zone` varchar(32) NOT NULL COMMENT '大区例如：wanba,changan,aiweiyou,360', 
    `server_id` int(11) DEFAULT '0' COMMENT '区服id', 
    `server_name` varchar(255) NOT NULL COMMENT '区服名称', 
    `open_time` varchar(64) DEFAULT NULL COMMENT '开服时间', 
    `service` varchar(30) DEFAULT NULL COMMENT '环境，test测试服，formal混服，wb玩吧', 
    `username` varchar(100) DEFAULT NULL COMMENT 'data管理员名称', 
    `submit_date` datetime DEFAULT NULL COMMENT '记录提交时间', 
    `status` tinyint(2) DEFAULT '0' COMMENT '状态，0未处理，1已处理，默认为0', 
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```

#### model中定义获取单个记录行和多个记录行的方法
数据操作的方法，基本围绕增删改查这些方法，这里关注查询方法 
* 查询单个记录行
方法的传入参数是表示要操作的数据的id等，
方法返回的是是对象指针，如
```
func GetUserById(id uint) *User {
	user := new(User)
	IsNotFound(Db.Where("username= ?", username).First(user).Error)
	return user
}
```
* 查询多个记录行
方法返回的是对象指针数组，如
```
func GetAllUsers(name, orderBy string, offset, limit int) []*User {
	var users []*User
	q := GetAll(name, orderBy, offset, limit)
	if err := q.Find(&users).Error; err != nil {
		color.Red(fmt.Sprintf("GetAllUserErr:%s \n ", err))
		return nil
	}
	return users
}
```

#### 数据库连接
```go
package  main
type ConnInfo struct {
    MyUser string
    Password string
    Host string
    Port int
    Db string
}
func main () {
    cn := ConnInfo{"root","123456", "127.0.0.1", 3306, "xd_data"}  
    db := DbConn(cn.MyUser,cn.Password,cn.Host,cn.Db,cn.Port)  
    defer db.Close()
    var rows []api.XzAutoServerConf
    //select
    db.Where("status=?", 0).Select([]string{"group_zone", "server_id", "open_time", "server_name"}).Find(&rows)
    //update 
    err := db.Model(&rows).Where("server_id=?", 80).Update("status", 1).Error
    if err !=nil {
        fmt.Println(err)
    }
        fmt.Println(rows)
}
```

#### 图片保存
图片保存到数据库通常有两种做法：一是把图片转化成二进制形式，然后保存到数据库中；二是将图片保存的路径存储在数据库中。
由于图片一般都在几M，所以当以二进制流的形式保存到数据库中时，在高并发的情况下，会加重数据库的负担。（至于更详细的原因，后期文章会做分享）因此大多数会选择第二种做法。


#### 客户端频繁tcp连接，导致端口耗尽
* 底层报错 error：cannot assign requested address
* 原因
并发场景下 client 频繁请求端口建立tcp连接导致端口被耗尽
* 解决方案
让TIME-WAIT状态的sockets快速回收
root执行
```shell script
# 开启对于TCP时间戳的支持,若该项设置为0，则下面一项设置不起作用
sysctl -w net.ipv4.tcp_timestamps=1 
# 表示开启TCP连接中TIME-WAIT sockets的快速回收
sysctl -w net.ipv4.tcp_tw_recycle=1 
```
