// Package dto 提供数据传输对象定义
// 用于在不同层之间传递数据，确保数据格式的一致性
package dto

// ApplicationStorageConfigDto 应用存储配置数据传输对象
// 用于在业务逻辑层和HTTP处理层之间传递数据
type ApplicationStorageConfigDto struct {
	ID            string `json:"id"`             // 主键ID
	Type          string `json:"type"`           // 存储类型
	ApplicationID string `json:"application_id"` // 所属应用ID
	// Type为file_system时的字段
	RootPath string `json:"root_path"` // 文件系统根路径
	// Type为s3时的字段
	Endpoint   string `json:"endpoint"`    // S3存储桶endpoint
	Region     string `json:"region"`      // S3存储桶区域
	BucketName string `json:"bucket_name"` // S3存储桶名称
	SecretId   string `json:"secret_id"`   // S3存储安全ID
	SecretKey  string `json:"secret_key"`  // S3存储密钥
	KeyPrefix  string `json:"key_prefix"`  // S3存储文件key前缀
	CreatedAt  string `json:"created_at"`  // 创建时间
	UpdatedAt  string `json:"updated_at"`  // 更新时间
}

// SaveApplicationStorageConfigRequest 保存应用存储配置请求
// 用于前端保存存储配置的请求数据
type SaveApplicationStorageConfigRequest struct {
	ID            *string `json:"id,omitempty"`   // 主键ID（更新时提供）
	Type          string  `json:"type"`           // 存储类型
	ApplicationID string  `json:"application_id"` // 所属应用ID
	// Type为file_system时的字段
	RootPath string `json:"root_path"` // 文件系统根路径
	// Type为s3时的字段
	Endpoint   string `json:"endpoint"`    // S3存储桶endpoint
	Region     string `json:"region"`      // S3存储桶区域
	BucketName string `json:"bucket_name"` // S3存储桶名称
	SecretId   string `json:"secret_id"`   // S3存储安全ID
	SecretKey  string `json:"secret_key"`  // S3存储密钥
	KeyPrefix  string `json:"key_prefix"`  // S3存储文件key前缀
}

// SingleApplicationStorageConfigResponse 单个应用存储配置响应
// 用于返回单个存储配置数据
type SingleApplicationStorageConfigResponse struct {
	ApplicationStorageConfig ApplicationStorageConfigDto `json:"application_storage_config"` // 存储配置数据
}
