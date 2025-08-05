// Package service 提供业务逻辑层功能
// 负责处理业务规则和逻辑，连接 Handler 层和 Repository 层
package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"lemon-tree-core/internal/define"
	"lemon-tree-core/internal/models"
	"lemon-tree-core/internal/repository"
	"time"

	"github.com/google/uuid"
)

// UserService User 业务逻辑层接口
// 定义了 User 相关的所有业务操作接口
// 包含用户认证、会话管理和用户信息管理
type UserService interface {
	Login(ctx context.Context, number, password string) (*models.SystemUser, string, error) // 用户登录
	SaveUser(ctx context.Context, user *models.SystemUser) error                            // 保存用户（创建或更新）
	DeleteUser(ctx context.Context, id uuid.UUID) error                                     // 删除用户
	GetAllUsers(ctx context.Context) ([]*models.SystemUser, error)                          // 获取所有用户
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.SystemUser, error)              // 根据ID获取用户详情
	GetUserByToken(ctx context.Context, token string) (*models.SystemUser, error)           // 根据Token获取当前登录用户
	GetCurrentUser(ctx context.Context) (*models.SystemUser, error)                         // 获取当前登录用户
	Logout(ctx context.Context, token string) error                                         // 用户登出
}

// userService User 业务逻辑层实现
// 实现了 UserService 接口的所有方法
// 包含用户认证、会话管理和用户信息管理
type userService struct {
	userRepo    repository.SystemUserRepository        // 用户数据访问层接口
	sessionRepo repository.SystemUserSessionRepository // 会话数据访问层接口
}

// NewUserService 创建 User Service 实例
// 返回 UserService 接口的实现
// 参数：userRepo - 用户数据访问层接口，sessionRepo - 会话数据访问层接口
func NewUserService(userRepo repository.SystemUserRepository, sessionRepo repository.SystemUserSessionRepository) UserService {
	return &userService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

// hashPassword 使用SHA256加密密码
// 格式：SHA256(密码 + '_' + 盐)
// 参数：password - 原始密码，salt - 密码盐
// 返回：加密后的密码
func hashPassword(password, salt string) string {
	data := password + "_" + salt
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// Login 用户登录
// 验证用户账号密码，创建会话并返回Token
// 参数：ctx - 上下文，number - 用户账号，password - 用户密码
// 返回：用户对象、Token和错误信息
func (s *userService) Login(ctx context.Context, number, password string) (*models.SystemUser, string, error) {
	// 根据账号获取用户
	user, err := s.userRepo.GetByNumber(ctx, number)
	if err != nil {
		return nil, "", fmt.Errorf("用户不存在或账号错误")
	}

	// 验证密码：使用SHA256(密码 + '_' + 盐)进行验证
	hashedPassword := hashPassword(password, user.PasswordSalt)
	if user.Password != hashedPassword {
		return nil, "", fmt.Errorf("密码错误")
	}

	// 生成Token：sha256(随机UUID_用户ID_13位毫秒unix时间戳)
	randomUUID := uuid.New().String()
	userID := user.ID.String()
	timestamp := time.Now().UnixMilli()
	tokenInput := fmt.Sprintf("%s_%s_%d", randomUUID, userID, timestamp)

	hash := sha256.Sum256([]byte(tokenInput))
	token := hex.EncodeToString(hash[:])

	// 创建会话
	session := &models.SystemUserSession{
		Token:          token,
		UserID:         user.ID,
		LoginExpiredAt: time.Now().Add(24 * time.Hour), // 24小时过期
	}

	err = s.sessionRepo.Save(ctx, session)
	if err != nil {
		return nil, "", fmt.Errorf("创建会话失败: %w", err)
	}

	return user, token, nil
}

// SaveUser 保存用户（创建或更新）
// 如果用户存在则更新，不存在则创建
// 参数：ctx - 上下文，user - 要保存的用户对象
// 返回：错误信息
func (s *userService) SaveUser(ctx context.Context, user *models.SystemUser) error {
	// 检查用户是否已存在
	if user.ID != uuid.Nil {
		// 更新现有用户
		existingUser, err := s.userRepo.GetByID(ctx, user.ID)
		if err != nil {
			return fmt.Errorf("用户不存在: %w", err)
		}

		// 验证Number唯一性（如果Number发生变化）
		if existingUser.Number != user.Number {
			if _, err := s.userRepo.GetByNumber(ctx, user.Number); err == nil {
				return fmt.Errorf("用户账号已存在")
			}
		}

		// 验证Email唯一性（如果Email发生变化）
		if existingUser.Email != user.Email {
			if _, err := s.userRepo.GetByEmail(ctx, user.Email); err == nil {
				return fmt.Errorf("用户邮箱已存在")
			}
		}

		// 基于existingUser修改值
		existingUser.Name = user.Name
		existingUser.Number = user.Number
		existingUser.Email = user.Email

		// 处理密码（如果提供了新密码，需要加密）
		if user.Password != "" {
			existingUser.Password = hashPassword(user.Password, existingUser.PasswordSalt)
		}

		// 保存修改后的existingUser
		return s.userRepo.Save(ctx, existingUser)
	} else {
		// 创建新用户

		// 验证Number唯一性
		if _, err := s.userRepo.GetByNumber(ctx, user.Number); err == nil {
			return fmt.Errorf("用户账号已存在")
		}

		// 验证Email唯一性
		if _, err := s.userRepo.GetByEmail(ctx, user.Email); err == nil {
			return fmt.Errorf("用户邮箱已存在")
		}

		user.ID = uuid.New()
		// 生成密码盐
		if user.PasswordSalt == "" {
			user.PasswordSalt = uuid.New().String()
		}
		// 加密密码
		if user.Password != "" {
			user.Password = hashPassword(user.Password, user.PasswordSalt)
		}

		return s.userRepo.Save(ctx, user)
	}
}

// GetAllUsers 获取所有用户
// 获取所有用户的信息列表
// 参数：ctx - 上下文
// 返回：用户列表和错误信息
func (s *userService) GetAllUsers(ctx context.Context) ([]*models.SystemUser, error) {
	return s.userRepo.ListAll(ctx)
}

// GetUserByID 根据ID获取用户详情
// 根据 UUID 获取指定的用户信息
// 参数：ctx - 上下文，id - 用户的 UUID
// 返回：用户对象和错误信息
func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.SystemUser, error) {
	return s.userRepo.GetByID(ctx, id)
}

// GetCurrentUser 获取当前登录用户
// 获取当前登录用户信息
// 参数：ctx - 上下文
// 返回：当前登录用户对象
func (s *userService) GetCurrentUser(ctx context.Context) (*models.SystemUser, error) {
	user := ctx.Value(define.AppContextKeyCurrentUser)
	if user != nil {
		return user.(*models.SystemUser), nil
	}
	return nil, nil
}

// GetUserByToken 根据Token获取当前登录用户
// 验证Token有效性并返回当前登录用户信息
// 参数：ctx - 上下文，token - 用户Token
// 返回：用户对象和错误信息
func (s *userService) GetUserByToken(ctx context.Context, token string) (*models.SystemUser, error) {
	// 根据Token获取会话
	session, err := s.sessionRepo.GetByToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("无效的Token")
	}

	// 检查会话是否过期
	if time.Now().After(session.LoginExpiredAt) {
		// 删除过期会话
		s.sessionRepo.DeleteByID(ctx, session.ID)
		return nil, fmt.Errorf("会话已过期")
	}

	// 获取用户信息
	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("用户不存在")
	}

	return user, nil
}

