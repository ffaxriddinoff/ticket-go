.DEFAULT_GOAL := help

GO ?= go
ARTISAN = $(GO) run . artisan
SSH ?= ssh
SCP ?= scp

-include deploy.mk

DEPLOY_SSH ?=
DEPLOY_PORT ?= 22
DEPLOY_PATH ?= /var/www/ticket-backend
DEPLOY_ARCH ?= amd64
DEPLOY_SERVICE ?= ticket-backend

export DEPLOY_SSH DEPLOY_PORT DEPLOY_PATH DEPLOY_ARCH DEPLOY_SERVICE

.PHONY: help run build tidy fmt vet artisan key-generate jwt-secret routes migrate-up migrate-status
.PHONY: setup-deploy check-arch check-deploy build-linux deploy upload restart status logs

help: ## Mavjud buyruqlarni ko‘rsatish
	@awk 'BEGIN {FS = ":.*## "; print "Buyruqlar:"} /^[[:alnum:]_-]+:.*## / {printf "  make %-16s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Serverni ishga tushirish
	$(GO) run .

build: ## Dasturni bin/ticket-backend fayliga yig‘ish
	mkdir -p bin
	$(GO) build -o bin/ticket-backend .

check-arch:
	@case "$$DEPLOY_ARCH" in \
		amd64|arm64) ;; \
		*) printf '%s\n' 'DEPLOY_ARCH amd64 yoki arm64 bo‘lishi kerak.' >&2; exit 1 ;; \
	esac

build-linux: check-arch ## Ubuntu uchun Linux binarini yig‘ish
	mkdir -p "bin/linux-$$DEPLOY_ARCH"
	CGO_ENABLED=0 GOOS=linux GOARCH="$$DEPLOY_ARCH" $(GO) build -trimpath -ldflags='-s -w' -o "bin/linux-$$DEPLOY_ARCH/ticket-backend" .

setup-deploy: ## deploy.mk sozlama faylini yaratish
	@if [ -e deploy.mk ]; then \
		printf '%s\n' 'deploy.mk mavjud. Sozlamalar saqlandi.'; \
	else \
		cp deploy.example.mk deploy.mk; \
		printf '%s\n' 'deploy.mk yaratildi. DEPLOY_SSH qiymatini to‘ldiring.'; \
	fi

