// Package converter 提供模型与DTO之间的转换功能
// 负责将内部模型转换为前端DTO，以及将DTO转换为内部模型
package converter

import (
	"lemon-tree-core/internal/dto"
	"lemon-tree-core/internal/models"

	"github.com/google/uuid"
)

// SystemUserModelToSystemUserDto 将 SystemUser 模型转换为 SystemUserDto
// 参数：user - 系统用户模型
// 返回：SystemUserDto - 系统用户DTO
func SystemUserModelToSystemUserDto(user *models.SystemUser) *dto.SystemUserDto {
	if user == nil {
		return nil
	}

	var deletedAt *int64
	if !user.DeletedAt.Time.IsZero() {
		timestamp := user.DeletedAt.Time.UnixMilli()
		deletedAt = &timestamp
	}

	return &dto.SystemUserDto{
		BaseModelDto: dto.BaseModelDto{
			ID:        user.ID,
			CreatedAt: user.CreatedAt.UnixMilli(),
			UpdatedAt: user.UpdatedAt.UnixMilli(),
			DeletedAt: deletedAt,
		},
		Name:   user.Name,
		Number: user.Number,
		Email:  user.Email,
	}
}

// SystemUserModelListToSystemUserDtoList 将 SystemUser 模型列表转换为 SystemUserDto 列表
// 参数：users - 系统用户模型列表
// 返回：[]*SystemUserDto - 系统用户DTO列表
func SystemUserModelListToSystemUserDtoList(users []*models.SystemUser) []*dto.SystemUserDto {
	if users == nil {
		return nil
	}

	result := make([]*dto.SystemUserDto, len(users))
	for i, user := range users {
		result[i] = SystemUserModelToSystemUserDto(user)
	}
	return result
}

// SystemUserSaveDtoToSystemUserModel 将 SystemUserSaveDto 转换为 SystemUser 模型
// 参数：userDto - 系统用户保存DTO
// 返回：*SystemUser - 系统用户模型
func SystemUserSaveDtoToSystemUserModel(userDto *dto.SystemUserSaveDto) *models.SystemUser {
	if userDto == nil {
		return nil
	}

	user := &models.SystemUser{
		Name:     userDto.Name,
		Number:   userDto.Number,
		Email:    userDto.Email,
		Password: userDto.Password,
	}

	// 如果提供了ID，则解析UUID
	if userDto.ID != "" {
		if id, err := uuid.Parse(userDto.ID); err == nil {
			user.ID = id
		}
	}

	return user
}
