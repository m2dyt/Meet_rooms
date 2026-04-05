.PHONY: up down build seed test cover lint swagger clean migrate-up migrate-down

# Запуск проекта (пересборка и запуск)
up:
	docker-compose up --build

# Запуск в фоновом режиме
up-d:
	docker-compose up --build -d

# Остановка и удаление контейнеров с volumes
down:
	docker-compose down -v

# Сборка без запуска
build:
	docker-compose build --no-cache

# Наполнение БД тестовыми данными
seed:
	docker-compose exec app ./booking-app -seed

# Запуск тестов (юнит + интеграционные)
test:
	docker-compose exec app go test -v -cover ./...

# Покрытие тестами (генерация HTML отчёта)
cover:
	docker-compose exec app go test -coverprofile=coverage.out ./...
	docker-compose exec app go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Запуск линтера (golangci-lint)
lint:
	docker-compose exec app golangci-lint run

# Генерация Swagger документации (локально, на хосте)
swagger:
	swag init -g cmd/server/main.go --parseDependency --parseInternal

# Полная очистка (контейнеры, образы, volumes)
clean:
	docker-compose down -v
	docker system prune -af --volumes

# Выполнить миграции (создать таблицы) через migrate-контейнер
migrate-up:
	docker-compose run --rm migrate -path /migrations -database "postgres://postgres:postgres@postgres:5432/booking?sslmode=disable" up

# Откатить миграции (удалить таблицы)
migrate-down:
	docker-compose run --rm migrate -path /migrations -database "postgres://postgres:postgres@postgres:5432/booking?sslmode=disable" down

# Логи приложения
logs:
	docker-compose logs -f app