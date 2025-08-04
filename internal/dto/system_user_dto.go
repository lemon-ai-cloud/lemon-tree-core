// Package dto 提供数据传输对象定义
package dto

// SystemUserDto 系统用户DTO
// 用于前后端数据传输，避免直接暴露内部模型结构
type SystemUserDto struct {
	BaseModelDto        // 继承基础DTO，包含 ID、时间戳等通用字段
	Name         string `json:"name"`   // 用户名字
	Number       string `json:"number"` // 用户账号
	Email        string `json:"email"`  // 用户邮箱
	// 注意：密码相关字段不包含在DTO中，避免安全风险
}

// SystemUserLoginDto 用户登录DTO
// 用于登录请求的数据传输
type SystemUserLoginDto struct {
	Number   string `json:"number" binding:"required"`   // 用户账号
	Password string `json:"password" binding:"required"` // 用户密码
}

// SystemUserSaveDto 用户保存DTO（创建或更新）
// 用于创建或更新用户时的数据传输
type SystemUserSaveDto struct {
	ID       string `json:"id,omitempty"`                 // 用户ID（更新时提供）
	Name     string `json:"name" binding:"required"`      // 用户名字
	Number   string `json:"number" binding:"required"`    // 用户账号
	Email    string `json:"email" binding:"required"`     // 用户邮箱
	Password string `json:"password" binding:"omitempty"` // 用户密码
}
