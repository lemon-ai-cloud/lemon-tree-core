// Package models 提供应用程序的数据模型定义
package models

import (
	"github.com/google/uuid"
	"lemon-tree-core/internal/base"
)

// ApplicationStorageConfig 应用存储配置
type ApplicationStorageConfig struct {
	base.BaseModel           // 继承基础模型，包含 ID、时间戳等通用字段
	Type           string    `json:"type" gorm:"type:varchar(64);not null;comment:存储类型"`
	ApplicationID  uuid.UUID `json:"application_id" gorm:"type:char(36);not null;comment:所属应用ID"`
	// Type为file_system
	RootPath string `json:"root_path" gorm:"type:varchar(512);not null;comment:文件系统根路径"`
	// Type为s3
	Endpoint   string `json:"endpoint" gorm:"type:varchar(512);not null;comment:S3存储桶endpoint"`
	Region     string `json:"region" gorm:"type:varchar(512);not null;comment:S3存储桶区域"`
	BucketName string `json:"bucket_name" gorm:"type:varchar(512);not null;comment:S3存储桶名称"`
	SecretId   string `json:"secret_id" gorm:"type:varchar(512);not null;comment:S3存储安全ID"`
	SecretKey  string `json:"secret_key" gorm:"type:varchar(512);not null;comment:S3存储密钥"`
	KeyPrefix  string `json:"key_prefix" gorm:"type:varchar(512);not null;comment:S3存储文件key前缀，用于设置根路径"`
}

// TableName 指定数据库表名
// 返回该模型对应的数据库表名
func (ApplicationStorageConfig) TableName() string {
	return "ltc_application_storage_config"
}
