package models

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	"shop/validates"
)

type Permission struct {
	gorm.Model
	Name        string `gorm:"not null VARCHAR(191)"`
	DisplayName string `gorm:"VARCHAR(191)"`
	Description string `gorm:"VARCHAR(191)"`
	Act         string `gorm:"VARCHAR(191)"`
}

/**
 * 通过 id 获取 permission 记录
 * @method GetPermissionById
 * @param  {[type]}       permission  *Permission [description]
 */
func GetPermissionById(id uint) *Permission {
	permission := new(Permission)
	IsNotFound(DB.Where("id = ?", id).First(permission).Error)

	return permission
}

/**
 * 通过 name 获取 permission 记录
 * @method GetPermissionByName
 * @param  {[type]}       permission  *Permission [description]
 */
func GetPermissionByNameAct(name, act string) *Permission {
	permission := new(Permission)
	IsNotFound(DB.Where("name = ?", name).Where("act = ?", act).First(permission).Error)
	return permission
}

/**
 * 通过 id 删除权限
 * @method DeletePermissionById
 */
func DeletePermissionById(id uint) {
	u := new(Permission)
	u.ID = id

	if err := DB.Delete(u).Error; err != nil {
		color.Red(fmt.Sprintf("DeletePermissionByIdError:%s \n", err))
	}
}

/**
 * 获取所有的权限
 * @method GetAllPermissions
 * @param  {[type]} name string [description]
 * @param  {[type]} orderBy string [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func GetAllPermissions(name, orderBy string, offset, limit int) (permissions []*Permission) {
	if err := GetAll(name, orderBy, offset, limit).Find(&permissions).Error; err != nil {
		color.Red(fmt.Sprintf("GetAllPermissionsError:%s \n", err))
	}

	return
}

/**
 * 创建
 * @method CreatePermission
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func CreatePermission(aul *validates.PermissionRequest) (permission *Permission) {
	permission = new(Permission)
	permission.Name = aul.Name
	permission.DisplayName = aul.DisplayName
	permission.Description = aul.Description
	permission.Act = aul.Act
	if err := DB.Create(permission).Error; err != nil {
		color.Red(fmt.Sprintf("CreatePermissionError:%s \n", err))
	}
	return
}

/**
 * 更新
 * @method UpdatePermission
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func UpdatePermission(pj *validates.PermissionRequest, id uint) (permission *Permission) {
	permission = new(Permission)
	permission.ID = id

	if err := DB.Model(&permission).Updates(pj).Error; err != nil {
		color.Red(fmt.Sprintf("UpdatePermissionError:%s \n", err))
	}

	return
}


