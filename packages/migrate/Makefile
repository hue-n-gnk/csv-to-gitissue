.PHONY: install migrate

# Set default env as develop.
ifeq ($(ENV),)
ENV := develop
endif

-include ./.env.${ENV}

quotestr = $(subst ",,$1)

POSTGRES_USER := $(call quotestr,$(SECRET_POSTGRES_USER))
POSTGRES_PWD := $(call quotestr,$(SECRET_POSTGRES_PASSWORD))
POSTGRES_HOST := $(call quotestr,$(SECRET_POSTGRES_HOSTNAME))
POSTGRES_DB := $(call quotestr,$(SECRET_POSTGRES_DATABASE))

MIGRATE  := migrate -source file://ddl  -database 'postgres://$(POSTGRES_USER):$(POSTGRES_PWD)@$(POSTGRES_HOST):5432/$(POSTGRES_DB)?sslmode=disable'

install:
	go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate:
	@echo "Migrating..."
	$(MIGRATE) up
	@echo "Done!"

drop:
	@echo "Droping..."
	$(MIGRATE) down
	@echo "Done!"