check-deploy: check-arch
	@set -eu; \
	case "$$DEPLOY_SSH" in \
		''|-*|*[!a-zA-Z0-9_.@-]*) printf '%s\n' 'DEPLOY_SSH ni user@host yoki SSH alias ko‘rinishida yozing.' >&2; exit 1 ;; \
	esac; \
	case "$$DEPLOY_PORT" in \
		''|*[!0-9]*|??????*) printf '%s\n' 'DEPLOY_PORT 1–65535 oralig‘idagi son bo‘lishi kerak.' >&2; exit 1 ;; \
	esac; \
	if [ "$$DEPLOY_PORT" -lt 1 ] || [ "$$DEPLOY_PORT" -gt 65535 ]; then \
		printf '%s\n' 'DEPLOY_PORT 1–65535 oralig‘ida bo‘lishi kerak.' >&2; exit 1; \
	fi; \
	case "$$DEPLOY_PATH" in \
		/*) ;; \
		*) printf '%s\n' 'DEPLOY_PATH serverdagi to‘liq papka yo‘li bo‘lishi kerak.' >&2; exit 1 ;; \
	esac; \
	case "$$DEPLOY_PATH" in \
		/|*[!a-zA-Z0-9_./-]*) printf '%s\n' 'DEPLOY_PATH uchun alohida backend papkasini yozing; bo‘shliq ishlatmang.' >&2; exit 1 ;; \
	esac; \
	case "$$DEPLOY_PATH" in \
		*[!/]*) ;; \
		*) printf '%s\n' 'DEPLOY_PATH ildiz papka bo‘lmasligi kerak.' >&2; exit 1 ;; \
	esac; \
	case "/$$DEPLOY_PATH/" in \
		*/../*|*/./*) printf '%s\n' 'DEPLOY_PATH ichida . yoki .. bo‘lmasin.' >&2; exit 1 ;; \
	esac; \
	case "$$DEPLOY_SERVICE" in \
		''|-*|*[!a-zA-Z0-9_-]*) printf '%s\n' 'DEPLOY_SERVICE ni .service qo‘shimchasiz yozing.' >&2; exit 1 ;; \
	esac; \
	command -v "$(SSH)" >/dev/null; \
	command -v "$(SCP)" >/dev/null

deploy: check-deploy ## Linux build, serverga yuklash va xizmatni qayta ishga tushirish
	$(MAKE) build-linux
	$(MAKE) upload
	$(MAKE) restart

upload: check-deploy ## Tayyor Linux binari, public va resources papkalarini yuklash
	@set -eu; \
	if [ ! -s "bin/linux-$$DEPLOY_ARCH/ticket-backend" ]; then \
		printf '%s\n' 'Linux binari topilmadi. Avval make build-linux qiling.' >&2; exit 1; \
	fi; \
	$(SSH) -p "$$DEPLOY_PORT" "$$DEPLOY_SSH" "test -d '$$DEPLOY_PATH' && test -w '$$DEPLOY_PATH' || { printf '%s\\n' 'Backend papkasi mavjud emas yoki unga yozish ruxsati yo‘q. DEPLOY.md bo‘yicha tayyorlang.' >&2; exit 1; }"; \
	printf 'Yuklanmoqda: %s:%s\n' "$$DEPLOY_SSH" "$$DEPLOY_PATH"; \
	$(SCP) -r -P "$$DEPLOY_PORT" public resources .env.example "$$DEPLOY_SSH:$$DEPLOY_PATH/"; \
	$(SCP) -P "$$DEPLOY_PORT" "bin/linux-$$DEPLOY_ARCH/ticket-backend" "$$DEPLOY_SSH:$$DEPLOY_PATH/.ticket-backend.upload"; \
	$(SSH) -p "$$DEPLOY_PORT" "$$DEPLOY_SSH" "set -eu; chmod 755 '$$DEPLOY_PATH/.ticket-backend.upload'; mv -f '$$DEPLOY_PATH/.ticket-backend.upload' '$$DEPLOY_PATH/ticket-backend'"; \
	printf '%s\n' 'Backend yuklandi. Ishga tushirish yoki yangilash uchun make restart qiling.'

restart: check-deploy ## systemd xizmatini qayta ishga tushirish
	@$(SSH) -t -p "$$DEPLOY_PORT" "$$DEPLOY_SSH" "sudo systemctl restart '$$DEPLOY_SERVICE.service' && sudo systemctl is-active --quiet '$$DEPLOY_SERVICE.service' && printf '%s\\n' 'Backend xizmati faol.'"

status: check-deploy ## Serverdagi backend holatini ko‘rish
	@$(SSH) -t -p "$$DEPLOY_PORT" "$$DEPLOY_SSH" "sudo systemctl status '$$DEPLOY_SERVICE.service' --no-pager"

logs: check-deploy ## Serverdagi backend jarayoni loglarini kuzatish
	@$(SSH) -t -p "$$DEPLOY_PORT" "$$DEPLOY_SSH" "sudo journalctl -u '$$DEPLOY_SERVICE.service' -n 100 -f"

tidy: ## Go kutubxonalarini tartibga keltirish
	$(GO) mod tidy

fmt: ## Go fayllarini formatlash
	gofmt -w main.go app bootstrap config database routes

vet: ## Statik tekshiruv
	$(GO) vet ./...

artisan: ## Istalgan Artisan buyrug‘i: make artisan ARGS="make:model User"
	$(ARTISAN) $(ARGS)

key-generate: ## APP_KEY yaratish
	$(ARTISAN) key:generate

jwt-secret: ## JWT_SECRET yaratish
	$(ARTISAN) jwt:secret

routes: ## HTTP marshrutlarini ko‘rsatish
	$(ARTISAN) route:list

migrate-up: ## Yangi migratsiyalarni qo‘llash
	$(ARTISAN) migrate

migrate-status: ## Migratsiyalar holatini ko‘rish
	$(ARTISAN) migrate:status
