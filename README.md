# bondy-backend

# 启动数据库
docker compose up -d

# 停止数据库（数据不丢失）
docker compose down

# 重新生成 Swagger 文档
swag init -g cmd/api/main.go -o docs

# 启动 API 服务
go run cmd/api/main.go

# 整理依赖
go mod tidy

# 查看数据库容器日志
docker compose logs postgres