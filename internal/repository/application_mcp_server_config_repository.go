// Package repository 提供数据访问层功能
// 负责与数据库的交互，实现数据持久化操作
package repository

import (
	"lemon-tree-core/internal/base"
	"lemon-tree-core/internal/models"

	"gorm.io/gorm"
)

// ApplicationMcpConfigConfigRepository ApplicationMcpConfigConfig 数据访问层接口
// 定义了 ApplicationMcpConfigConfig 模型的所有数据操作接口
// 继承自 BaseRepository，包含基本的增删改查功能
type ApplicationMcpConfigConfigRepository interface {
	base.BaseRepository[models.ApplicationMcpConfigConfig] // 继承基础仓库接口
}

// applicationMcpConfigConfigRepository ApplicationMcpConfigConfig 数据访问层实现
// 实现了 ApplicationMcpConfigConfigRepository 接口的所有方法
// 通过组合 baseRepository 来复用基础功能
type applicationMcpConfigConfigRepository struct {
	base.BaseRepository[models.ApplicationMcpConfigConfig] // 组合基础仓库实现
}

// NewApplicationMcpConfigConfigRepository 创建 ApplicationMcpConfigConfig Repository 实例
// 返回 ApplicationMcpConfigConfigRepository 接口的实现
// 参数：db - GORM 数据库连接实例
func NewApplicationMcpConfigConfigRepository(db *gorm.DB) ApplicationMcpConfigConfigRepository {
	return &applicationMcpConfigConfigRepository{
		BaseRepository: base.NewBaseRepository[models.ApplicationMcpConfigConfig](db),
	}
}
