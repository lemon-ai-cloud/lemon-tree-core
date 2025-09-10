// Package config 提供应用程序的配置管理功能
// 使用 Viper 库来读取和管理配置文件
package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config 应用程序的主配置结构体
// 包含服务器配置、数据库配置和AI客户端配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`   // 服务器配置
	Database DatabaseConfig `mapstructure:"database"` // 数据库配置
	AI       AIConfig       `mapstructure:"ai"`       // AI客户端配置
}

// ServerConfig 服务器配置结构体
// 定义服务器的端口和运行模式
type ServerConfig struct {
	Port string `mapstructure:"port"` // 服务器监听端口，如 ":8080"
	Mode string `mapstructure:"mode"` // 服务器运行模式，如 "debug" 或 "release"
}

// DatabaseConfig 数据库配置结构体
// 定义数据库连接的相关参数
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`     // 数据库主机地址
	Port     string `mapstructure:"port"`     // 数据库端口号
	Username string `mapstructure:"username"` // 数据库用户名
	Password string `mapstructure:"password"` // 数据库密码
	Database string `mapstructure:"database"` // 数据库名称
	Charset  string `mapstructure:"charset"`  // 数据库字符集
}

// AIConfig AI客户端配置结构体
// 定义AI客户端的相关参数
type AIConfig struct {
	Type    string `mapstructure:"type"`     // AI客户端类型，如 "openai", "ollama"
	APIKey  string `mapstructure:"api_key"`  // API密钥
	BaseURL string `mapstructure:"base_url"` // 基础URL（可选）
	Model   string `mapstructure:"model"`    // 模型名称（可选）
}

// AppConfig 全局配置变量
// 用于在整个应用程序中访问配置信息
var AppConfig *Config

// LoadConfig 加载配置文件
// 优先从环境变量读取配置，支持 .env 文件
// 返回配置对象，如果加载失败则程序退出
func LoadConfig() *Config {
	// 加载 .env 文件（如果存在）
	if err := godotenv.Load(); err != nil {
		// 如果 .env 文件不存在，继续使用环境变量
		log.Println("No .env file found, using environment variables")
	}

	// 设置环境变量前缀
	viper.SetEnvPrefix("")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// 设置默认值
	setDefaults()

	// 创建配置对象
	AppConfig = &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", ":8080"),
			Mode: getEnv("SERVER_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			Username: getEnv("DB_USERNAME", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Database: getEnv("DB_DATABASE", "lemon_tree_core"),
			Charset:  getEnv("DB_CHARSET", "utf8mb4"),
		},
		AI: AIConfig{
			Type:    getEnv("AI_TYPE", "openai"),
			APIKey:  getEnv("AI_API_KEY", ""),
			BaseURL: getEnv("AI_BASE_URL", ""),
			Model:   getEnv("AI_MODEL", ""),
		},
	}

	return AppConfig
}

// setDefaults 设置默认配置值
func setDefaults() {
	// 服务器默认配置
	viper.SetDefault("server.port", ":8080")
	viper.SetDefault("server.mode", "debug")

	// 数据库默认配置
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "3306")
	viper.SetDefault("database.username", "root")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.database", "lemon_tree_core")
	viper.SetDefault("database.charset", "utf8mb4")

	// AI客户端默认配置
	viper.SetDefault("ai.type", "openai")
	viper.SetDefault("ai.api_key", "")
	viper.SetDefault("ai.base_url", "")
	viper.SetDefault("ai.model", "")
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
