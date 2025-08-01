// Package repository 提供数据访问层功能
// 负责与数据库的交互，实现数据持久化操作
package repository

import (
	"lemon-tree-core/internal/base"
	"lemon-tree-core/internal/models"

	"gorm.io/gorm"
)

// ApplicationInternalToolNetSearchConfigRepository ApplicationInternalToolNetSearchConfig 数据访问层接口
// 定义了 ApplicationInternalToolNetSearchConfig 模型的所有数据操作接口
// 继承自 BaseRepository，包含基本的增删改查功能
type ApplicationInternalToolNetSearchConfigRepository interface {
	base.BaseRepository[models.ApplicationInternalToolNetSearchConfig] // 继承基础仓库接口
}

// applicationInternalToolNetSearchConfigRepository ApplicationInternalToolNetSearchConfig 数据访问层实现
// 实现了 ApplicationInternalToolNetSearchConfigRepository 接口的所有方法
// 通过组合 baseRepository 来复用基础功能
type applicationInternalToolNetSearchConfigRepository struct {
	base.BaseRepository[models.ApplicationInternalToolNetSearchConfig] // 组合基础仓库实现
}

// NewApplicationInternalToolNetSearchConfigRepository 创建 ApplicationInternalToolNetSearchConfig Repository 实例
// 返回 ApplicationInternalToolNetSearchConfigRepository 接口的实现
// 参数：db - GORM 数据库连接实例
func NewApplicationInternalToolNetSearchConfigRepository(db *gorm.DB) ApplicationInternalToolNetSearchConfigRepository {
	return &applicationInternalToolNetSearchConfigRepository{
		BaseRepository: base.NewBaseRepository[models.ApplicationInternalToolNetSearchConfig](db),
	}
}
