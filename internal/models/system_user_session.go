// Package models 提供应用程序的数据模型定义
package models

import (
	"github.com/google/uuid"
	"lemon-tree-core/internal/base"
	"time"
)

// SystemUserSession 系统用户会话
type SystemUserSession struct {
	base.BaseModel // 继承基础模型，包含 ID、时间戳等通用字段
	// Token生成算法：sha256(随机UUID_用户ID_13位毫秒unix时间戳)
	Token          string    `json:"token" gorm:"type:varchar(512);not null;comment:Token"`
	UserID         uuid.UUID `json:"user_id" gorm:"type:char(36);not null;comment:用户ID"`
	LoginExpiredAt time.Time `json:"login_expired_at" gorm:"type:datetime;not null;comment:登录过期时间"`
}

// TableName 指定数据库表名
// 返回该模型对应的数据库表名
func (SystemUserSession) TableName() string {
	return "ltc_system_user_session"
}
