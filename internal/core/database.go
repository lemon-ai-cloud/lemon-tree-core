// Package database 提供数据库连接和配置功能
// 负责建立数据库连接、配置 GORM 和自动迁移表结构
package core

import (
	"fmt"
	"lemon-tree-core/internal/config"
	"lemon-tree-core/internal/models"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDatabase 创建数据库连接
// 根据配置信息建立 MySQL 数据库连接
// 配置 GORM 日志和自动迁移表结构
// 参数：config - 应用程序配置
// 返回：GORM 数据库连接实例和错误信息
func NewDatabase(config *config.Config) (*gorm.DB, error) {
	// 构建数据库连接字符串（DSN）
	// 格式：username:password@tcp(host:port)/database?charset=charset&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		config.Database.Username, // 数据库用户名
		config.Database.Password, // 数据库密码
		config.Database.Host,     // 数据库主机地址
		config.Database.Port,     // 数据库端口号
		config.Database.Database, // 数据库名称
		config.Database.Charset,  // 数据库字符集
	)

	// 配置 GORM 日志记录器
	// 使用结构化日志记录 SQL 查询和错误信息
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // 日志输出到标准输出
		logger.Config{
			SlowThreshold:             time.Second, // 慢查询阈值（1秒）
			LogLevel:                  logger.Info, // 日志级别（Info）
			IgnoreRecordNotFoundError: true,        // 忽略记录未找到错误
			Colorful:                  true,        // 启用彩色输出
		},
	)

	// 使用 MySQL 驱动打开数据库连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger, // 使用配置的日志记录器
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 自动迁移表结构
	// 根据模型定义自动创建或更新数据库表
	// 迁移所有模型对应的表
	if err := db.AutoMigrate(
		&models.Application{},                            // 应用表
		&models.SystemUser{},                             // 系统用户表
		&models.SystemUserSession{},                      // 系统用户会话表
		&models.ApplicationLlm{},                         // 应用模型表
		&models.ChatAgent{},                              // 聊天智能体表
		&models.ChatAgentConversation{},                  // 聊天智能体会话表
		&models.ChatAgentMessage{},                       // 聊天智能体消息表
		&models.ChatAgentAttachment{},                    // 聊天智能体附件表
		&models.ChatAgentApiKey{},                        // 聊天智能体API Key表
		&models.ChatConversation{},                       // 聊天会话表
		&models.LlmProvider{},                            // 大语言模型供应商表
		&models.ApplicationStorageConfig{},               // 应用存储配置表
		&models.ApplicationInternalToolNetSearchConfig{}, // 应用内部工具网络搜索配置表
		&models.ApplicationMcpConfigConfig{},             // 应用MCP服务器配置表
	); err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	return db, nil
}
