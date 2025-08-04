// Package dto 提供数据传输对象定义
// 用于前后端数据传输，避免直接暴露内部模型结构
package dto

import (
	"github.com/google/uuid"
)

// BaseModelDto 基础DTO，包含通用字段
// 所有DTO都应该继承此结构，确保一致性
type BaseModelDto struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt int64     `json:"created_at"`           // Unix 13位毫秒时间戳
	UpdatedAt int64     `json:"updated_at"`           // Unix 13位毫秒时间戳
	DeletedAt *int64    `json:"deleted_at,omitempty"` // Unix 13位毫秒时间戳
}
