.PHONY: up down restart build logs all front back

# Запуск всего проекта в фоновом режиме
up:
	docker compose up -d

# Остановка и удаление всех контейнеров (данные в БД и логи сохранятся)
down:
	docker compose down

# Перезапуск всех контейнеров
restart:
	docker compose restart

# Сборка всех образов с нуля (нужно, если поменяли Dockerfile или зависимости)
build:
	docker compose build

# Читать логи только бэкенда
logs:
	docker compose logs -f backend

# Обновить все
all:
	git pull
	git -C ../vietio-ui pull
	docker compose build
	docker compose down
	docker compose up -d

# Обновить только фронт
front:
	git -C ../vietio-ui pull
	docker compose build frontend-builder
	docker compose up --force-recreate -d frontend-builder

# Обновить только бэк
back:
	git pull
	docker compose build backend
	docker compose up --force-recreate -d backend
	docker compose restart nginx
