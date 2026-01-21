# MagTrade Makefile
# 專案管理命令集

.PHONY: dev prod build test lint clean logs backup

# 開發環境
dev:
	cd docker && docker compose -f docker-compose.dev.yml up -d

# 生產環境
prod:
	cd docker && docker compose up -d

# 重建並部署（零停機）
build:
	cd docker && docker compose up -d --build --no-deps backend

# 完全重建
rebuild:
	cd docker && docker compose down && docker compose up -d --build

# 運行測試
test:
	go test -v -race ./internal/model/... ./internal/pkg/...

# 測試覆蓋率
coverage:
	go test -coverprofile=coverage.out ./internal/model/... ./internal/pkg/...
	go tool cover -html=coverage.out -o coverage.html

# 代碼檢查
lint:
	golangci-lint run ./...

# 查看日誌
logs:
	cd docker && docker compose logs -f --tail=100

# 查看後端日誌
logs-backend:
	cd docker && docker compose logs -f backend

# 資料庫備份
backup:
	@mkdir -p backups
	docker exec mt-postgres pg_dump -U postgres magtrade > backups/magtrade_$$(date +%Y%m%d_%H%M%S).sql
	@echo "備份完成: backups/magtrade_$$(date +%Y%m%d_%H%M%S).sql"

# 資料庫還原（使用: make restore FILE=backups/xxx.sql）
restore:
	docker exec -i mt-postgres psql -U postgres magtrade < $(FILE)

# 清理未使用的 Docker 資源
clean:
	docker system prune -f
	docker volume prune -f

# 進入資料庫
db:
	docker exec -it mt-postgres psql -U postgres -d magtrade

# 進入 Redis
redis:
	docker exec -it mt-redis redis-cli

# 健康檢查
health:
	curl -s https://gcptw.yukiyuki.cfd/health | jq

# 幫助
help:
	@echo "可用命令："
	@echo "  make dev      - 啟動開發環境"
	@echo "  make prod     - 啟動生產環境"
	@echo "  make build    - 零停機更新後端"
	@echo "  make rebuild  - 完全重建所有服務"
	@echo "  make test     - 運行單元測試"
	@echo "  make lint     - 代碼靜態檢查"
	@echo "  make logs     - 查看服務日誌"
	@echo "  make backup   - 備份資料庫"
	@echo "  make db       - 進入資料庫 CLI"
	@echo "  make health   - 檢查服務健康狀態"
