ENV_FILE = ./.env
COMPOSE = docker compose -f ./docker-compose.yml --env-file=$(ENV_FILE)

.PHONY: help
help:
	@echo "\033[33mДоступные команды:\033[0m"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


#--------------- ОСНОВНЫЕ КОМАНДЫ ---------------#

.PHONY: env
env: ## Создать .env файл из .env.example
	@cp .env.example .env

.PHONY: build
build: ## Собрать Docker-образы
	@$(COMPOSE) up -d --build

.PHONY: up
up: ## Запустить контейнеры Docker
	@$(COMPOSE) up -d --remove-orphans

.PHONY: down
down: ## Остановить и удалить контейнеры Docker
	@$(COMPOSE) down

.PHONY: clean
clean: ## Полностью очистить окружение (удалить контейнеры и образы)
	@$(COMPOSE) down --rmi all

.PHONY: clean-volumes
clean-volumes: ## Очистить окружение с удалением томов
	@$(COMPOSE) down -v

.PHONY: restart
restart: ## Перезапустить контейнеры Docker
	@$(COMPOSE) restart


#--------------- ПРОВЕРКА КАЧЕСТВА КОДА ---------------#

.PHONY: lint
lint: ## Анализ кода с помощью golangci-lint (поиск ошибок)
	@$(COMPOSE) exec app golangci-lint run

.PHONY: fix
fix: ## Автоматическое исправление стиля кода (golangci-lint)
	@$(COMPOSE) exec app golangci-lint run --fix
	@$(COMPOSE) exec app gofumpt -l -w .
	@$(COMPOSE) exec app gci write . --skip-generated -s standard -s default
