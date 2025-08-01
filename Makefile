# Makefile - 项目构建和管理脚本
# 提供常用的构建、测试、部署等命令

# .PHONY 声明伪目标，避免与同名文件冲突
.PHONY: build run test clean

# build - 构建项目
# 编译 Go 代码生成可执行文件
build:
	go build -o lemon-tree-core .

# run - 运行项目
# 直接运行 Go 代码，不生成可执行文件
run:
	go run main.go

# test - 测试项目
# 运行所有测试用例
test:
	go test ./...

# clean - 清理构建文件
# 删除生成的可执行文件
clean:
	rm -f lemon-tree-core

# deps - 安装依赖
# 下载和整理 Go 模块依赖
deps:
	go mod tidy

# fmt - 格式化代码
# 使用 gofmt 格式化所有 Go 代码
fmt:
	go fmt ./...

# lint - 代码检查
# 使用 golangci-lint 进行代码质量检查
lint:
	golangci-lint run

# docs - 生成文档
# 启动 godoc 服务器查看代码文档
docs:
	godoc -http=:6060

# migrate - 数据库迁移
# 运行数据库迁移脚本
migrate:
	go run main.go migrate

# dev - 开发模式运行
# 在开发模式下运行应用程序
dev:
	go run main.go 