.PHONY: swagger-api swagger-admin run-api run-admin

# 生成 App API 文档
swagger-api:
	swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal --parseGoList=false

# 生成 Admin 文档
swagger-admin:
	swag init -g cmd/api/main.go -o docs/admin --instanceName admin --parseDependency --parseInternal --parseGoList=false
# 启动 API 服务
run-api:
	go run cmd/api/main.go

# 启动 Admin 服务
run-admin:
	go run cmd/admin/main.go

# 整理代码
tidy:
	go mod tidy