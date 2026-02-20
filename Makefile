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

# Посмотреть статус всех контейнеров
ps:
	docker compose ps

# Читать логи ТОЛЬКО бэкенда
logs-back:
	docker compose logs -f backend

# Обновить ТОЛЬКО фронтенд
front:
	docker compose build frontend-builder
	docker compose up --force-recreate -d frontend-builder

# Обновить ТОЛЬКО бэкенд
back:
	docker compose build backend
	docker compose up --force-recreate -d backend
