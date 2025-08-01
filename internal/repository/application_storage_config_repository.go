// Package repository 提供数据访问层功能
// 负责与数据库的交互，实现数据持久化操作
package repository

import (
	"lemon-tree-core/internal/base"
	"lemon-tree-core/internal/models"

	"gorm.io/gorm"
)

// ApplicationStorageConfigRepository ApplicationStorageConfig 数据访问层接口
// 定义了 ApplicationStorageConfig 模型的所有数据操作接口
// 继承自 BaseRepository，包含基本的增删改查功能
type ApplicationStorageConfigRepository interface {
	base.BaseRepository[models.ApplicationStorageConfig] // 继承基础仓库接口
}

// applicationStorageConfigRepository ApplicationStorageConfig 数据访问层实现
// 实现了 ApplicationStorageConfigRepository 接口的所有方法
// 通过组合 baseRepository 来复用基础功能
type applicationStorageConfigRepository struct {
	base.BaseRepository[models.ApplicationStorageConfig] // 组合基础仓库实现
}

// NewApplicationStorageConfigRepository 创建 ApplicationStorageConfig Repository 实例
// 返回 ApplicationStorageConfigRepository 接口的实现
// 参数：db - GORM 数据库连接实例
func NewApplicationStorageConfigRepository(db *gorm.DB) ApplicationStorageConfigRepository {
	return &applicationStorageConfigRepository{
		BaseRepository: base.NewBaseRepository[models.ApplicationStorageConfig](db),
	}
}
