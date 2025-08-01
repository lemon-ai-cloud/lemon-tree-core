// Package models 提供应用程序的数据模型定义
package models

import "lemon-tree-core/internal/base"

// SystemUser 系统用户
type SystemUser struct {
	base.BaseModel        // 继承基础模型，包含 ID、时间戳等通用字段
	Name           string `json:"name" gorm:"type:varchar(64);not null;comment:用户名字"`
	Number         string `json:"number" gorm:"type:varchar(64);not null;comment:用户账号"`
	Email          string `json:"email" gorm:"type:varchar(128);not null;comment:用户邮箱"`
	Password       string `json:"password" gorm:"type:varchar(512);not null;comment:用户密码"`
	PasswordSalt   string `json:"password_salt" gorm:"type:varchar(512);not null;comment:用户密码盐"`
}

// TableName 指定数据库表名
// 返回该模型对应的数据库表名
func (SystemUser) TableName() string {
	return "ltc_system_user"
}
