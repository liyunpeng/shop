package sys

import (
	"github.com/kataras/iris/v12"
	"shop/client"
	"shop/config"
	"shop/models"
	"shop/transformer"
	"shop/validates"
	"shop/web/routes"
)

func createEtcdKv(transformConfiguration *transformer.Conf) {
	etcdKV := make(map[string]string, 2)

	etcdKV["192.168.0.1"] = `
[
	{
		"topic":"nginx_log",
		"log_path":"/Users/admin1/goworkspace/shop/log1.txt",
		"service":"test_service",
		"send_rate":1000
	},
		
	{
		"topic":"nginx_log1",
		"log_path":"/Users/admin1/goworkspace/shop/log2.txt",
		"service":"test_service1",
		"send_rate":1000
	},

	{
		"topic":"nginx_log",
		"log_path":"D:\\goworkspace\\shop\\log1.txt",
		"service":"test_service1",
		"send_rate":1000
	}
]`
	etcdKV["192.168.0.2"] = `
[
	{
		"topic":"nginx_log",
		"log_path":"/Users/admin1/goworkspace/shop/123.txt",
		"service":"test_service",
		"send_rate":2000
	}
]`
	for k, v := range etcdKV {
		client.EtcdClientInsance.PutKV(k, v)
	}
}

func CreateOrderTestData() {

	models.OrderDelete()
	aul := &validates.CreateOrderRequest{
		Username:    "aa",
		Title:       "订单标题1",
		Price:       1000.00,
		Description: "订单描述1",
		Status:      "待发货",
		ImagePath:   "/public/images/classify-ph03.png",
	}
	models.CreateOrder(aul)

	aul = &validates.CreateOrderRequest{
		Username:    "aa",
		Title:       "订单标题2",
		Price:       3000.00,
		Description: "订单描述2",
		Status:      "已发货",
		ImagePath:   "/public/images/classify-ph01.png",
	}
	models.CreateOrder(aul)

	aul = &validates.CreateOrderRequest{
		Username:    "aa",
		Title:       "订单标题5",
		Price:       3000.00,
		Description: "订单描述3",
		Status:      "待接收",
		ImagePath:   "/public/images/classify-ph02.png",
	}
	models.CreateOrder(aul)
}

func CreateGoodsTestData() {
	goods := &models.Goods{
		Name:        "goods1",
		Description: "desription1",
		Price:       10,
		Type:        "营养快线",
		ImagePath:   "/public/images/classify-ph03.png",
		Stock:       100,
	}
	models.CreateGoods(goods)

	goods = &models.Goods{
		Name:        "goods2",
		Description: "desription2",
		Price:       100,
		Type:        "营养快线",
		ImagePath:   "/public/images/classify-ph03.png",
		Stock:       200,
	}
	models.CreateGoods(goods)

}

/**
 * 创建系统权限
 * @return
 */
func CreateSystemAdminPermission(perms []*validates.PermissionRequest) []uint {
	var permIds []uint
	for _, perm := range perms {
		p := models.GetPermissionByNameAct(perm.Name, perm.Act)
		if p.ID != 0 {
			continue
		}
		pp := models.CreatePermission(perm)
		permIds = append(permIds, pp.ID)
	}
	return permIds
}

func CreateSystemData(app *iris.Application) {
	//if rc.App.CreateSysData == 1 {
	if true {
		apiRoutes := routes.GetRoutes(app)
		permIds := CreateSystemAdminPermission(apiRoutes) //初始化权限
		role := CreateSystemAdminRole(permIds)        //初始化角色
		if role.ID != 0 {
			CreateSystemAdmin(role.ID) //初始化管理员
		}
		CreateOrderTestData()
		CreateGoodsTestData()
		createEtcdKv(config.TransformConfiguration)
	}

}

/**
*创建系统管理员
*@return   *models.AdminRoleTranform api格式化后的数据格式
 */
func CreateSystemAdminRole(permIds []uint) *models.Role {
	aul := &validates.RoleRequest{
		Name:        "admin",
		DisplayName: "超级管理员",
		Description: "超级管理员",
	}

	//role := models.GetRoleByName(aul.Name)
	//if role.ID == 0 {
	//  role.ID == 0 表示角色不存在 有时casbin没有配置好，就错过了在casbin中为此用户创建的权限策略
	//  所以为角色添加casbin权限的操作，不应该放在CreateRole里面。

	//	return models.CreateRole(aul, permIds)
	//} else {
	//	return role
	//}

	return models.CreateRole(aul, permIds)
	//return role
}
func CreateSystemAdmin(roleId uint) *models.User {
	aul := &validates.CreateUpdateUserRequest{
		Username: "admin@126.com",
		Password: "123",
		Name:     "admin",
		RoleIds:  []uint{roleId},
	}

	if models.IsUserExist(aul.Username) == false {
		return models.CreateUser(aul)
	} else {
		user := models.UserFindByName(aul.Username)
		return user
	}
	//user := UserFindByName(aul.Username)
	//if user.ID == 0 {
	//	return CreateUser(aul)
	//} else {
	//	return user
	//}
}
