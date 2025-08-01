// Package main 应用程序的主入口包
// 负责启动应用程序和初始化依赖注入容器
package main

import (
	"context"
	"lemon-tree-core/internal/core"
	"log"
)

// main 应用程序的主函数
// 程序的入口点，负责启动整个应用程序
func main() {
	// 创建依赖注入容器
	// 配置所有组件的依赖关系和生命周期
	app := core.NewContainer()

	// 启动应用程序
	// 开始依赖注入容器的生命周期管理
	if err := app.Start(context.Background()); err != nil {
		log.Fatal("Failed to start application:", err)
	}

	// 等待应用程序结束
	// 阻塞主线程，直到应用程序被终止
	<-app.Done()
}