// Logout 用户登出
// 删除用户的会话记录
// 参数：ctx - 上下文，token - 用户Token
// 返回：错误信息
func (s *userService) Logout(ctx context.Context, token string) error {
	// 根据Token获取会话
	session, err := s.sessionRepo.GetByToken(ctx, token)
	if err != nil {
		return fmt.Errorf("无效的Token")
	}

	// 删除会话
	return s.sessionRepo.DeleteByID(ctx, session.ID)
}

// DeleteUser 删除用户
// 删除指定用户及其所有会话记录
// 参数：ctx - 上下文，id - 要删除的用户ID
// 返回：错误信息
func (s *userService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	// 1. 检查用户是否存在
	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("用户不存在: %w", err)
	}

	// 2. 检查是否为当前登录用户（不能删除自己）
	currentUser, err := s.GetCurrentUser(ctx)
	if err == nil && currentUser != nil && currentUser.ID == id {
		return fmt.Errorf("不能删除当前登录的用户")
	}

	// 3. 检查是否为系统最后一个用户
	allUsers, err := s.userRepo.ListAll(ctx)
	if err != nil {
		return fmt.Errorf("获取用户列表失败: %w", err)
	}

	// 过滤掉已删除的用户
	var activeUsers []*models.SystemUser
	for _, u := range allUsers {
		if u.DeletedAt.Time.IsZero() {
			activeUsers = append(activeUsers, u)
		}
	}

	if len(activeUsers) <= 1 {
		return fmt.Errorf("系统至少需要保留一个用户")
	}

	// 4. 删除用户的所有会话记录
	err = s.sessionRepo.DeleteByUserID(ctx, id)
	if err != nil {
		return fmt.Errorf("删除用户会话失败: %w", err)
	}

	// 5. 删除用户
	err = s.userRepo.DeleteByID(ctx, id)
	if err != nil {
		return fmt.Errorf("删除用户失败: %w", err)
	}

	return nil
}
