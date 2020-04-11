package models

import (
	"shop/validates"
)

func CreateOrderTestData() {

	OrderDelete()
	aul := &validates.CreateOrderRequest{
		Username: "aa",
		Title: "订单标题1",
		Price: 1000.00,
		Description: "订单描述1",
		Status: "待发货",
		ImagePath: "/public/images/classify-ph03.png",
	}
	CreateOrder(aul)

	aul = &validates.CreateOrderRequest{
		Username: "aa",
		Title: "订单标题2",
		Price: 3000.00,
		Description: "订单描述2",
		Status: "已发货",
		ImagePath: "/public/images/classify-ph01.png",
	}
	CreateOrder(aul)

	aul = &validates.CreateOrderRequest{
		Username: "aa",
		Title: "订单标题5",
		Price: 3000.00,
		Description: "订单描述3",
		Status: "待接收",
		ImagePath: "/public/images/classify-ph02.png",
	}
	CreateOrder(aul)
}

func CreateGoodsTestData(){
	goods := &Goods{
		Name:        "goods1",
		Description: "desription1",
		Price:       10,
		Type:	"营养快线",
		ImagePath:   "/public/images/classify-ph03.png",
		Stock:       100,
	}
	CreateGoods(goods)

	goods = &Goods{
		Name:        "goods2",
		Description: "desription2",
		Price:       100,
		Type:	"营养快线",
		ImagePath:   "/public/images/classify-ph03.png",
		Stock:       200,
	}
	CreateGoods(goods)

}
/**
 * 创建系统权限
 * @return
 */
func CreateSystemAdminPermission(perms []*validates.PermissionRequest) []uint {
	var permIds []uint
	for _, perm := range perms {
		p := GetPermissionByNameAct(perm.Name, perm.Act)
		if p.ID != 0 {
			continue
		}
		pp := CreatePermission(perm)
		permIds = append(permIds, pp.ID)
	}
	return permIds
}

func CreateSystemData(perms []*validates.PermissionRequest) {
	//if rc.App.CreateSysData == 1 {
	if true {
		permIds := CreateSystemAdminPermission(perms) //初始化权限
		role := CreateSystemAdminRole(permIds)        //初始化角色
		if role.ID != 0 {
			CreateSystemAdmin(role.ID) //初始化管理员
		}
		CreateOrderTestData()
		CreateGoodsTestData()
	}
}

/**
*创建系统管理员
*@return   *models.AdminRoleTranform api格式化后的数据格式
 */
func CreateSystemAdminRole(permIds []uint) *Role {
	aul := &validates.RoleRequest{
		Name:        "admin",
		DisplayName: "超级管理员",
		Description: "超级管理员",
	}

	role := GetRoleByName(aul.Name)
	if role.ID == 0 {
		return CreateRole(aul, permIds)
	} else {
		return role
	}
}
func CreateSystemAdmin(roleId uint ) *User {
	aul := &validates.CreateUpdateUserRequest{
		Username: "admin@126.com",
		Password: "123",
		Name:    "admin",
		RoleIds:  []uint{roleId},
	}

	if ( IsUserExist( aul.Username) == false ) {
		return CreateUser(aul)
	}else{
		user := UserFindByName(aul.Username)
		return user
	}
	//user := UserFindByName(aul.Username)
	//if user.ID == 0 {
	//	return CreateUser(aul)
	//} else {
	//	return user
	//}
}
