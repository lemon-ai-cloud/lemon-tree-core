// Package base 提供基础组件功能
package base

import (
	"context"
	"reflect"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseRepository 基础仓库接口
// 定义了通用的数据访问操作接口
// 包含基本的增删改查功能和动态查询功能
type BaseRepository[T any] interface {
	Create(ctx context.Context, entity *T) error           // 创建实体
	Update(ctx context.Context, entity *T) error           // 更新实体
	Save(ctx context.Context, entity *T) error             // 保存实体（upsert）
	DeleteByID(ctx context.Context, id uuid.UUID) error    // 根据ID删除实体
	ListAll(ctx context.Context) ([]*T, error)             // 获取所有实体
	GetByID(ctx context.Context, id uuid.UUID) (*T, error) // 根据ID获取实体
	Query(ctx context.Context, query *T) ([]*T, error)     // 动态查询实体
}

// baseRepository 基础仓库实现
// 实现了 BaseRepository 接口的所有方法
// 使用泛型支持任意类型的实体
type baseRepository[T any] struct {
	db *gorm.DB // GORM 数据库连接实例
}

// NewBaseRepository 创建基础仓库实例
// 返回 BaseRepository 接口的实现
// 参数：db - GORM 数据库连接实例
func NewBaseRepository[T any](db *gorm.DB) BaseRepository[T] {
	return &baseRepository[T]{db: db}
}

// Create 创建实体
// 将新的实体信息保存到数据库中
// 参数：ctx - 上下文，entity - 要创建的实体对象
// 返回：错误信息
func (r *baseRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// Update 更新实体
// 更新数据库中指定实体的信息
// 参数：ctx - 上下文，entity - 要更新的实体对象
// 返回：错误信息
func (r *baseRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

// Save 保存实体（upsert）
// 如果实体存在则更新，不存在则创建
// 参数：ctx - 上下文，entity - 要保存的实体对象
// 返回：错误信息
func (r *baseRepository[T]) Save(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

// DeleteByID 根据ID删除实体
// 根据 UUID 删除指定的实体（软删除）
// 参数：ctx - 上下文，id - 要删除的实体 UUID
// 返回：错误信息
func (r *baseRepository[T]) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(new(T), "id = ?", id).Error
}

// ListAll 获取所有实体
// 从数据库中获取所有实体的信息列表（排除已删除的）
// 参数：ctx - 上下文
// 返回：实体列表和错误信息
func (r *baseRepository[T]) ListAll(ctx context.Context) ([]*T, error) {
	var entities []*T
	err := r.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&entities).Error
	return entities, err
}

// GetByID 根据ID获取实体
// 根据 UUID 查找并返回指定的实体信息（排除已删除的）
// 参数：ctx - 上下文，id - 实体的 UUID
// 返回：实体对象和错误信息
func (r *baseRepository[T]) GetByID(ctx context.Context, id uuid.UUID) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// Query 动态查询实体
// 根据传入的查询对象进行动态查询
// 只查询非零值的字段，排除已删除的
// 参数：ctx - 上下文，query - 查询条件对象
// 返回：匹配的实体列表和错误信息
func (r *baseRepository[T]) Query(ctx context.Context, query *T) ([]*T, error) {
	var entities []*T

	// 构建查询条件
	db := r.db.WithContext(ctx).Where("deleted_at IS NULL")

	// 使用反射获取查询对象的非零值字段
	queryMap := r.buildQueryMap(query)

	// 应用查询条件
	for field, value := range queryMap {
		db = db.Where(field+" = ?", value)
	}

	// 执行查询
	err := db.Find(&entities).Error
	return entities, err
}

// buildQueryMap 构建查询映射
// 使用反射获取结构体中非零值的字段
// 参数：query - 查询对象
// 返回：字段名到值的映射
func (r *baseRepository[T]) buildQueryMap(query *T) map[string]interface{} {
	queryMap := make(map[string]interface{})

	// 获取结构体的反射值
	v := reflect.ValueOf(query).Elem()
	t := v.Type()

	// 遍历结构体的所有字段
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// 获取字段名（转换为数据库字段名）
		fieldName := r.getDBFieldName(fieldType)

		// 检查字段是否为零值
		if !field.IsZero() {
			queryMap[fieldName] = field.Interface()
		}
	}

	return queryMap
}

// getDBFieldName 获取数据库字段名
// 从结构体标签中提取数据库字段名
// 参数：fieldType - 结构体字段类型
// 返回：数据库字段名
func (r *baseRepository[T]) getDBFieldName(fieldType reflect.StructField) string {
	// 尝试从 gorm 标签中获取字段名
	if gormTag := fieldType.Tag.Get("gorm"); gormTag != "" {
		// 这里可以解析 gorm 标签，简化处理直接使用 json 标签
	}

	// 使用 json 标签作为字段名
	if jsonTag := fieldType.Tag.Get("json"); jsonTag != "" {
		return jsonTag
	}

	// 如果没有标签，使用字段名（转换为下划线格式）
	return r.camelToSnake(fieldType.Name)
}

// camelToSnake 驼峰命名转下划线命名
// 将驼峰命名的字符串转换为下划线命名
// 参数：s - 驼峰命名字符串
// 返回：下划线命名字符串
func (r *baseRepository[T]) camelToSnake(s string) string {
	var result string
	for i, char := range s {
		if i > 0 && char >= 'A' && char <= 'Z' {
			result += "_"
		}
		result += string(char)
	}
	return result
}
